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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
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
func (s *Erroring) PrysmAttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorPubKeys []phase0.BLSPubKey) ([]*api.AttesterDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmAttesterDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.PrysmAttesterDuties(ctx, epoch, validatorPubKeys)
}

// PrysmProposerDuties obtains proposer duties with prysm-specific parameters.
func (s *Erroring) PrysmProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorPubKeys []phase0.BLSPubKey) ([]*api.ProposerDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmProposerDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.PrysmProposerDuties(ctx, epoch, validatorPubKeys)
}

// PrysmValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIDs is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Erroring) PrysmValidatorBalances(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]phase0.Gwei, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.PrysmValidatorBalancesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.PrysmValidatorBalances(ctx, stateID, validatorPubKeys)
}

// EpochFromStateID converts a state ID to its epoch.
func (s *Erroring) EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.EpochFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
func (s *Erroring) SlotFromStateID(ctx context.Context, stateID string) (phase0.Slot, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.SlotFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
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
		return "", errors.New("next does not support this call")
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
		return 0, errors.New("next does not support this call")
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
		return 0, errors.New("next does not support this call")
	}
	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Erroring) FarFutureEpoch(ctx context.Context) (phase0.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}
	next, isNext := s.next.(eth2client.FarFutureEpochProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
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
		return 0, errors.New("next does not support this call")
	}
	return next.TargetAggregatorsPerCommittee(ctx)
}

// AggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Erroring) AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AggregateAttestationProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AggregateAttestation(ctx, slot, attestationDataRoot)
}

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Erroring) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.AggregateAttestationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Erroring) AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttestationDataProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AttestationData(ctx, slot, committeeIndex)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Erroring) AttestationPool(ctx context.Context, slot phase0.Slot) ([]*phase0.Attestation, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttestationPoolProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AttestationPool(ctx, slot)
}

// SubmitAttestations submits attestations.
func (s *Erroring) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.AttestationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitAttestations(ctx, attestations)
}

// AttesterDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Erroring) AttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.AttesterDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.AttesterDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
	}
	return next.BeaconBlockHeader(ctx, blockID)
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Erroring) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconBlockProposalProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.BeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

// SubmitBeaconBlock submits a beacon block.
func (s *Erroring) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.BeaconBlockSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
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
		return errors.New("next does not support this call")
	}
	return next.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

// BeaconState fetches a beacon state.
func (s *Erroring) BeaconState(ctx context.Context, stateID string) (*phase0.BeaconState, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.BeaconStateProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
		return errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
	}
	return next.Finality(ctx, stateID)
}

// Fork fetches fork information for the given state.
func (s *Erroring) Fork(ctx context.Context, stateID string) (*phase0.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ForkProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Fork(ctx, stateID)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Erroring) ForkSchedule(ctx context.Context) ([]*phase0.Fork, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ForkScheduleProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
	}
	return next.NodeSyncing(ctx)
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
func (s *Erroring) ProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.ProposerDuty, error) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}
	next, isNext := s.next.(eth2client.ProposerDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
		return nil, errors.New("next does not support this call")
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
	next, isNext := s.next.(eth2client.ValidatorBalancesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
	next, isNext := s.next.(eth2client.ValidatorsProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
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
	next, isNext := s.next.(eth2client.ValidatorsProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.ValidatorsByPubKey(ctx, stateID, validatorPubKeys)
}

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Erroring) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *phase0.SignedVoluntaryExit) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}
	next, isNext := s.next.(eth2client.VoluntaryExitSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitVoluntaryExit(ctx, voluntaryExit)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Erroring) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	if err := s.maybeError(ctx); err != nil {
		return phase0.Domain{}, err
	}
	next, isNext := s.next.(eth2client.DomainProvider)
	if !isNext {
		return phase0.Domain{}, errors.New("next does not support this call")
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
		return time.Time{}, errors.New("next does not support this call")
	}
	return next.GenesisTime(ctx)
}
