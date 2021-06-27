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

package testclients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// Erroring is an Ethereum 2 client that errors at a given rate.
type Erroring struct {
	errorRate float64
	next      eth2client.Service
}

// NewErroring creates a new Ethereum 2 client that errors at a given rate.
func NewErroring(ctx context.Context,
	errorRate float64,
	next eth2client.Service,
) (eth2client.Service, error) {
	if next == nil {
		return nil, errors.New("no next service supplied")
	}
	if errorRate < 0 {
		return nil, errors.New("error rate cannot be less than 0")
	}
	if errorRate > 1 {
		return nil, errors.New("error rate cannot be more than 1")
	}

	rand.Seed(time.Now().UnixNano())

	return &Erroring{
		errorRate: errorRate,
		next:      next,
	}, nil
}

// Name returns the name of the client implementation.
func (s *Erroring) Name() string {
	nextName := s.next.Name()
	return fmt.Sprintf("erroring(%v,%s)", s.errorRate, nextName)
}

// Address returns the address of the client.
func (s *Erroring) Address() string {
	nextAddress := s.next.Address()
	return fmt.Sprintf("erroring:%v,%s", s.errorRate, nextAddress)
}

// maybeError may return an error depending on the error arte.
func (s *Erroring) maybeError(ctx context.Context) error {
	// #nosec G404
	roll := rand.Float64()
	if roll < s.errorRate {
		return errors.New("error")
	}
	return nil
}

// PrysmAttesterDuties obtains attester duties with prysm-specific parameters.
func (s *Erroring) PrysmAttesterDuties(ctx context.Context, epoch spec.Epoch, validatorPubKeys []spec.BLSPubKey) ([]*api.AttesterDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmAttesterDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.PrysmAttesterDuties(ctx, epoch, validatorPubKeys)
}

// PrysmProposerDuties obtains proposer duties with prysm-specific parameters.
func (s *Erroring) PrysmProposerDuties(ctx context.Context, epoch spec.Epoch, validatorPubKeys []spec.BLSPubKey) ([]*api.ProposerDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmProposerDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.PrysmProposerDuties(ctx, epoch, validatorPubKeys)
}

// PrysmValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIDs is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Erroring) PrysmValidatorBalances(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]spec.Gwei, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmValidatorBalancesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.PrysmValidatorBalances(ctx, stateID, validatorPubKeys)
}

// EpochFromStateID converts a state ID to its epoch.
func (s *Erroring) EpochFromStateID(ctx context.Context, stateID string) (spec.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.EpochFromStateIDProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
func (s *Erroring) SlotFromStateID(ctx context.Context, stateID string) (spec.Slot, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.SlotFromStateIDProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SlotFromStateID(ctx, stateID)
}

// NodeVersion returns a free-text string with the node version.
func (s *Erroring) NodeVersion(ctx context.Context) (string, error) {
	if err := s.maybeError(ctx); err != nil {
		return "", err
	}
	next, isNext := s.next.(eth2client.NodeVersionProvider)
	if !isNext {
		return "", fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.NodeVersion(ctx)
}

// SlotDuration provides the duration of a slot of the chain.
func (s *Erroring) SlotDuration(ctx context.Context) (time.Duration, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.SlotDurationProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SlotDuration(ctx)
}

// SlotsPerEpoch provides the slots per epoch of the chain.
func (s *Erroring) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.SlotsPerEpochProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Erroring) FarFutureEpoch(ctx context.Context) (spec.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.FarFutureEpochProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.FarFutureEpoch(ctx)
}

// GenesisValidatorsRoot provides the genesis validators root of the chain.
func (s *Erroring) GenesisValidatorsRoot(ctx context.Context) ([]byte, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.GenesisValidatorsRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.GenesisValidatorsRoot(ctx)
}

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
func (s *Erroring) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.TargetAggregatorsPerCommitteeProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.TargetAggregatorsPerCommittee(ctx)
}

// BeaconAttesterDomain provides the beacon attester domain.
func (s *Erroring) BeaconAttesterDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.BeaconAttesterDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconAttesterDomain(ctx)
}

// BeaconProposerDomain provides the beacon proposer domain.
func (s *Erroring) BeaconProposerDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.BeaconProposerDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconProposerDomain(ctx)
}

// RANDAODomain provides the RANDAO domain.
func (s *Erroring) RANDAODomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.RANDAODomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.RANDAODomain(ctx)
}

// DepositDomain provides the deposit domain.
func (s *Erroring) DepositDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.DepositDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.DepositDomain(ctx)
}

// VoluntaryExitDomain provides the voluntary exit domain.
func (s *Erroring) VoluntaryExitDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.VoluntaryExitDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.VoluntaryExitDomain(ctx)
}

// SelectionProofDomain provides the selection proof domain.
func (s *Erroring) SelectionProofDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.SelectionProofDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SelectionProofDomain(ctx)
}

// AggregateAndProofDomain provides the aggregate and proof domain.
func (s *Erroring) AggregateAndProofDomain(ctx context.Context) (spec.DomainType, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.DomainType{}, err
	}
	next, isNext := s.next.(eth2client.AggregateAndProofDomainProvider)
	if !isNext {
		return spec.DomainType{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AggregateAndProofDomain(ctx)
}

// AggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Erroring) AggregateAttestation(ctx context.Context, slot spec.Slot, attestationDataRoot spec.Root) (*spec.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AggregateAttestationProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AggregateAttestation(ctx, slot, attestationDataRoot)
}

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Erroring) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.AggregateAttestationsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Erroring) AttestationData(ctx context.Context, slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttestationDataProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AttestationData(ctx, slot, committeeIndex)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Erroring) AttestationPool(ctx context.Context, slot spec.Slot) ([]*spec.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttestationPoolProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AttestationPool(ctx, slot)
}

// SubmitAttestations submits attestations.
func (s *Erroring) SubmitAttestations(ctx context.Context, attestations []*spec.Attestation) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.AttestationsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitAttestations(ctx, attestations)
}

// AttesterDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Erroring) AttesterDuties(ctx context.Context, epoch spec.Epoch, validatorIndices []spec.ValidatorIndex) ([]*api.AttesterDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttesterDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AttesterDuties(ctx, epoch, validatorIndices)
}

// BeaconBlockHeader provides the block header of a given block ID.
func (s *Erroring) BeaconBlockHeader(ctx context.Context, blockID string) (*api.BeaconBlockHeader, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconBlockHeadersProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconBlockHeader(ctx, blockID)
}

// BeaconCommittees fetches all beacon committees for the epoch at the given state.
func (s *Erroring) BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconCommitteesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconCommittees(ctx, stateID)
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Erroring) BeaconBlockProposal(ctx context.Context, slot spec.Slot, randaoReveal spec.BLSSignature, graffiti []byte) (*spec.BeaconBlock, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconBlockProposalProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

// SubmitBeaconBlock submits a beacon block.
func (s *Erroring) SubmitBeaconBlock(ctx context.Context, block *spec.SignedBeaconBlock) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.BeaconBlockSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitBeaconBlock(ctx, block)
}

// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
func (s *Erroring) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*api.BeaconCommitteeSubscription) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.BeaconCommitteeSubscriptionsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

// BeaconState fetches a beacon state.
func (s *Erroring) BeaconState(ctx context.Context, stateID string) (*spec.BeaconState, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconStateProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconState(ctx, stateID)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Erroring) Events(ctx context.Context, topics []string, handler eth2client.EventHandlerFunc) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.EventsProvider)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Events(ctx, topics, handler)
}

// Finality provides the finality given a state ID.
func (s *Erroring) Finality(ctx context.Context, stateID string) (*api.Finality, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.FinalityProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Finality(ctx, stateID)
}

// Fork fetches fork information for the given state.
func (s *Erroring) Fork(ctx context.Context, stateID string) (*spec.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ForkProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Fork(ctx, stateID)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Erroring) ForkSchedule(ctx context.Context) ([]*spec.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ForkScheduleProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ForkSchedule(ctx)
}

// Genesis fetches genesis information for the chain.
func (s *Erroring) Genesis(ctx context.Context) (*api.Genesis, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.GenesisProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Genesis(ctx)
}

// NodeSyncing provides the state of the node's synchronization with the chain.
func (s *Erroring) NodeSyncing(ctx context.Context) (*api.SyncState, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.NodeSyncingProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.NodeSyncing(ctx)
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
func (s *Erroring) ProposerDuties(ctx context.Context, epoch spec.Epoch, validatorIndices []spec.ValidatorIndex) ([]*api.ProposerDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ProposerDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ProposerDuties(ctx, epoch, validatorIndices)
}

// Spec provides the spec information of the chain.
func (s *Erroring) Spec(ctx context.Context) (map[string]interface{}, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.SpecProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Spec(ctx)
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Erroring) ValidatorBalances(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]spec.Gwei, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ValidatorBalancesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ValidatorBalances(ctx, stateID, validatorIndices)
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
// will be applied.
func (s *Erroring) Validators(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ValidatorsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Validators(ctx, stateID, validatorIndices)
}

// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
// supplied no filter will be applied.
func (s *Erroring) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ValidatorsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ValidatorsByPubKey(ctx, stateID, validatorPubKeys)
}

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Erroring) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *spec.SignedVoluntaryExit) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.VoluntaryExitSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitVoluntaryExit(ctx, voluntaryExit)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Erroring) Domain(ctx context.Context, domainType spec.DomainType, epoch spec.Epoch) (spec.Domain, error) {
	if err := s.maybeError(ctx); err != nil {
		return spec.Domain{}, err
	}
	next, isNext := s.next.(eth2client.DomainProvider)
	if !isNext {
		return spec.Domain{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Domain(ctx, domainType, epoch)
}

// GenesisTime provides the genesis time of the chain.
func (s *Erroring) GenesisTime(ctx context.Context) (time.Time, error) {
	if err := s.maybeError(ctx); err != nil {
		return time.Time{}, err
	}
	next, isNext := s.next.(eth2client.GenesisTimeProvider)
	if !isNext {
		return time.Time{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.GenesisTime(ctx)
}

// DepositContract provides details of the Ethereum 1 deposit contract for the chain.
func (s *Erroring) DepositContract(ctx context.Context) (*api.DepositContract, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.DepositContractProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.DepositContract(ctx)
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Erroring) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.SignedBeaconBlock, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.SignedBeaconBlockProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SignedBeaconBlock(ctx, blockID)
}

// BeaconStateRoot fetches a beacon state root given a state ID.
func (s *Erroring) BeaconStateRoot(ctx context.Context, stateID string) (*spec.Root, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconStateRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconStateRoot(ctx, stateID)
}
