// Copyright © 2020 - 2025 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context, opts *api.EventsOpts) error {
	if err := s.assertIsActive(ctx); err != nil {
		return err
	}
	if opts == nil {
		return client.ErrNoOptions
	}
	if len(opts.Topics) == 0 {
		return errors.Join(errors.New("no topics supplied"), client.ErrInvalidOptions)
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Str("address", s.address).Logger()
	ctx = log.WithContext(ctx)

	if err := s.checkEventsOpts(opts); err != nil {
		return err
	}

	endpoint := "/eth/v1/events"
	query := "topics=" + strings.Join(opts.Topics, "&topics=")
	callURL := urlForCall(s.base, endpoint, query)
	log.Trace().Str("url", callURL.String()).Msg("GET request to events stream")

	sseClient := sse.NewClient(callURL.String())
	for k, v := range s.extraHeaders {
		sseClient.Headers[k] = v
	}
	if _, exists := sseClient.Headers["User-Agent"]; !exists {
		sseClient.Headers["User-Agent"] = defaultUserAgent
	}
	sseClient.Headers["Accept"] = "text/event-stream"
	sseClient.Connection.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
		}).Dial,
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				log.Trace().Msg("Connecting to events stream")
				if err := sseClient.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
					s.handleEvent(ctx, msg, opts)
				}); err != nil {
					log.Error().Err(err).Msg("Failed to subscribe to event stream")
				}
				log.Trace().Msg("Events stream disconnected")
			case <-ctx.Done():
				log.Debug().Msg("Context done")

				return
			}
		}
	}()

	return nil
}

func (s *Service) checkEventsOpts(opts *api.EventsOpts) error {
	// Ensure we support the requested topic(s), and have a handler for each.
	for _, topic := range opts.Topics {
		if _, exists := apiv1.SupportedEventTopics[topic]; !exists {
			return fmt.Errorf("unsupported event topic %s", topic)
		}
		if opts.Handler != nil {
			// There is a generic handler in place, no further checks for this topic required.
			continue
		}
		if err := s.checkEventSpecificHandler(opts, topic); err != nil {
			return err
		}
	}

	return nil
}

func (*Service) checkEventSpecificHandler(opts *api.EventsOpts, topic string) error {
	var hasHandler bool

	switch topic {
	case "attestation":
		hasHandler = opts.AttestationHandler != nil
	case "attester_slashing":
		hasHandler = opts.AttesterSlashingHandler != nil
	case "blob_sidecar":
		hasHandler = opts.BlobSidecarHandler != nil
	case "block":
		hasHandler = opts.BlockHandler != nil
	case "block_gossip":
		hasHandler = opts.BlockGossipHandler != nil
	case "bls_to_execution_change":
		hasHandler = opts.BLSToExecutionChangeHandler != nil
	case "chain_reorg":
		hasHandler = opts.ChainReorgHandler != nil
	case "contribution_and_proof":
		hasHandler = opts.ContributionAndProofHandler != nil
	case "data_column_sidecar":
		hasHandler = opts.DataColumnSidecarHandler != nil
	case "finalized_checkpoint":
		hasHandler = opts.FinalizedCheckpointHandler != nil
	case "head":
		hasHandler = opts.HeadHandler != nil
	case "payload_attributes":
		hasHandler = opts.PayloadAttributesHandler != nil
	case "proposer_slashing":
		hasHandler = opts.ProposerSlashingHandler != nil
	case "single_attestation":
		hasHandler = opts.SingleAttestationHandler != nil
	case "voluntary_exit":
		hasHandler = opts.VoluntaryExitHandler != nil
	default:
		return fmt.Errorf("unsupported event %s", topic)
	}

	if !hasHandler {
		return fmt.Errorf("no handler for %s event", topic)
	}

	return nil
}

// handleEvent handles all events.
func (s *Service) handleEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)

	if msg == nil {
		log.Debug().Msg("No message supplied; ignoring")

		return
	}

	switch string(msg.Event) {
	case "attestation":
		s.handleAttestationEvent(ctx, msg, opts)
	case "attester_slashing":
		s.handleAttesterSlashingEvent(ctx, msg, opts)
	case "blob_sidecar":
		s.handleBlobSidecarEvent(ctx, msg, opts)
	case "block":
		s.handleBlockEvent(ctx, msg, opts)
	case "block_gossip":
		s.handleBlockGossipEvent(ctx, msg, opts)
	case "bls_to_execution_change":
		s.handleBLSToExecutionChangeEvent(ctx, msg, opts)
	case "chain_reorg":
		s.handleChainReorgEvent(ctx, msg, opts)
	case "contribution_and_proof":
		s.handleContributionAndProofEvent(ctx, msg, opts)
	case "data_column_sidecar":
		s.handleDataColumnSidecarEvent(ctx, msg, opts)
	case "finalized_checkpoint":
		s.handleFinalizedCheckpointEvent(ctx, msg, opts)
	case "head":
		s.handleHeadEvent(ctx, msg, opts)
	case "payload_attributes":
		s.handlePayloadAttributesEvent(ctx, msg, opts)
	case "proposer_slashing":
		s.handleProposerSlashingEvent(ctx, msg, opts)
	case "single_attestation":
		s.handleSingleAttestationEvent(ctx, msg, opts)
	case "voluntary_exit":
		s.handleVoluntaryExitEvent(ctx, msg, opts)
	case "":
		// Used as keepalive.  Ignore.
	default:
		log.Warn().Str("topic", string(msg.Event)).Msg("Received message with unhandled topic; ignoring")
	}
}

func (*Service) handleAttestationEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &spec.VersionedAttestation{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attestation")

		return
	}

	switch {
	case opts.AttestationHandler != nil:
		opts.AttestationHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleAttesterSlashingEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &electra.AttesterSlashing{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attester slashing event")

		return
	}

	switch {
	case opts.AttesterSlashingHandler != nil:
		opts.AttesterSlashingHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleBlobSidecarEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.BlobSidecarEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse blob sidecar event")

		return
	}

	switch {
	case opts.BlobSidecarHandler != nil:
		opts.BlobSidecarHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleBlockEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.BlockEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse block event")

		return
	}

	switch {
	case opts.BlockHandler != nil:
		opts.BlockHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleBlockGossipEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.BlockGossipEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse block gossip event")

		return
	}

	switch {
	case opts.BlockGossipHandler != nil:
		opts.BlockGossipHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleBLSToExecutionChangeEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &capella.SignedBLSToExecutionChange{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse bls to execution change event")

		return
	}

	switch {
	case opts.BLSToExecutionChangeHandler != nil:
		opts.BLSToExecutionChangeHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleChainReorgEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.ChainReorgEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse chain reorg event")

		return
	}

	switch {
	case opts.ChainReorgHandler != nil:
		opts.ChainReorgHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleContributionAndProofEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &altair.SignedContributionAndProof{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse contribution and proof event")

		return
	}

	switch {
	case opts.ContributionAndProofHandler != nil:
		opts.ContributionAndProofHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleFinalizedCheckpointEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.FinalizedCheckpointEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse finalized checkpoint event")

		return
	}

	switch {
	case opts.FinalizedCheckpointHandler != nil:
		opts.FinalizedCheckpointHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleHeadEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.HeadEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse head event")

		return
	}

	switch {
	case opts.HeadHandler != nil:
		opts.HeadHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handlePayloadAttributesEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.PayloadAttributesEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse payload attributes event")

		return
	}

	switch {
	case opts.PayloadAttributesHandler != nil:
		opts.PayloadAttributesHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleProposerSlashingEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &phase0.ProposerSlashing{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse proposer slashing event")

		return
	}

	switch {
	case opts.ProposerSlashingHandler != nil:
		opts.ProposerSlashingHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleSingleAttestationEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &electra.SingleAttestation{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse single attestation")

		return
	}

	switch {
	case opts.SingleAttestationHandler != nil:
		opts.SingleAttestationHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleVoluntaryExitEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &phase0.SignedVoluntaryExit{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse voluntary exit")

		return
	}

	switch {
	case opts.VoluntaryExitHandler != nil:
		opts.VoluntaryExitHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}

func (*Service) handleDataColumnSidecarEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)
	data := &apiv1.DataColumnSidecarEvent{}
	err := json.Unmarshal(msg.Data, data)
	if err != nil {
		log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse data column sidecar event")

		return
	}

	switch {
	case opts.DataColumnSidecarHandler != nil:
		opts.DataColumnSidecarHandler(ctx, data)
	case opts.Handler != nil:
		opts.Handler(&apiv1.Event{
			Topic: string(msg.Event),
			Data:  data,
		})
	default:
		log.Debug().Msg("No specific or generic handler supplied; ignoring")
	}
}
