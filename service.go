// Copyright © 2020 - 2023 Attestant Limited.
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
	//
	// Deprecated: will be removed in a future release.
	EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error)
}

// SlotFromStateIDProvider is the interface for providing slots from state IDs.
type SlotFromStateIDProvider interface {
	// SlotFromStateID converts a state ID to its slot.
	//
	// Deprecated: will be removed in a future release.
	SlotFromStateID(ctx context.Context, stateID string) (phase0.Slot, error)
}

// SlotDurationProvider is the interface for providing the duration of each slot of a chain.
type SlotDurationProvider interface {
	// SlotDuration provides the duration of a slot of the chain.
	//
	// Deprecated: use Spec()
	SlotDuration(ctx context.Context) (time.Duration, error)
}

// SlotsPerEpochProvider is the interface for providing the number of slots in each epoch of a chain.
type SlotsPerEpochProvider interface {
	// SlotsPerEpoch provides the slots per epoch of the chain.
	//
	// Deprecated: use Spec()
	SlotsPerEpoch(ctx context.Context) (uint64, error)
}

// FarFutureEpochProvider is the interface for providing the far future epoch of a chain.
type FarFutureEpochProvider interface {
	// FarFutureEpoch provides the far future epoch of the chain.
	FarFutureEpoch(ctx context.Context) (phase0.Epoch, error)
}

// TargetAggregatorsPerCommitteeProvider is the interface for providing the target number of
// aggregators in each attestation committee.
type TargetAggregatorsPerCommitteeProvider interface {
	// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
	//
	// Deprecated: use Spec()
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

// SignedBeaconBlockProvider is the interface for providing beacon blocks.
type SignedBeaconBlockProvider interface {
	// SignedBeaconBlock fetches a signed beacon block given a block ID.
	SignedBeaconBlock(ctx context.Context, opts *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error)
}

// BlobSidecarsProvider is the interface for providing blobs for a given beacon block.
type BlobSidecarsProvider interface {
	// BlobSidecars fetches the blobs given a block ID.
	BlobSidecars(ctx context.Context, opts *api.BlobSidecarsOpts) (*api.Response[[]*deneb.BlobSidecar], error)
}

// BeaconCommitteesProvider is the interface for providing beacon committees.
type BeaconCommitteesProvider interface {
	// BeaconCommittees fetches all beacon committees for the given options.
	BeaconCommittees(ctx context.Context, opts *api.BeaconCommitteesOpts) (*api.Response[[]*apiv1.BeaconCommittee], error)
}

// SyncCommitteesProvider is the interface for providing sync committees.
type SyncCommitteesProvider interface {
	// SyncCommittee fetches the sync committee for the given state.
	SyncCommittee(ctx context.Context, opts *api.SyncCommitteeOpts) (*api.Response[*apiv1.SyncCommittee], error)
}

// EventHandlerFunc is the handler for events.
type EventHandlerFunc func(*apiv1.Event)

//
// Standard API
//

// AggregateAttestationProvider is the interface for providing aggregate attestations.
type AggregateAttestationProvider interface {
	// AggregateAttestation fetches the aggregate attestation for the given options.
	AggregateAttestation(ctx context.Context, opts *api.AggregateAttestationOpts) (*api.Response[*phase0.Attestation], error)
}

// AggregateAttestationsSubmitter is the interface for submitting aggregate attestations.
type AggregateAttestationsSubmitter interface {
	// SubmitAggregateAttestations submits aggregate attestations.
	SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error
}

// AttestationDataProvider is the interface for providing attestation data.
type AttestationDataProvider interface {
	// AttestationData fetches the attestation data for the given options.
	AttestationData(ctx context.Context, opts *api.AttestationDataOpts) (*api.Response[*phase0.AttestationData], error)
}

// AttestationPoolProvider is the interface for providing attestation pools.
type AttestationPoolProvider interface {
	// AttestationPool fetches the attestation pool for the given options.
	AttestationPool(ctx context.Context, opts *api.AttestationPoolOpts) (*api.Response[[]*phase0.Attestation], error)
}

// AttestationsSubmitter is the interface for submitting attestations.
type AttestationsSubmitter interface {
	// SubmitAttestations submits attestations.
	SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error
}

// AttesterSlashingSubmitter is the interface for submitting attester slashings.
type AttesterSlashingSubmitter interface {
	// SubmitAttesterSlashing submits an attester slashing
	SubmitAttesterSlashing(ctx context.Context, slashing *phase0.AttesterSlashing) error
}

// AttesterDutiesProvider is the interface for providing attester duties.
type AttesterDutiesProvider interface {
	// AttesterDuties obtains attester duties.
	AttesterDuties(ctx context.Context, opts *api.AttesterDutiesOpts) (*api.Response[[]*apiv1.AttesterDuty], error)
}

// DepositContractProvider is the interface for providing details about the deposit contract.
type DepositContractProvider interface {
	// DepositContract provides details of the execution deposit contract for the chain.
	DepositContract(ctx context.Context, opts *api.DepositContractOpts) (*api.Response[*apiv1.DepositContract], error)
}

// SyncCommitteeDutiesProvider is the interface for providing sync committee duties.
type SyncCommitteeDutiesProvider interface {
	// SyncCommitteeDuties obtains sync committee duties.
	// If validatorIndicess is nil it will return all duties for the given epoch.
	SyncCommitteeDuties(ctx context.Context, opts *api.SyncCommitteeDutiesOpts) (*api.Response[[]*apiv1.SyncCommitteeDuty], error)
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
	SyncCommitteeContribution(ctx context.Context, opts *api.SyncCommitteeContributionOpts) (*api.Response[*altair.SyncCommitteeContribution], error)
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
	BeaconBlockHeader(ctx context.Context, opts *api.BeaconBlockHeaderOpts) (*api.Response[*apiv1.BeaconBlockHeader], error)
}

// ProposalProvider is the interface for providing proposals.
type ProposalProvider interface {
	// Proposal fetches a proposal for signing.
	Proposal(ctx context.Context, opts *api.ProposalOpts) (*api.Response[*api.VersionedProposal], error)
}

// ProposalSlashingSubmitter is the interface for submitting proposal slashings.
type ProposalSlashingSubmitter interface {
	SubmitProposalSlashing(ctx context.Context, slashing *phase0.ProposerSlashing) error
}

// BeaconBlockRootProvider is the interface for providing beacon block roots.
type BeaconBlockRootProvider interface {
	// BeaconBlockRoot fetches a block's root given a set of options.
	BeaconBlockRoot(ctx context.Context, opts *api.BeaconBlockRootOpts) (*api.Response[*phase0.Root], error)
}

// BeaconBlockSubmitter is the interface for submitting beacon blocks.
type BeaconBlockSubmitter interface {
	// SubmitBeaconBlock submits a beacon block.
	//
	// Deprecated: this will not work from the deneb hard-fork onwards.  Use ProposalSubmitter.SubmitProposal() instead.
	SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error
}

// ProposalSubmitter is the interface for submitting proposals.
type ProposalSubmitter interface {
	// SubmitProposal submits a proposal.
	SubmitProposal(ctx context.Context, block *api.VersionedSignedProposal) error
}

// BeaconCommitteeSubscriptionsSubmitter is the interface for submitting beacon committee subnet subscription requests.
type BeaconCommitteeSubscriptionsSubmitter interface {
	// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
	SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error
}

// BeaconStateProvider is the interface for providing beacon state.
type BeaconStateProvider interface {
	// BeaconState fetches a beacon state given a state ID.
	BeaconState(ctx context.Context, opts *api.BeaconStateOpts) (*api.Response[*spec.VersionedBeaconState], error)
}

// BeaconStateRandaoProvider is the interface for providing beacon state RANDAOs.
type BeaconStateRandaoProvider interface {
	// BeaconStateRandao fetches a beacon state RANDAO given a state ID.
	BeaconStateRandao(ctx context.Context, opts *api.BeaconStateRandaoOpts) (*api.Response[*phase0.Root], error)
}

// BeaconStateRootProvider is the interface for providing beacon state roots.
type BeaconStateRootProvider interface {
	// BeaconStateRoot fetches a beacon state root given a state ID.
	BeaconStateRoot(ctx context.Context, opts *api.BeaconStateRootOpts) (*api.Response[*phase0.Root], error)
}

// BlindedProposalProvider is the interface for providing blinded beacon block proposals.
type BlindedProposalProvider interface {
	// BlindedProposal fetches a blinded proposed beacon block for signing.
	BlindedProposal(ctx context.Context, opts *api.BlindedProposalOpts) (*api.Response[*api.VersionedBlindedProposal], error)
}

// BlindedBeaconBlockSubmitter is the interface for submitting blinded beacon blocks.
type BlindedBeaconBlockSubmitter interface {
	// SubmitBlindedBeaconBlock submits a beacon block.
	//
	// Deprecated: this will not work from the deneb hard-fork onwards.  Use BlindedProposalSubmitter.SubmitBlindedProposal() instead.
	SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error
}

// BlindedProposalSubmitter is the interface for submitting blinded proposals.
type BlindedProposalSubmitter interface {
	// SubmitBlindedProposal submits a beacon block.
	SubmitBlindedProposal(ctx context.Context, block *api.VersionedSignedBlindedProposal) error
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
	Finality(ctx context.Context, opts *api.FinalityOpts) (*api.Response[*apiv1.Finality], error)
}

// ForkChoiceProvider is the interface for providing fork choice information.
type ForkChoiceProvider interface {
	// Fork fetches all current fork choice context.
	ForkChoice(ctx context.Context, opts *api.ForkChoiceOpts) (*api.Response[*apiv1.ForkChoice], error)
}

// ForkProvider is the interface for providing fork information.
type ForkProvider interface {
	// Fork fetches fork information for the given state.
	Fork(ctx context.Context, opts *api.ForkOpts) (*api.Response[*phase0.Fork], error)
}

// ForkScheduleProvider is the interface for providing fork schedule data.
type ForkScheduleProvider interface {
	// ForkSchedule provides details of past and future changes in the chain's fork version.
	ForkSchedule(ctx context.Context, opts *api.ForkScheduleOpts) (*api.Response[[]*phase0.Fork], error)
}

// GenesisProvider is the interface for providing genesis information.
type GenesisProvider interface {
	// Genesis fetches genesis information for the chain.
	Genesis(ctx context.Context, opts *api.GenesisOpts) (*api.Response[*apiv1.Genesis], error)
}

// NodePeersProvider is the interface for providing peer information.
type NodePeersProvider interface {
	// NodePeers provides the peers of the node.
	NodePeers(ctx context.Context, opts *api.NodePeersOpts) (*api.Response[[]*apiv1.Peer], error)
}

// NodeSyncingProvider is the interface for providing synchronization state.
type NodeSyncingProvider interface {
	// NodeSyncing provides the state of the node's synchronization with the chain.
	NodeSyncing(ctx context.Context, opts *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error)
}

// NodeVersionProvider is the interface for providing the node version.
type NodeVersionProvider interface {
	// NodeVersion returns a free-text string with the node version.
	NodeVersion(ctx context.Context, opts *api.NodeVersionOpts) (*api.Response[string], error)
}

// ProposalPreparationsSubmitter is the interface for submitting proposal preparations.
type ProposalPreparationsSubmitter interface {
	// SubmitProposalPreparations provides the beacon node with information required if a proposal for the given validators
	// shows up in the next epoch.
	SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error
}

// ProposerDutiesProvider is the interface for providing proposer duties.
type ProposerDutiesProvider interface {
	// ProposerDuties obtains proposer duties for the given options.
	ProposerDuties(ctx context.Context, opts *api.ProposerDutiesOpts) (*api.Response[[]*apiv1.ProposerDuty], error)
}

// SpecProvider is the interface for providing spec data.
type SpecProvider interface {
	// Spec provides the spec information of the chain.
	Spec(ctx context.Context, opts *api.SpecOpts) (*api.Response[map[string]any], error)
}

// SyncStateProvider is the interface for providing synchronization state.
type SyncStateProvider interface {
	// SyncState provides the state of the node's synchronization with the chain.
	//
	// Deprecated: use NodeSyncing()
	SyncState(ctx context.Context) (*apiv1.SyncState, error)
}

// ValidatorBalancesProvider is the interface for providing validator balances.
type ValidatorBalancesProvider interface {
	// ValidatorBalances provides the validator balances for the given options.
	ValidatorBalances(ctx context.Context, opts *api.ValidatorBalancesOpts) (*api.Response[map[phase0.ValidatorIndex]phase0.Gwei], error)
}

// ValidatorsProvider is the interface for providing validator information.
type ValidatorsProvider interface {
	// Validators provides the validators, with their balance and status, for the given options.
	Validators(ctx context.Context, opts *api.ValidatorsOpts) (*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator], error)
}

// VoluntaryExitSubmitter is the interface for submitting voluntary exits.
type VoluntaryExitSubmitter interface {
	// SubmitVoluntaryExit submits a voluntary exit.
	SubmitVoluntaryExit(ctx context.Context, voluntaryExit *phase0.SignedVoluntaryExit) error
}

// VoluntaryExitPoolProvider is the interface for providing voluntary exit pools.
type VoluntaryExitPoolProvider interface {
	// VoluntaryExitPool fetches the voluntary exit pool.
	VoluntaryExitPool(ctx context.Context, opts *api.VoluntaryExitPoolOpts) (*api.Response[[]*phase0.SignedVoluntaryExit], error)
}

//
// Local extensions
//

// DomainProvider provides a domain for a given domain type at an epoch.
type DomainProvider interface {
	// Domain provides a domain for a given domain type at a given epoch.
	Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error)

	// GenesisDomain returns the domain for the given domain type at genesis.
	// N.B. this is not always the same as the domain at epoch 0.  It is possible
	// for a chain's fork schedule to have multiple forks at genesis.  In this situation,
	// GenesisDomain() will return the first, and Domain() will return the last.
	GenesisDomain(ctx context.Context, domainType phase0.DomainType) (phase0.Domain, error)
}

// GenesisTimeProvider is the interface for providing the genesis time of a chain.
//
// Deprecated: use Genesis().
type GenesisTimeProvider interface {
	// GenesisTime provides the genesis time of the chain.
	GenesisTime(ctx context.Context) (time.Time, error)
}

// NodeClientProvider provides the client for the node.
type NodeClientProvider interface {
	// NodeClient provides the client for the node.
	NodeClient(ctx context.Context) (*api.Response[string], error)
}
