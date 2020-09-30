// Copyright © 2020 Attestant Limited.
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

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// Service is the service providing a connection to an Ethereum 2 client.
type Service interface {
	// Name returns the name of the client.
	Name() string

	// EpochFromStateID converts a state ID to its epoch.
	EpochFromStateID(ctx context.Context, stateID string) (uint64, error)

	// SlotFromStateID converts a state ID to its slot.
	SlotFromStateID(ctx context.Context, stateID string) (uint64, error)
}

// NodeVersionProvider is the interface for providing the node version.
type NodeVersionProvider interface {
	// NodeVersion returns a free-text string with the node version.
	NodeVersion(ctx context.Context) (string, error)
}

// GenesisTimeProvider is the interface for providing the genesis time of a chain.
type GenesisTimeProvider interface {
	// GenesisTime provides the genesis time of the chain.
	GenesisTime(ctx context.Context) (time.Time, error)
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
	BeaconAttesterDomain(ctx context.Context) ([]byte, error)
}

// BeaconProposerDomainProvider is the interface for providing the beacon proposer domain.
type BeaconProposerDomainProvider interface {
	// BeaconProposerDomain provides the beacon proposer domain.
	BeaconProposerDomain(ctx context.Context) ([]byte, error)
}

// RANDAODomainProvider is the interface for providing the RANDAO domain.
type RANDAODomainProvider interface {
	// RANDAODomain provides the RANDAO domain.
	RANDAODomain(ctx context.Context) ([]byte, error)
}

// DepositDomainProvider is the interface for providing the deposit domain.
type DepositDomainProvider interface {
	// DepositDomain provides the deposit domain.
	DepositDomain(ctx context.Context) ([]byte, error)
}

// VoluntaryExitDomainProvider is the interface for providing the voluntary exit domain.
type VoluntaryExitDomainProvider interface {
	// VoluntaryExitDomain provides the voluntary exit domain.
	VoluntaryExitDomain(ctx context.Context) ([]byte, error)
}

// SelectionProofDomainProvider is the interface for providing the selection proof domain.
type SelectionProofDomainProvider interface {
	// SelectionProofDomain provides the selection proof domain.
	SelectionProofDomain(ctx context.Context) ([]byte, error)
}

// AggregateAndProofDomainProvider is the interface for providing the aggregate and proof domain.
type AggregateAndProofDomainProvider interface {
	// AggregateAndProofDomain provides the aggregate and proof domain.
	AggregateAndProofDomain(ctx context.Context) ([]byte, error)
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
	Index(ctx context.Context) (uint64, error)
}

// ValidatorPubKeyProvider is the interface for entities that can provide the public key of a validator.
type ValidatorPubKeyProvider interface {
	// PubKey provides the public key of the validator.
	PubKey(ctx context.Context) ([]byte, error)
}

// ValidatorIDProvider is the interface that provides the identifiers (pubkey, index) of a validator.
type ValidatorIDProvider interface {
	ValidatorIndexProvider
	ValidatorPubKeyProvider
}

// DepositContractProvider is the interface for providng details about the deposit contract.
type DepositContractProvider interface {
	// DepositContractAddress provides the Ethereum 1 address of the deposit contract.
	DepositContractAddress(ctx context.Context) ([]byte, error)

	// DepositContractChainID provides the Ethereum 1 chain ID of the deposit contract.
	DepositContractChainID(ctx context.Context) (uint64, error)

	// DepositContractNetworkID provides the Ethereum 1 network ID of the deposit contract.
	DepositContractNetworkID(ctx context.Context) (uint64, error)
}

// NonSpecAggregateAttestationProvider is the interface for providing aggregate attestations.
type NonSpecAggregateAttestationProvider interface {
	// NonSpecAggregateAttestation fetches the aggregate attestation given an attestation.
	NonSpecAggregateAttestation(ctx context.Context, attestation *spec.Attestation, validatorPubKey []byte, slotSignature []byte) (*spec.Attestation, error)
}

// SignatureDomainProvider provides a full signature domain for a given domain at an epoch.
type SignatureDomainProvider interface {
	// SignatureDomain provides a signature domain for a given domain at a given epoch.
	SignatureDomain(ctx context.Context, domain []byte, epoch uint64) ([]byte, error)
}

// SignedBeaconBlockProvider is the interface for providing beacon blocks.
type SignedBeaconBlockProvider interface {
	// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
	SignedBeaconBlockBySlot(ctx context.Context, slot uint64) (*spec.SignedBeaconBlock, error)
}

// BeaconBlockRootProvider is the interface for providing beacon block roots.
type BeaconBlockRootProvider interface {
	// BeaconBlockRootBySlot fetches a block's root given its slot.
	BeaconBlockRootBySlot(ctx context.Context, slot uint64) ([]byte, error)
}

// BeaconCommitteesProvider is the interface for providing beacon committees.
type BeaconCommitteesProvider interface {
	// BeaconCommittees fetches the chain's beacon committees given a state.
	BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error)
}

// ValidatorsWithoutBalanceProvider is the interface for providing validator information, minus the balance.
type ValidatorsWithoutBalanceProvider interface {
	// ValidatorsWithoutBalance provides the validators, with their status, for a given state.
	// Balances are set to 0.
	// This is a non-standard all, only to be used if fetching balances results in the call being too slow.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
	ValidatorsWithoutBalance(ctx context.Context, stateID string, validators []ValidatorIDProvider) (map[uint64]*api.Validator, error)
}

//
// Standard API
//

// GenesisProvider is the interface for providing genesis information.
type GenesisProvider interface {
	// Genesis fetches genesis information for the chain.
	Genesis(ctx context.Context) (*api.Genesis, error)
}

// ForkProvider is the interface for providing fork information.
type ForkProvider interface {
	// Fork fetches fork information for the given state.
	Fork(ctx context.Context, stateID string) (*spec.Fork, error)
}

// ValidatorsProvider is the interface for providing validator information.
type ValidatorsProvider interface {
	// Validators provides the validators, with their balance and status, for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
	Validators(ctx context.Context, stateID string, validators []ValidatorIDProvider) (map[uint64]*api.Validator, error)
}

// ValidatorBalancesProvider is the interface for providing validator balances.
type ValidatorBalancesProvider interface {
	// ValidatorBalances provides the validator balances for a given state.
	// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
	ValidatorBalances(ctx context.Context, stateID string, validators []ValidatorIDProvider) (map[uint64]uint64, error)
}

// BeaconBlockSubmitter is the interface for submitting beacon blocks.
type BeaconBlockSubmitter interface {
	// SubmitBeaconBlock submits a beacon block.
	SubmitBeaconBlock(ctx context.Context, block *spec.SignedBeaconBlock) error
}

// AttestationSubmitter is the interface for submitting attestations.
type AttestationSubmitter interface {
	// SubmitAttestation submits an attestation.
	SubmitAttestation(ctx context.Context, attestation *spec.Attestation) error
}

// AggregateAttestationsSubmitter is the interface for submitting aggregate attestations.
type AggregateAttestationsSubmitter interface {
	// SubmitAggregateAttestations submits aggregate attestations.
	SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error
}

// BeaconCommitteeSubscriptionsSubmitter is the interface for submitting beacon committee subnet subscription requests.
type BeaconCommitteeSubscriptionsSubmitter interface {
	// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
	SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*BeaconCommitteeSubscription) error
}

// SyncStateProvider is the interface for providing synchronization state.
type SyncStateProvider interface {
	// SyncState provides the state of the node's synchronization with the chain.
	SyncState(ctx context.Context) (*api.SyncState, error)
}

// AttesterDutiesProvider is the interface for providing attester duties
type AttesterDutiesProvider interface {
	// AttesterDuties obtains attester duties.
	// If validators is nil it will return all duties for the given epoch.
	AttesterDuties(ctx context.Context, epoch uint64, validators []ValidatorIDProvider) ([]*api.AttesterDuty, error)
}

// ProposerDutiesProvider is the interface for providing proposer duties.
type ProposerDutiesProvider interface {
	// ProposerDuties obtains proposer duties.
	// If validators is nil it will return all duties for the given epoch.
	ProposerDuties(ctx context.Context, epoch uint64, validators []ValidatorIDProvider) ([]*api.ProposerDuty, error)
}

// BeaconBlockProposalProvider is the interface for providing beacon block proposals.
type BeaconBlockProposalProvider interface {
	// BeaconBlockProposal fetches a proposed beacon block for signing.
	BeaconBlockProposal(ctx context.Context, slot uint64, randaoReveal []byte, graffiti []byte) (*spec.BeaconBlock, error)
}

// AttestationDataProvider is the interface for providing attestation data.
type AttestationDataProvider interface {
	// AttestationData fetches the attestation data for the given slot and committee index.
	AttestationData(ctx context.Context, slot uint64, committeeIndex uint64) (*spec.AttestationData, error)
}

// AggregateAttestationProvider is the interface for providing aggregate attestations.
type AggregateAttestationProvider interface {
	// AggregateAttestation fetches the aggregate attestation given an attestation.
	NonSpecAggregateAttestation(ctx context.Context, attestation *spec.Attestation, validatorPubKey []byte, slotSignature []byte) (*spec.Attestation, error)
	AggregateAttestation(ctx context.Context, slot uint64, attestationDataRoot []byte) (*spec.Attestation, error)
}

// type DepositContractProvider interface {
// 	// DepositContract provides details of the Ethereum 1 deposit contract for the chain.
// 	DepositContract(ctx context.Context) (*api.DepositContract, error)
// }

// ForkScheduleProvider is the interface for providing fork schedule data.
type ForkScheduleProvider interface {
	// ForkSchedule provides details of past and future changes in the chain's fork version.
	ForkSchedule(ctx context.Context) ([]*spec.Fork, error)
}

// SpecProvider is the interface for providing spec data.
type SpecProvider interface {
	// Spec provides the spec information of the chain.
	Spec(ctx context.Context) (map[string]interface{}, error)
}
