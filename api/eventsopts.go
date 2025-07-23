// Copyright © 2025 Attestant Limited.
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

package api

import (
	"context"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// EventsOpts are the options for obtaining events.
type EventsOpts struct {
	Common CommonOpts

	// Topics are the topics of events to which we want to listen.
	Topics []string

	// Handler is a generic handler function to which to send all events.
	// In general, it is better to use event-specific handlers as they avoid casting, and also provide a context.
	Handler EventHandlerFunc

	// AttestationHandler is a handler for the attestation event.
	AttestationHandler AttestationEventHandlerFunc
	// AttesterSlashingHandler is a handler for the attester_slashing event.
	AttesterSlashingHandler AttesterSlashingEventHandlerFunc
	// BlobSidecarHandler is a handler for the blob_sidecar event.
	BlobSidecarHandler BlobSidecarEventHandlerFunc
	// BlockHandler is a handler for the block event.
	BlockHandler BlockEventHandlerFunc
	// BlockGossipHandler is a handler for the block_gossip event.
	BlockGossipHandler BlockGossipEventHandlerFunc
	// BLSToExecutionChangeHandler is a handler for the bls_to_execution_change event.
	BLSToExecutionChangeHandler BLSToExecutionChangeEventHandlerFunc
	// ChainReorgHandler is a handler for the chain_reorg event.
	ChainReorgHandler ChainReorgEventHandlerFunc
	// ContributionAndProofHandler is a handler for the contribution_and_proof event.
	ContributionAndProofHandler ContributionAndProofEventHandlerFunc
	// DataColumnSidecarHandler is a handler for the data_column_sidecar event.
	DataColumnSidecarHandler DataColumnSidecarEventHandlerFunc
	// FinalizedCheckpointHandler is a handler for the finalized_checkpoint event.
	FinalizedCheckpointHandler FinalizedCheckpointEventHandlerFunc
	// HeadHandler is a handler for the head event.
	HeadHandler HeadEventHandlerFunc
	// PayloadAttributesHandler is a handler for the payload_attributes event.
	PayloadAttributesHandler PayloadAttributesEventHandlerFunc
	// ProposerSlashingHandler is a handler for the proposer_slashing event.
	ProposerSlashingHandler ProposerSlashingEventHandlerFunc
	// SingleAttestationHandler is a handler for the single_attestation event.
	SingleAttestationHandler SingleAttestationEventHandlerFunc
	// VoluntaryExitHandler is a handler for the voluntary_exit event.
	VoluntaryExitHandler VoluntaryExitEventHandlerFunc
}

// EventHandlerFunc is the handler for generic events.
type EventHandlerFunc func(*apiv1.Event)

// AttestationEventHandlerFunc is the handler for attestation events.
type AttestationEventHandlerFunc func(context.Context, *spec.VersionedAttestation)

// AttesterSlashingEventHandlerFunc is the handler for attestation_slashing events.
type AttesterSlashingEventHandlerFunc func(context.Context, *electra.AttesterSlashing)

// BlobSidecarEventHandlerFunc is the handler for blob_sidecar events.
type BlobSidecarEventHandlerFunc func(context.Context, *apiv1.BlobSidecarEvent)

// BlockEventHandlerFunc is the handler for block events.
type BlockEventHandlerFunc func(context.Context, *apiv1.BlockEvent)

// BlockGossipEventHandlerFunc is the handler for block_gossip events.
type BlockGossipEventHandlerFunc func(context.Context, *apiv1.BlockGossipEvent)

// BLSToExecutionChangeEventHandlerFunc is the handler for bls_to_execution_change events.
type BLSToExecutionChangeEventHandlerFunc func(context.Context, *capella.SignedBLSToExecutionChange)

// ChainReorgEventHandlerFunc is the handler for chain_reorg events.
type ChainReorgEventHandlerFunc func(context.Context, *apiv1.ChainReorgEvent)

// ContributionAndProofEventHandlerFunc is the handler for contribution_and_proof events.
type ContributionAndProofEventHandlerFunc func(context.Context, *altair.SignedContributionAndProof)

// FinalizedCheckpointEventHandlerFunc is the handler for finalized_checkpoint events.
type FinalizedCheckpointEventHandlerFunc func(context.Context, *apiv1.FinalizedCheckpointEvent)

// HeadEventHandlerFunc is the handler for head events.
type HeadEventHandlerFunc func(context.Context, *apiv1.HeadEvent)

// PayloadAttributesEventHandlerFunc is the handler for payload_attributes events.
type PayloadAttributesEventHandlerFunc func(context.Context, *apiv1.PayloadAttributesEvent)

// ProposerSlashingEventHandlerFunc is the handler for proposer_slashing events.
type ProposerSlashingEventHandlerFunc func(context.Context, *phase0.ProposerSlashing)

// SingleAttestationEventHandlerFunc is the handler for single_attestation events.
type SingleAttestationEventHandlerFunc func(context.Context, *electra.SingleAttestation)

// VoluntaryExitEventHandlerFunc is the handler for voluntary_exit events.
type VoluntaryExitEventHandlerFunc func(context.Context, *phase0.SignedVoluntaryExit)

// DataColumnSidecarEventHandlerFunc is the handler for data_column_sidecar events.
type DataColumnSidecarEventHandlerFunc func(context.Context, *apiv1.DataColumnSidecarEvent)
