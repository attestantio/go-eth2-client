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

	consensusclient "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Erroring is an Ethereum 2 client that errors at a given rate.
type Erroring struct {
	errorRate float64
	next      consensusclient.Service
}

// NewErroring creates a new Ethereum 2 client that errors at a given rate.
func NewErroring(ctx context.Context,
	errorRate float64,
	next consensusclient.Service,
) (consensusclient.Service, error) {
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

// EpochFromStateID converts a state ID to its epoch.
func (s *Erroring) EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(consensusclient.EpochFromStateIDProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
func (s *Erroring) SlotFromStateID(ctx context.Context, stateID string) (phase0.Slot, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(consensusclient.SlotFromStateIDProvider)
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
	next, isNext := s.next.(consensusclient.NodeVersionProvider)
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
	next, isNext := s.next.(consensusclient.SlotDurationProvider)
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
	next, isNext := s.next.(consensusclient.SlotsPerEpochProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Erroring) FarFutureEpoch(ctx context.Context) (phase0.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(consensusclient.FarFutureEpochProvider)
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
	next, isNext := s.next.(consensusclient.GenesisValidatorsRootProvider)
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
	next, isNext := s.next.(consensusclient.TargetAggregatorsPerCommitteeProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.TargetAggregatorsPerCommittee(ctx)
}

// AggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Erroring) AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.AggregateAttestationProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AggregateAttestation(ctx, slot, attestationDataRoot)
}

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Erroring) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.AggregateAttestationsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Erroring) AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.AttestationDataProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AttestationData(ctx, slot, committeeIndex)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Erroring) AttestationPool(ctx context.Context, slot phase0.Slot) ([]*phase0.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.AttestationPoolProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.AttestationPool(ctx, slot)
}

// SubmitAttestations submits attestations.
func (s *Erroring) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.AttestationsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitAttestations(ctx, attestations)
}

// SubmitSyncCommitteeMessages submits sync committee messages.
func (s *Erroring) SubmitSyncCommitteeMessages(ctx context.Context, messages []*altair.SyncCommitteeMessage) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.SyncCommitteeMessagesSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitSyncCommitteeMessages(ctx, messages)
}

// AttesterDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Erroring) AttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.AttesterDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.AttesterDutiesProvider)
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
	next, isNext := s.next.(consensusclient.BeaconBlockHeadersProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconBlockHeader(ctx, blockID)
}

// BeaconBlockRoot fetches a block's root given a block ID.
func (s *Erroring) BeaconBlockRoot(ctx context.Context, blockID string) (*phase0.Root, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconBlockRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconBlockRoot(ctx, blockID)
}

// BeaconCommittees fetches all beacon committees for the epoch at the given state.
func (s *Erroring) BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconCommitteesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconCommittees(ctx, stateID)
}

// BeaconCommitteesAtEpoch fetches all beacon committees for the given epoch at the given state.
func (s *Erroring) BeaconCommitteesAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) ([]*api.BeaconCommittee, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconCommitteesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconCommitteesAtEpoch(ctx, stateID, epoch)
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Erroring) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconBlockProposalProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

// SubmitBeaconBlock submits a beacon block.
func (s *Erroring) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.BeaconBlockSubmitter)
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
	next, isNext := s.next.(consensusclient.BeaconCommitteeSubscriptionsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

// SubmitSyncCommitteeSubscriptions subscribes to sync committees.
func (s *Erroring) SubmitSyncCommitteeSubscriptions(ctx context.Context, subscriptions []*api.SyncCommitteeSubscription) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.SyncCommitteeSubscriptionsSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitSyncCommitteeSubscriptions(ctx, subscriptions)
}

// BeaconState fetches a beacon state.
func (s *Erroring) BeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconStateProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconState(ctx, stateID)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Erroring) Events(ctx context.Context, topics []string, handler consensusclient.EventHandlerFunc) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.EventsProvider)
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
	next, isNext := s.next.(consensusclient.FinalityProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Finality(ctx, stateID)
}

// Fork fetches fork information for the given state.
func (s *Erroring) Fork(ctx context.Context, stateID string) (*phase0.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ForkProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Fork(ctx, stateID)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Erroring) ForkSchedule(ctx context.Context) ([]*phase0.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ForkScheduleProvider)
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
	next, isNext := s.next.(consensusclient.GenesisProvider)
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
	next, isNext := s.next.(consensusclient.NodeSyncingProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.NodeSyncing(ctx)
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
func (s *Erroring) ProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.ProposerDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ProposerDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ProposerDuties(ctx, epoch, validatorIndices)
}

// SyncCommitteeDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Erroring) SyncCommitteeDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.SyncCommitteeDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.SyncCommitteeDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SyncCommitteeDuties(ctx, epoch, validatorIndices)
}

// Spec provides the spec information of the chain.
func (s *Erroring) Spec(ctx context.Context) (map[string]interface{}, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.SpecProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Spec(ctx)
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Erroring) ValidatorBalances(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]phase0.Gwei, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ValidatorBalancesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ValidatorBalances(ctx, stateID, validatorIndices)
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
// will be applied.
func (s *Erroring) Validators(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]*api.Validator, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Validators(ctx, stateID, validatorIndices)
}

// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
// supplied no filter will be applied.
func (s *Erroring) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*api.Validator, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.ValidatorsByPubKey(ctx, stateID, validatorPubKeys)
}

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Erroring) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *phase0.SignedVoluntaryExit) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(consensusclient.VoluntaryExitSubmitter)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SubmitVoluntaryExit(ctx, voluntaryExit)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Erroring) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	if err := s.maybeError(ctx); err != nil {
		return phase0.Domain{}, err
	}
	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return phase0.Domain{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.Domain(ctx, domainType, epoch)
}

// GenesisTime provides the genesis time of the chain.
func (s *Erroring) GenesisTime(ctx context.Context) (time.Time, error) {
	if err := s.maybeError(ctx); err != nil {
		return time.Time{}, err
	}
	next, isNext := s.next.(consensusclient.GenesisTimeProvider)
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
	next, isNext := s.next.(consensusclient.DepositContractProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.DepositContract(ctx)
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Erroring) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.SignedBeaconBlockProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.SignedBeaconBlock(ctx, blockID)
}

// BeaconStateRoot fetches a beacon state root given a state ID.
func (s *Erroring) BeaconStateRoot(ctx context.Context, stateID string) (*phase0.Root, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(consensusclient.BeaconStateRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}
	return next.BeaconStateRoot(ctx, stateID)
}
