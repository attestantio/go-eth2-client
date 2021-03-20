// Copyright © 2021 Attestant Limited.
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

package v2

import (
	"context"
	"time"

	api "github.com/attestantio/go-eth2-client/api/v2"
	spec "github.com/attestantio/go-eth2-client/spec/altair"
)

// Service is the service providing a connection to an Ethereum 2 client.
type Service interface {
	// Name returns the name of the client implementation.
	Name() string

	// Address returns the address of the client.
	Address() string
}

// PrysmAttesterDutiesProvider is the interface for providing attester duties with prysm-specific parameters.
type PrysmAttesterDutiesProvider interface {
	// PrysmAttesterDuties obtains attester duties with prysm-specific parameters.
	PrysmAttesterDuties(ctx context.Context, epoch spec.Epoch, validatorPubKeys []spec.BLSPubKey) ([]*api.AttesterDuty, error)
}

// PrysmProposerDutiesProvider is the interface for providing proposer duties with prysm-specific parameters.
type PrysmProposerDutiesProvider interface {
	// PrysmProposerDuties obtains proposer duties with prysm-specific parameters.
	PrysmProposerDuties(ctx context.Context, epoch spec.Epoch, validatorPubKeys []spec.BLSPubKey) ([]*api.ProposerDuty, error)
}

// PrysmValidatorBalancesProvider is the interface for providing validator balances.
type PrysmValidatorBalancesProvider interface {
	// PrysmValidatorBalances provides the validator balances for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIDs is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
	// will be applied.
	PrysmValidatorBalances(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]spec.Gwei, error)
}

// EpochFromStateIDProvider is the interface for providing epochs from state IDs.
type EpochFromStateIDProvider interface {
	// EpochFromStateID converts a state ID to its epoch.
	EpochFromStateID(ctx context.Context, stateID string) (spec.Epoch, error)
}

// SlotFromStateIDProvider is the interface for providing slots from state IDs.
type SlotFromStateIDProvider interface {
	// SlotFromStateID converts a state ID to its slot.
	SlotFromStateID(ctx context.Context, stateID string) (spec.Slot, error)
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
	FarFutureEpoch(ctx context.Context) (spec.Epoch, error)
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

// BeaconAttesterDomainProvider is the interface for providing the beacon attester domain.
type BeaconAttesterDomainProvider interface {
	// BeaconAttesterDomain provides the beacon attester domain.
	BeaconAttesterDomain(ctx context.Context) (spec.DomainType, error)
}

// BeaconProposerDomainProvider is the interface for providing the beacon proposer domain.
type BeaconProposerDomainProvider interface {
	// BeaconProposerDomain provides the beacon proposer domain.
	BeaconProposerDomain(ctx context.Context) (spec.DomainType, error)
}

// RANDAODomainProvider is the interface for providing the RANDAO domain.
type RANDAODomainProvider interface {
	// RANDAODomain provides the RANDAO domain.
	RANDAODomain(ctx context.Context) (spec.DomainType, error)
}

// DepositDomainProvider is the interface for providing the deposit domain.
type DepositDomainProvider interface {
	// DepositDomain provides the deposit domain.
	DepositDomain(ctx context.Context) (spec.DomainType, error)
}

// VoluntaryExitDomainProvider is the interface for providing the voluntary exit domain.
type VoluntaryExitDomainProvider interface {
	// VoluntaryExitDomain provides the voluntary exit domain.
	VoluntaryExitDomain(ctx context.Context) (spec.DomainType, error)
}

// SelectionProofDomainProvider is the interface for providing the selection proof domain.
type SelectionProofDomainProvider interface {
	// SelectionProofDomain provides the selection proof domain.
	SelectionProofDomain(ctx context.Context) (spec.DomainType, error)
}

// AggregateAndProofDomainProvider is the interface for providing the aggregate and proof domain.
type AggregateAndProofDomainProvider interface {
	// AggregateAndProofDomain provides the aggregate and proof domain.
	AggregateAndProofDomain(ctx context.Context) (spec.DomainType, error)
}

// SyncCommitteeSelectionProofDomainProvider is the interface for providing the sync committee selection proof domain.
type SyncCommitteeSelectionProofDomainProvider interface {
	// SyncCommitteeSelectionProofDomain provides the sync committee selection proof domain.
	SyncCommitteeSelectionProofDomain(ctx context.Context) (spec.DomainType, error)
}

// ContributionAndProofDomainProvider is the interface for providing the contribution and proof domain.
type ContributionAndProofDomainProvider interface {
	// ContributionAndProofDomain provides the contribution and proof domain.
	ContributionAndProofDomain(ctx context.Context) (spec.DomainType, error)
}

// SyncCommitteeDomainProvider is the interface for providing the sync committee domain.
type SyncCommitteeDomainProvider interface {
	// SyncCommitteeDomain provides the sync committee domain.
	SyncCommitteeDomain(ctx context.Context) (spec.DomainType, error)
}

// BeaconChainHeadUpdatedSource is the interface for a service that provides beacon chain head updates.
type BeaconChainHeadUpdatedSource interface {
	// AddOnBeaconChainHeadUpdatedHandler adds a handler provided with beacon chain head updates.
	AddOnBeaconChainHeadUpdatedHandler(ctx context.Context, handler BeaconChainHeadUpdatedHandler) error
}

// BeaconChainHeadUpdatedHandler is the interface that needs to be implemented by processes that wish
// to receive beacon chain head updated messages.
type BeaconChainHeadUpdatedHandler interface {
	// OnBeaconChainHeadUpdated is called whenever we receive a notification of an update to the beacon chain head.
	OnBeaconChainHeadUpdated(ctx context.Context, slot uint64, blockRoot []byte, stateRoot []byte, epochTransition bool)
}

// ValidatorIndexProvider is the interface for entities that can provide the index of a validator.
type ValidatorIndexProvider interface {
	// Index provides the index of the validator.
	Index(ctx context.Context) (spec.ValidatorIndex, error)
}

// ValidatorPubKeyProvider is the interface for entities that can provide the public key of a validator.
type ValidatorPubKeyProvider interface {
	// PubKey provides the public key of the validator.
	PubKey(ctx context.Context) (spec.BLSPubKey, error)
}

// ValidatorIDProvider is the interface that provides the identifiers (pubkey, index) of a validator.
type ValidatorIDProvider interface {
	ValidatorIndexProvider
	ValidatorPubKeyProvider
}

// DepositContractProvider is the interface for providing details about the deposit contract.
type DepositContractProvider interface {
	// DepositContract provides details of the Ethereum 1 deposit contract for the chain.
	DepositContract(ctx context.Context) (*api.DepositContract, error)
}

// PrysmAggregateAttestationProvider is the interface for providing aggregate attestations.
type PrysmAggregateAttestationProvider interface {
	// PrysmAggregateAttestation fetches the aggregate attestation given an attestation.
	PrysmAggregateAttestation(ctx context.Context, attestation *spec.Attestation, validatorPubKey spec.BLSPubKey, slotSignature spec.BLSSignature) (*spec.Attestation, error)
}

// SignedBeaconBlockProvider is the interface for providing beacon blocks.
type SignedBeaconBlockProvider interface {
	// SignedBeaconBlock fetches a signed beacon block given a block ID.
	SignedBeaconBlock(ctx context.Context, blockID string) (*spec.SignedBeaconBlock, error)
	// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
	// SignedBeaconBlockBySlot(ctx context.Context, slot uint64) (*spec.SignedBeaconBlock, error)
}

// BeaconBlockRootProvider is the interface for providing beacon block roots.
type BeaconBlockRootProvider interface {
	// BeaconBlockRootBySlot fetches a block's root given its slot.
	BeaconBlockRootBySlot(ctx context.Context, slot uint64) ([]byte, error)
}

// BeaconCommitteesProvider is the interface for providing beacon committees.
type BeaconCommitteesProvider interface {
	// BeaconCommittees fetches all beacon committees for the epoch at the given state.
	BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error)
}

// ValidatorsWithoutBalanceProvider is the interface for providing validator information, minus the balance.
type ValidatorsWithoutBalanceProvider interface {
	// ValidatorsWithoutBalance provides the validators, with their status, for a given state.
	// Balances are set to 0.
	// This is a non-standard call, only to be used if fetching balances results in the call being too slow.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
	// will be applied.
	ValidatorsWithoutBalance(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error)

	// ValidatorsWithoutBalanceByPubKey provides the validators, with their status, for a given state.
	// This is a non-standard call, only to be used if fetching balances results in the call being too slow.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
	// supplied no filter will be applied.
	ValidatorsWithoutBalanceByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error)
}

// EventHandlerFunc is the handler for events.
type EventHandlerFunc func(*api.Event)

//
// Standard API
//

// AggregateAttestationProvider is the interface for providing aggregate attestations.
type AggregateAttestationProvider interface {
	// AggregateAttestation fetches the aggregate attestation given an attestation.
	AggregateAttestation(ctx context.Context, slot spec.Slot, attestationDataRoot spec.Root) (*spec.Attestation, error)
}

// AggregateAttestationsSubmitter is the interface for submitting aggregate attestations.
type AggregateAttestationsSubmitter interface {
	// SubmitAggregateAttestations submits aggregate attestations.
	SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error
}

// AttestationDataProvider is the interface for providing attestation data.
type AttestationDataProvider interface {
	// AttestationData fetches the attestation data for the given slot and committee index.
	AttestationData(ctx context.Context, slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error)
}

// AttestationPoolProvider is the interface for providing attestation pools.
type AttestationPoolProvider interface {
	// AttestationPool fetches the attestation pool for the given slot.
	AttestationPool(ctx context.Context, slot spec.Slot) ([]*spec.Attestation, error)
}

// AttestationsSubmitter is the interface for submitting attestations.
type AttestationsSubmitter interface {
	// SubmitAttestations submits attestations.
	SubmitAttestations(ctx context.Context, attestations []*spec.Attestation) error
}

// AttesterDutiesProvider is the interface for providing attester duties
type AttesterDutiesProvider interface {
	// AttesterDuties obtains attester duties.
	// If validatorIndicess is nil it will return all duties for the given epoch.
	AttesterDuties(ctx context.Context, epoch spec.Epoch, validatorIndices []spec.ValidatorIndex) ([]*api.AttesterDuty, error)
}

// BeaconBlockHeadersProvider is the interface for providing beacon block headers.
type BeaconBlockHeadersProvider interface {
	// BeaconBlockHeader provides the block header of a given block ID.
	BeaconBlockHeader(ctx context.Context, blockID string) (*api.BeaconBlockHeader, error)
}

// BeaconBlockProposalProvider is the interface for providing beacon block proposals.
type BeaconBlockProposalProvider interface {
	// BeaconBlockProposal fetches a proposed beacon block for signing.
	BeaconBlockProposal(ctx context.Context, slot spec.Slot, randaoReveal spec.BLSSignature, graffiti []byte) (*spec.BeaconBlock, error)
}

// BeaconBlockSubmitter is the interface for submitting beacon blocks.
type BeaconBlockSubmitter interface {
	// SubmitBeaconBlock submits a beacon block.
	SubmitBeaconBlock(ctx context.Context, block *spec.SignedBeaconBlock) error
}

// BeaconCommitteeSubscriptionsSubmitter is the interface for submitting beacon committee subnet subscription requests.
type BeaconCommitteeSubscriptionsSubmitter interface {
	// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
	SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*api.BeaconCommitteeSubscription) error
}

// BeaconStateProvider is the interface for providing beacon state.
type BeaconStateProvider interface {
	// BeaconState fetches a beacon state.
	BeaconState(ctx context.Context, stateID string) (*spec.BeaconState, error)
}

// EventsProvider is the interface for providing events.
type EventsProvider interface {
	// Events feeds requested events with the given topics to the supplied handler.
	Events(ctx context.Context, topics []string, handler EventHandlerFunc) error
}

// FinalityProvider is the interface for providing finality information.
type FinalityProvider interface {
	// Finality provides the finality given a state ID.
	Finality(ctx context.Context, stateID string) (*api.Finality, error)
}

// ForkProvider is the interface for providing fork information.
type ForkProvider interface {
	// Fork fetches fork information for the given state.
	Fork(ctx context.Context, stateID string) (*spec.Fork, error)
}

// ForkScheduleProvider is the interface for providing fork schedule data.
type ForkScheduleProvider interface {
	// ForkSchedule provides details of past and future changes in the chain's fork version.
	ForkSchedule(ctx context.Context) ([]*spec.Fork, error)
}

// GenesisProvider is the interface for providing genesis information.
type GenesisProvider interface {
	// Genesis fetches genesis information for the chain.
	Genesis(ctx context.Context) (*api.Genesis, error)
}

// NodeSyncingProvider is the interface for providing synchronization state.
type NodeSyncingProvider interface {
	// NodeSyncing provides the state of the node's synchronization with the chain.
	NodeSyncing(ctx context.Context) (*api.SyncState, error)
}

// ProposerDutiesProvider is the interface for providing proposer duties.
type ProposerDutiesProvider interface {
	// ProposerDuties obtains proposer duties for the given epoch.
	// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
	ProposerDuties(ctx context.Context, epoch spec.Epoch, validatorIndices []spec.ValidatorIndex) ([]*api.ProposerDuty, error)
}

// SpecProvider is the interface for providing spec data.
type SpecProvider interface {
	// Spec provides the spec information of the chain.
	Spec(ctx context.Context) (map[string]interface{}, error)
}

// SyncStateProvider is the interface for providing synchronization state.
type SyncStateProvider interface {
	// SyncState provides the state of the node's synchronization with the chain.
	SyncState(ctx context.Context) (*api.SyncState, error)
}

// ValidatorBalancesProvider is the interface for providing validator balances.
type ValidatorBalancesProvider interface {
	// ValidatorBalances provides the validator balances for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
	// will be applied.
	ValidatorBalances(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]spec.Gwei, error)
}

// ValidatorsProvider is the interface for providing validator information.
type ValidatorsProvider interface {
	// Validators provides the validators, with their balance and status, for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
	// will be applied.
	Validators(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error)

	// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
	// supplied no filter will be applied.
	ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error)
}

// VoluntaryExitSubmitter is the interface for submitting voluntary exits.
type VoluntaryExitSubmitter interface {
	// SubmitVoluntaryExit submits a voluntary exit.
	SubmitVoluntaryExit(ctx context.Context, voluntaryExit *spec.SignedVoluntaryExit) error
}

//
// Local extensions
//

// DomainProvider provides a domain for a given domain type at an epoch.
type DomainProvider interface {
	// Domain provides a domain for a given domain type at a given epoch.
	Domain(ctx context.Context, domainType spec.DomainType, epoch spec.Epoch) (spec.Domain, error)
}

// GenesisTimeProvider is the interface for providing the genesis time of a chain.
type GenesisTimeProvider interface {
	// GenesisTime provides the genesis time of the chain.
	GenesisTime(ctx context.Context) (time.Time, error)
}
