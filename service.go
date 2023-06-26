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

package client

import (
	"context"
	"time"

	api "github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Service is the service providing a connection to an Ethereum 2 client.
type Service interface {
	// Name returns the name of the client implementation.
	Name() string

	// Address returns the address of the client.
	Address() string
}

// EpochFromStateIDProvider is the interface for providing epochs from state IDs.
type EpochFromStateIDProvider interface {
	// EpochFromStateID converts a state ID to its epoch.
	EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error)
}

// SlotFromStateIDProvider is the interface for providing slots from state IDs.
type SlotFromStateIDProvider interface {
	// SlotFromStateID converts a state ID to its slot.
	SlotFromStateID(ctx context.Context, stateID string) (phase0.Slot, error)
}

// NodeVersionProvider is the interface for providing the node version.
type NodeVersionProvider interface {
	// NodeVersion returns a free-text string with the node version.
	NodeVersion(ctx context.Context) (string, error)
}

// SlotDurationProvider is the interface for providing the duration of each slot of a chain.
type SlotDurationProvider interface {
	// SlotDuration provides the duration of a slot of the chain.
	SlotDuration(ctx context.Context) (time.Duration, error)
}

// SlotsPerEpochProvider is the interface for providing the number of slots in each epoch of a chain.
type SlotsPerEpochProvider interface {
	// SlotsPerEpoch provides the slots per epoch of the chain.
	SlotsPerEpoch(ctx context.Context) (uint64, error)
}

// FarFutureEpochProvider is the interface for providing the far future epoch of a chain.
type FarFutureEpochProvider interface {
	// FarFutureEpoch provides the far future epoch of the chain.
	FarFutureEpoch(ctx context.Context) (phase0.Epoch, error)
}

// GenesisValidatorsRootProvider is the interface for providing the genesis validators root of a chain.
type GenesisValidatorsRootProvider interface {
	// GenesisValidatorsRoot provides the genesis validators root of the chain.
	GenesisValidatorsRoot(ctx context.Context) ([]byte, error)
}

// TargetAggregatorsPerCommitteeProvider is the interface for providing the target number of
// aggregators in each attestation committee.
type TargetAggregatorsPerCommitteeProvider interface {
	// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
	TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error)
}

// ValidatorIndexProvider is the interface for entities that can provide the index of a validator.
type ValidatorIndexProvider interface {
	// Index provides the index of the validator.
	Index(ctx context.Context) (phase0.ValidatorIndex, error)
}

// ValidatorPubKeyProvider is the interface for entities that can provide the public key of a validator.
type ValidatorPubKeyProvider interface {
	// PubKey provides the public key of the validator.
	PubKey(ctx context.Context) (phase0.BLSPubKey, error)
}

// ValidatorIDProvider is the interface that provides the identifiers (pubkey, index) of a validator.
type ValidatorIDProvider interface {
	ValidatorIndexProvider
	ValidatorPubKeyProvider
}

// DepositContractProvider is the interface for providng details about the deposit contract.
type DepositContractProvider interface {
	// DepositContract provides details of the Ethereum 1 deposit contract for the chain.
	DepositContract(ctx context.Context) (*apiv1.DepositContract, error)
}

// SignedBeaconBlockProvider is the interface for providing beacon blocks.
type SignedBeaconBlockProvider interface {
	// SignedBeaconBlock fetches a signed beacon block given a block ID.
	SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error)
}

// BeaconBlockBlobsProvider is the interface for providing blobs for a given beacon block.
type BeaconBlockBlobsProvider interface {
	// BeaconBlockBlobs fetches the blobs given a block ID.
	BeaconBlockBlobs(ctx context.Context, blockID string) ([]*deneb.BlobSidecar, error)
}

// BeaconCommitteesProvider is the interface for providing beacon committees.
type BeaconCommitteesProvider interface {
	// BeaconCommittees fetches all beacon committees for the epoch at the given state.
	BeaconCommittees(ctx context.Context, stateID string) ([]*apiv1.BeaconCommittee, error)

	// BeaconCommitteesAtEpoch fetches all beacon committees for the given epoch at the given state.
	BeaconCommitteesAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) ([]*apiv1.BeaconCommittee, error)
}

// SyncCommitteesProvider is the interface for providing sync committees.
type SyncCommitteesProvider interface {
	// SyncCommittee fetches the sync committee for the given state.
	SyncCommittee(ctx context.Context, stateID string) (*apiv1.SyncCommittee, error)

	// SyncCommitteeAtEpoch fetches the sync committee for the given epoch at the given state.
	SyncCommitteeAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) (*apiv1.SyncCommittee, error)
}

// EventHandlerFunc is the handler for events.
type EventHandlerFunc func(*apiv1.Event)

//
// Standard API
//

// AggregateAttestationProvider is the interface for providing aggregate attestations.
type AggregateAttestationProvider interface {
	// AggregateAttestation fetches the aggregate attestation given an attestation.
	AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error)
}

// AggregateAttestationsSubmitter is the interface for submitting aggregate attestations.
type AggregateAttestationsSubmitter interface {
	// SubmitAggregateAttestations submits aggregate attestations.
	SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error
}

// AttestationDataProvider is the interface for providing attestation data.
type AttestationDataProvider interface {
	// AttestationData fetches the attestation data for the given slot and committee index.
	AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error)
}

// AttestationPoolProvider is the interface for providing attestation pools.
type AttestationPoolProvider interface {
	// AttestationPool fetches the attestation pool for the given slot.
	AttestationPool(ctx context.Context, slot phase0.Slot) ([]*phase0.Attestation, error)
}

// AttestationsSubmitter is the interface for submitting attestations.
type AttestationsSubmitter interface {
	// SubmitAttestations submits attestations.
	SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error
}

// AttesterDutiesProvider is the interface for providing attester duties.
type AttesterDutiesProvider interface {
	// AttesterDuties obtains attester duties.
	// If validatorIndicess is nil it will return all duties for the given epoch.
	AttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*apiv1.AttesterDuty, error)
}

// SyncCommitteeDutiesProvider is the interface for providing sync committee duties.
type SyncCommitteeDutiesProvider interface {
	// SyncCommitteeDuties obtains sync committee duties.
	// If validatorIndicess is nil it will return all duties for the given epoch.
	SyncCommitteeDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*apiv1.SyncCommitteeDuty, error)
}

// SyncCommitteeMessagesSubmitter is the interface for submitting sync committee messages.
type SyncCommitteeMessagesSubmitter interface {
	// SubmitSyncCommitteeMessages submits sync committee messages.
	SubmitSyncCommitteeMessages(ctx context.Context, messages []*altair.SyncCommitteeMessage) error
}

// SyncCommitteeSubscriptionsSubmitter is the interface for submitting sync committee subnet subscription requests.
type SyncCommitteeSubscriptionsSubmitter interface {
	// SubmitSyncCommitteeSubscriptions subscribes to sync committees.
	SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.SyncCommitteeSubscription) error
}

// SyncCommitteeContributionProvider is the interface for providing sync committee contributions.
type SyncCommitteeContributionProvider interface {
	// SyncCommitteeContribution provides a sync committee contribution.
	SyncCommitteeContribution(ctx context.Context, slot phase0.Slot, subcommitteeIndex uint64, beaconBlockRoot phase0.Root) (*altair.SyncCommitteeContribution, error)
}

// SyncCommitteeContributionsSubmitter is the interface for submitting sync committee contributions.
type SyncCommitteeContributionsSubmitter interface {
	// SubmitSyncCommitteeContributions submits sync committee contributions.
	SubmitSyncCommitteeContributions(ctx context.Context, contributionAndProofs []*altair.SignedContributionAndProof) error
}

// BLSToExecutionChangesSubmitter is the interface for submitting BLS to execution address changes.
type BLSToExecutionChangesSubmitter interface {
	// SubmitBLSToExecutionChanges submits BLS to execution address change operations.
	SubmitBLSToExecutionChanges(ctx context.Context, blsToExecutionChanges []*capella.SignedBLSToExecutionChange) error
}

// BeaconBlockHeadersProvider is the interface for providing beacon block headers.
type BeaconBlockHeadersProvider interface {
	// BeaconBlockHeader provides the block header of a given block ID.
	BeaconBlockHeader(ctx context.Context, blockID string) (*apiv1.BeaconBlockHeader, error)
}

// BeaconBlockProposalProvider is the interface for providing beacon block proposals.
type BeaconBlockProposalProvider interface {
	// BeaconBlockProposal fetches a proposed beacon block for signing.
	BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error)
}

// BeaconBlockRootProvider is the interface for providing beacon block roots.
type BeaconBlockRootProvider interface {
	// BeaconBlockRoot fetches a block's root given a block ID.
	BeaconBlockRoot(ctx context.Context, blockID string) (*phase0.Root, error)
}

// BeaconBlockSubmitter is the interface for submitting beacon blocks.
type BeaconBlockSubmitter interface {
	// SubmitBeaconBlock submits a beacon block.
	SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error
}

// BeaconCommitteeSubscriptionsSubmitter is the interface for submitting beacon committee subnet subscription requests.
type BeaconCommitteeSubscriptionsSubmitter interface {
	// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
	SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error
}

// BeaconStateProvider is the interface for providing beacon state.
type BeaconStateProvider interface {
	// BeaconState fetches a beacon state given a state ID.
	BeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error)
}

// BeaconStateRandaoProvider is the interface for providing beacon state RANDAOs.
type BeaconStateRandaoProvider interface {
	// BeaconStateRandao fetches a beacon state RANDAO given a state ID.
	BeaconStateRandao(ctx context.Context, stateID string) (*phase0.Root, error)
}

// BeaconStateRootProvider is the interface for providing beacon state roots.
type BeaconStateRootProvider interface {
	// BeaconStateRoot fetches a beacon state root given a state ID.
	BeaconStateRoot(ctx context.Context, stateID string) (*phase0.Root, error)
}

// BlindedBeaconBlockProposalProvider is the interface for providing blinded beacon block proposals.
type BlindedBeaconBlockProposalProvider interface {
	// BlindedBeaconBlockProposal fetches a blinded proposed beacon block for signing.
	BlindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*api.VersionedBlindedBeaconBlock, error)
}

// BlindedBeaconBlockSubmitter is the interface for submitting blinded beacon blocks.
type BlindedBeaconBlockSubmitter interface {
	// SubmitBlindedBeaconBlock submits a beacon block.
	SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error
}

// ValidatorRegistrationsSubmitter is the interface for submitting validator registrations.
type ValidatorRegistrationsSubmitter interface {
	// SubmitValidatorRegistrations submits a validator registration.
	SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error
}

// EventsProvider is the interface for providing events.
type EventsProvider interface {
	// Events feeds requested events with the given topics to the supplied handler.
	Events(ctx context.Context, topics []string, handler EventHandlerFunc) error
}

// FinalityProvider is the interface for providing finality information.
type FinalityProvider interface {
	// Finality provides the finality given a state ID.
	Finality(ctx context.Context, stateID string) (*apiv1.Finality, error)
}

// ForkChoiceProvider is the interface for providing fork choice information.
type ForkChoiceProvider interface {
	// Fork fetches all current fork choice context.
	ForkChoice(ctx context.Context) (*apiv1.ForkChoice, error)
}

// ForkProvider is the interface for providing fork information.
type ForkProvider interface {
	// Fork fetches fork information for the given state.
	Fork(ctx context.Context, stateID string) (*phase0.Fork, error)
}

// ForkScheduleProvider is the interface for providing fork schedule data.
type ForkScheduleProvider interface {
	// ForkSchedule provides details of past and future changes in the chain's fork version.
	ForkSchedule(ctx context.Context) ([]*phase0.Fork, error)
}

// GenesisProvider is the interface for providing genesis information.
type GenesisProvider interface {
	// Genesis fetches genesis information for the chain.
	Genesis(ctx context.Context) (*apiv1.Genesis, error)
}

// NodeSyncingProvider is the interface for providing synchronization state.
type NodeSyncingProvider interface {
	// NodeSyncing provides the state of the node's synchronization with the chain.
	NodeSyncing(ctx context.Context) (*apiv1.SyncState, error)
}

// ProposalPreparationsSubmitter is the interface for submitting proposal preparations.
type ProposalPreparationsSubmitter interface {
	// SubmitProposalPreparations provides the beacon node with information required if a proposal for the given validators
	// shows up in the next epoch.
	SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error
}

// ProposerDutiesProvider is the interface for providing proposer duties.
type ProposerDutiesProvider interface {
	// ProposerDuties obtains proposer duties for the given epoch.
	// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
	ProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*apiv1.ProposerDuty, error)
}

// SpecProvider is the interface for providing spec data.
type SpecProvider interface {
	// Spec provides the spec information of the chain.
	Spec(ctx context.Context) (map[string]interface{}, error)
}

// SyncStateProvider is the interface for providing synchronization state.
type SyncStateProvider interface {
	// SyncState provides the state of the node's synchronization with the chain.
	SyncState(ctx context.Context) (*apiv1.SyncState, error)
}

// ValidatorBalancesProvider is the interface for providing validator balances.
type ValidatorBalancesProvider interface {
	// ValidatorBalances provides the validator balances for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
	// will be applied.
	ValidatorBalances(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]phase0.Gwei, error)
}

// ValidatorsProvider is the interface for providing validator information.
type ValidatorsProvider interface {
	// Validators provides the validators, with their balance and status, for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
	// will be applied.
	Validators(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]*apiv1.Validator, error)

	// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
	// supplied no filter will be applied.
	ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*apiv1.Validator, error)
}

// VoluntaryExitSubmitter is the interface for submitting voluntary exits.
type VoluntaryExitSubmitter interface {
	// SubmitVoluntaryExit submits a voluntary exit.
	SubmitVoluntaryExit(ctx context.Context, voluntaryExit *phase0.SignedVoluntaryExit) error
}

//
// Local extensions
//

// DomainProvider provides a domain for a given domain type at an epoch.
type DomainProvider interface {
	// Domain provides a domain for a given domain type at a given epoch.
	Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error)

	// GenesisDomain returns the domain for the given domain type at genesis.
	// N.B. this is not always the same as the the domain at epoch 0.  It is possible
	// for a chain's fork schedule to have multiple forks at genesis.  In this situation,
	// GenesisDomain() will return the first, and Domain() will return the last.
	GenesisDomain(ctx context.Context, domainType phase0.DomainType) (phase0.Domain, error)
}

// GenesisTimeProvider is the interface for providing the genesis time of a chain.
type GenesisTimeProvider interface {
	// GenesisTime provides the genesis time of the chain.
	GenesisTime(ctx context.Context) (time.Time, error)
}

// NodeClientProvider provides the client for the node.
type NodeClientProvider interface {
	// NodeClient provides the client for the node.
	NodeClient(ctx context.Context) (string, error)
}
