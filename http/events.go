// Copyright © 2020 - 2022 Attestant Limited.
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
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context, topics []string, handler client.EventHandlerFunc) error {
	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Str("address", s.address).Logger()
	ctx = log.WithContext(ctx)

	if len(topics) == 0 {
		return errors.New("no topics supplied")
	}

	// Ensure we support the requested topic(s).
	for i := range topics {
		if _, exists := api.SupportedEventTopics[topics[i]]; !exists {
			return fmt.Errorf("unsupported event topic %s", topics[i])
		}
	}

	reference, err := url.Parse(fmt.Sprintf("eth/v1/events?topics=%s", strings.Join(topics, "&topics=")))
	if err != nil {
		return errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()
	log.Trace().Str("url", url).Msg("GET request to events stream")

	client := sse.NewClient(url)
	client.Connection.Transport = &http.Transport{
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
				if err := client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
					s.handleEvent(ctx, msg, handler)
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

// handleEvent parses an event and passes it on to the handler.
func (s *Service) handleEvent(ctx context.Context, msg *sse.Event, handler client.EventHandlerFunc) {
	log := zerolog.Ctx(ctx)

	if handler == nil {
		log.Debug().Msg("No handler supplied; ignoring")

		return
	}
	if msg == nil {
		log.Debug().Msg("No message supplied; ignoring")

		return
	}

	event := &api.Event{
		Topic: string(msg.Event),
	}
	switch string(msg.Event) {
	case "head":
		headEvent := &api.HeadEvent{}
		err := json.Unmarshal(msg.Data, headEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse head event")

			return
		}
		event.Data = headEvent
	case "block":
		blockEvent := &api.BlockEvent{}
		err := json.Unmarshal(msg.Data, blockEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse block event")

			return
		}
		event.Data = blockEvent
	case "attestation":
		attestation := &phase0.Attestation{}
		err := json.Unmarshal(msg.Data, attestation)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attestation")

			return
		}
		event.Data = attestation
	case "voluntary_exit":
		voluntaryExit := &phase0.SignedVoluntaryExit{}
		err := json.Unmarshal(msg.Data, voluntaryExit)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse voluntary exit")

			return
		}
		event.Data = voluntaryExit
	case "finalized_checkpoint":
		finalizedCheckpointEvent := &api.FinalizedCheckpointEvent{}
		err := json.Unmarshal(msg.Data, finalizedCheckpointEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse finalized checkpoint event")

			return
		}
		event.Data = finalizedCheckpointEvent
	case "chain_reorg":
		chainReorgEvent := &api.ChainReorgEvent{}
		err := json.Unmarshal(msg.Data, chainReorgEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse chain reorg event")

			return
		}
		event.Data = chainReorgEvent
	case "contribution_and_proof":
		contributionAndProofEvent := &altair.SignedContributionAndProof{}
		err := json.Unmarshal(msg.Data, contributionAndProofEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse contribution and proof event")

			return
		}
		event.Data = contributionAndProofEvent
	case "payload_attributes":
		payloadAttributesEvent := &api.PayloadAttributesEvent{}
		err := json.Unmarshal(msg.Data, payloadAttributesEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse payload attributes event")

			return
		}
		event.Data = payloadAttributesEvent
	case "proposer_slashing":
		proposerSlashingEvent := &phase0.ProposerSlashing{}
		err := json.Unmarshal(msg.Data, proposerSlashingEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse proposer slashing event")

			return
		}
		event.Data = proposerSlashingEvent
	case "attester_slashing":
		attesterSlashingEvent := &phase0.AttesterSlashing{}
		err := json.Unmarshal(msg.Data, attesterSlashingEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attester slashing event")

			return
		}
		event.Data = attesterSlashingEvent
	case "bls_to_execution_change":
		blsToExecutionChangeEvent := &capella.BLSToExecutionChange{}
		err := json.Unmarshal(msg.Data, blsToExecutionChangeEvent)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse bls to execution change event")

			return
		}
		event.Data = blsToExecutionChangeEvent
	case "blob_sidecar":
		blobSidecar := &api.BlobSidecarEvent{}
		err := json.Unmarshal(msg.Data, blobSidecar)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse blob sidecar event")

			return
		}
		event.Data = blobSidecar
	case "":
		// Used as keepalive.  Ignore.
		return
	default:
		log.Warn().Str("topic", string(msg.Event)).Msg("Received message with unhandled topic; ignoring")

		return
	}
	handler(event)
}
