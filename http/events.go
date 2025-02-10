// Copyright © 2020 - 2024 Attestant Limited.
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

	consensusclient "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context, topics []string, handler consensusclient.EventHandlerFunc) error {
	if err := s.assertIsActive(ctx); err != nil {
		return err
	}
	if len(topics) == 0 {
		return errors.Join(errors.New("no topics supplied"), consensusclient.ErrInvalidOptions)
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Str("address", s.address).Logger()
	ctx = log.WithContext(ctx)

	// Ensure we support the requested topic(s).
	for i := range topics {
		if _, exists := api.SupportedEventTopics[topics[i]]; !exists {
			return fmt.Errorf("unsupported event topic %s", topics[i])
		}
	}

	endpoint := "/eth/v1/events"
	query := "topics=" + strings.Join(topics, "&topics=")
	callURL := urlForCall(s.base, endpoint, query)
	log.Trace().Str("url", callURL.String()).Msg("GET request to events stream")

	client := sse.NewClient(callURL.String())
	for k, v := range s.extraHeaders {
		client.Headers[k] = v
	}
	if _, exists := client.Headers["User-Agent"]; !exists {
		client.Headers["User-Agent"] = defaultUserAgent
	}
	client.Headers["Accept"] = "text/event-stream"
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
func (*Service) handleEvent(ctx context.Context, msg *sse.Event, handler consensusclient.EventHandlerFunc) {
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
	case "attestation":
		data := &phase0.Attestation{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attestation")

			return
		}
		event.Data = data
	case "attester_slashing":
		data := &phase0.AttesterSlashing{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse attester slashing event")

			return
		}
		event.Data = data
	case "blob_sidecar":
		data := &api.BlobSidecarEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse blob sidecar event")

			return
		}
		event.Data = data
	case "block":
		data := &api.BlockEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse block event")

			return
		}
		event.Data = data
	case "block_gossip":
		data := &api.BlockGossipEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse block gossip event")

			return
		}
		event.Data = data
	case "bls_to_execution_change":
		data := &capella.SignedBLSToExecutionChange{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse bls to execution change event")

			return
		}
		event.Data = data
	case "chain_reorg":
		data := &api.ChainReorgEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse chain reorg event")

			return
		}
		event.Data = data
	case "contribution_and_proof":
		data := &altair.SignedContributionAndProof{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse contribution and proof event")

			return
		}
		event.Data = data
	case "finalized_checkpoint":
		data := &api.FinalizedCheckpointEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse finalized checkpoint event")

			return
		}
		event.Data = data
	case "head":
		data := &api.HeadEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse head event")

			return
		}
		event.Data = data
	case "payload_attributes":
		data := &api.PayloadAttributesEvent{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse payload attributes event")

			return
		}
		event.Data = data
	case "proposer_slashing":
		data := &phase0.ProposerSlashing{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse proposer slashing event")

			return
		}
		event.Data = data
	case "voluntary_exit":
		data := &phase0.SignedVoluntaryExit{}
		err := json.Unmarshal(msg.Data, data)
		if err != nil {
			log.Error().Err(err).RawJSON("data", msg.Data).Msg("Failed to parse voluntary exit")

			return
		}
		event.Data = data
	case "":
		// Used as keepalive.  Ignore.
		return
	default:
		log.Warn().Str("topic", string(msg.Event)).Msg("Received message with unhandled topic; ignoring")

		return
	}
	handler(event)
}
