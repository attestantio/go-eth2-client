// Copyright © 2021, 2025 Attestant Limited.
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

package multi

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context,
	opts *api.EventsOpts,
) error {
	if opts == nil {
		return consensusclient.ErrNoOptions
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Logger()

	// Because events are streams we treat them differently from all other calls.
	// We listen to all active clients, and only pass along events from the currently active provider.

	// Grab local copy of both active and inactive clients in case it is updated whilst we are using it.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	inactiveClients := s.inactiveClients
	s.clientsMu.RUnlock()

	// Call all active clients immediately.
	for _, client := range activeClients {
		ah := &activeHandler{
			s:       s,
			log:     log.With().Logger(),
			address: client.Address(),
			opts: &api.EventsOpts{
				Common: opts.Common,
			},
		}
		ah.opts.Handler = ah.genericHandler
		ah.opts.AttestationHandler = ah.attestationHandler
		ah.opts.AttesterSlashingHandler = ah.attesterSlashingHandler
		ah.opts.BlobSidecarHandler = ah.blobSidecarHandler
		ah.opts.BLSToExecutionChangeHandler = ah.blsToExecutionChangeHandler
		ah.opts.ChainReorgHandler = ah.chainReorgHandler
		ah.opts.ContributionAndProofHandler = ah.contributionAndProofHandler
		ah.opts.FinalizedCheckpointHandler = ah.finalizedCheckpointHandler
		ah.opts.HeadHandler = ah.headHandler
		ah.opts.PayloadAttributesHandler = ah.payloadAttributesHandler
		ah.opts.ProposerSlashingHandler = ah.proposerSlashingHandler
		ah.opts.SingleAttestationHandler = ah.singleAttestationHandler
		ah.opts.VoluntaryExitHandler = ah.voluntaryExitHandler

		if err := client.(consensusclient.EventsProvider).Events(ctx, ah.opts); err != nil {
			inactiveClients = append(inactiveClients, client)

			continue
		}

		log.Trace().Str("address", ah.address).Strs("topics", opts.Topics).Msg("Events handler active")
	}

	// Periodically try all inactive clients, quitting as they become active.
	for _, inactiveClient := range inactiveClients {
		ah := &activeHandler{
			s:       s,
			log:     log.With().Logger(),
			address: inactiveClient.Address(),
			opts:    opts,
		}
		go func(c consensusclient.Service, ah *activeHandler) {
			for {
				provider, isProvider := c.(consensusclient.NodeSyncingProvider)
				if !isProvider {
					ah.log.Error().
						Str("address", ah.address).
						Strs("topics", opts.Topics).
						Msg("Not a node syncing provider")

					return
				}

				syncResponse, err := provider.NodeSyncing(ctx, &api.NodeSyncingOpts{})
				if err != nil {
					ah.log.Error().
						Str("address", ah.address).
						Strs("topics", opts.Topics).
						Err(err).
						Msg("Failed to obtain sync state from node")

					return
				}

				if !syncResponse.Data.IsSyncing {
					// Client is now synced, set up the events call.
					if err := c.(consensusclient.EventsProvider).Events(ctx, opts); err != nil {
						ah.log.Error().
							Str("address", ah.address).
							Strs("topics", opts.Topics).
							Err(err).
							Msg("Failed to set up events handler")
					}

					// Return either way.
					return
				}

				time.Sleep(5 * time.Second)
			}
		}(inactiveClient, ah)
	}

	return nil
}

type activeHandler struct {
	s       *Service
	log     zerolog.Logger
	address string
	opts    *api.EventsOpts
}

func (h *activeHandler) attestationHandler(ctx context.Context, data *spec.VersionedAttestation) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Attestation event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.AttestationHandler(ctx, data)
}

func (h *activeHandler) attesterSlashingHandler(ctx context.Context, data *electra.AttesterSlashing) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Attester slashing event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.AttesterSlashingHandler(ctx, data)
}

func (h *activeHandler) blobSidecarHandler(ctx context.Context, data *apiv1.BlobSidecarEvent) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Blob sidecar event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.BlobSidecarHandler(ctx, data)
}

func (h *activeHandler) blsToExecutionChangeHandler(ctx context.Context, data *capella.SignedBLSToExecutionChange) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("BLS to execution change event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.BLSToExecutionChangeHandler(ctx, data)
}

func (h *activeHandler) chainReorgHandler(ctx context.Context, data *apiv1.ChainReorgEvent) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Chain reorg event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.ChainReorgHandler(ctx, data)
}

func (h *activeHandler) contributionAndProofHandler(ctx context.Context, data *altair.SignedContributionAndProof) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Chain reorg event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.ContributionAndProofHandler(ctx, data)
}

func (h *activeHandler) finalizedCheckpointHandler(ctx context.Context, data *apiv1.FinalizedCheckpointEvent) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Finalized checkpoint event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.FinalizedCheckpointHandler(ctx, data)
}

func (h *activeHandler) headHandler(ctx context.Context, data *apiv1.HeadEvent) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Head event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.HeadHandler(ctx, data)
}

func (h *activeHandler) payloadAttributesHandler(ctx context.Context, data *apiv1.PayloadAttributesEvent) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Payload attributes event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.PayloadAttributesHandler(ctx, data)
}

func (h *activeHandler) proposerSlashingHandler(ctx context.Context, data *phase0.ProposerSlashing) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Proposer slashing event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.ProposerSlashingHandler(ctx, data)
}

func (h *activeHandler) singleAttestationHandler(ctx context.Context, data *electra.SingleAttestation) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Single attestation event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.SingleAttestationHandler(ctx, data)
}

func (h *activeHandler) voluntaryExitHandler(ctx context.Context, data *phase0.SignedVoluntaryExit) {
	log := h.log.With().Str("address", h.address).Logger()
	log.Trace().Msg("Voluntary exit event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	h.opts.VoluntaryExitHandler(ctx, data)
}

func (h *activeHandler) genericHandler(event *apiv1.Event) {
	log := h.log.With().Str("address", h.address).Str("topic", event.Topic).Logger()
	log.Trace().Msg("Event received")

	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() != h.address {
		return
	}

	log.Trace().Msg("Forwarding due to primary active address")

	if h.opts.Handler != nil {
		h.opts.Handler(event)
	}
}
