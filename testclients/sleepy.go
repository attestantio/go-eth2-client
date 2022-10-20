// Copyright © 2021, 2022 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Sleepy is an Ethereum 2 client that sleeps for a random amount of time within a
// set of bounds before continuing.
type Sleepy struct {
	minSleep time.Duration
	maxSleep time.Duration
	next     consensusclient.Service
}

// NewSleepy creates a new Ethereum 2 client that sleeps for random amount of time
// within a set of bounds between minSleep and maxSleep before continuing.
func NewSleepy(ctx context.Context,
	minSleep time.Duration,
	maxSleep time.Duration,
	next consensusclient.Service,
) (consensusclient.Service, error) {
	if next == nil {
		return nil, errors.New("no next service supplied")
	}
	if maxSleep < minSleep {
		return nil, errors.New("max sleep less than min sleep")
	}

	rand.Seed(time.Now().UnixNano())

	return &Sleepy{
		minSleep: minSleep,
		maxSleep: maxSleep,
		next:     next,
	}, nil
}

// Name returns the name of the client implementation.
func (s *Sleepy) Name() string {
	nextName := s.next.Name()
	return fmt.Sprintf("sleepy(%v,%v,%s)", s.minSleep, s.maxSleep, nextName)
}

// Address returns the address of the client.
func (s *Sleepy) Address() string {
	nextAddress := s.next.Address()
	return fmt.Sprintf("sleepy:%v,%v,%s", s.minSleep, s.maxSleep, nextAddress)
}

// sleep sleeps for a bounded amount of time.
func (s *Sleepy) sleep(ctx context.Context) {
	// #nosec G404
	duration := time.Duration(s.minSleep.Milliseconds()+rand.Int63n(s.maxSleep.Milliseconds()-s.minSleep.Milliseconds())) * time.Millisecond
	time.Sleep(duration)
}

// EpochFromStateID converts a state ID to its epoch.
func (s *Sleepy) EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.EpochFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
func (s *Sleepy) SlotFromStateID(ctx context.Context, stateID string) (phase0.Slot, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.SlotFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.SlotFromStateID(ctx, stateID)
}

// NodeVersion returns a free-text string with the node version.
func (s *Sleepy) NodeVersion(ctx context.Context) (string, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.NodeVersionProvider)
	if !isNext {
		return "", errors.New("next does not support this call")
	}
	return next.NodeVersion(ctx)
}

// SlotDuration provides the duration of a slot of the chain.
func (s *Sleepy) SlotDuration(ctx context.Context) (time.Duration, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.SlotDurationProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.SlotDuration(ctx)
}

// SlotsPerEpoch provides the slots per epoch of the chain.
func (s *Sleepy) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.SlotsPerEpochProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Sleepy) FarFutureEpoch(ctx context.Context) (phase0.Epoch, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.FarFutureEpochProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.FarFutureEpoch(ctx)
}

// GenesisValidatorsRoot provides the genesis validators root of the chain.
func (s *Sleepy) GenesisValidatorsRoot(ctx context.Context) ([]byte, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.GenesisValidatorsRootProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.GenesisValidatorsRoot(ctx)
}

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
func (s *Sleepy) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.TargetAggregatorsPerCommitteeProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}
	return next.TargetAggregatorsPerCommittee(ctx)
}

// AggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Sleepy) AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AggregateAttestationProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AggregateAttestation(ctx, slot, attestationDataRoot)
}

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Sleepy) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*phase0.SignedAggregateAndProof) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AggregateAttestationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitAggregateAttestations(ctx, aggregateAndProofs)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Sleepy) AttestationData(ctx context.Context, slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.AttestationData, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AttestationDataProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AttestationData(ctx, slot, committeeIndex)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Sleepy) AttestationPool(ctx context.Context, slot phase0.Slot) ([]*phase0.Attestation, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AttestationPoolProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AttestationPool(ctx, slot)
}

// SubmitAttestations submits attestations.
func (s *Sleepy) SubmitAttestations(ctx context.Context, attestations []*phase0.Attestation) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AttestationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitAttestations(ctx, attestations)
}

// AttesterDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Sleepy) AttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*apiv1.AttesterDuty, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.AttesterDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.AttesterDuties(ctx, epoch, validatorIndices)
}

// BeaconBlockHeader provides the block header of a given block ID.
func (s *Sleepy) BeaconBlockHeader(ctx context.Context, blockID string) (*apiv1.BeaconBlockHeader, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BeaconBlockHeadersProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.BeaconBlockHeader(ctx, blockID)
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Sleepy) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BeaconBlockProposalProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.BeaconBlockProposal(ctx, slot, randaoReveal, graffiti)
}

// SubmitBeaconBlock submits a beacon block.
func (s *Sleepy) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BeaconBlockSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitBeaconBlock(ctx, block)
}

// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
func (s *Sleepy) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*apiv1.BeaconCommitteeSubscription) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BeaconCommitteeSubscriptionsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitBeaconCommitteeSubscriptions(ctx, subscriptions)
}

// SubmitProposalPreparations submits proposal preparations.
func (s *Sleepy) SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ProposalPreparationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitProposalPreparations(ctx, preparations)
}

// SubmitBlindedBeaconBlock submits a beacon block.
func (s *Sleepy) SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BlindedBeaconBlockSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitBlindedBeaconBlock(ctx, block)
}

// SubmitValidatorRegistrations submits a validator registration.
func (s *Sleepy) SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ValidatorRegistrationsSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitValidatorRegistrations(ctx, registrations)
}

// BeaconState fetches a beacon state.
func (s *Sleepy) BeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.BeaconStateProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.BeaconState(ctx, stateID)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Sleepy) Events(ctx context.Context, topics []string, handler consensusclient.EventHandlerFunc) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.EventsProvider)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.Events(ctx, topics, handler)
}

// Finality provides the finality given a state ID.
func (s *Sleepy) Finality(ctx context.Context, stateID string) (*apiv1.Finality, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.FinalityProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Finality(ctx, stateID)
}

// Fork fetches fork information for the given state.
func (s *Sleepy) Fork(ctx context.Context, stateID string) (*phase0.Fork, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ForkProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Fork(ctx, stateID)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Sleepy) ForkSchedule(ctx context.Context) ([]*phase0.Fork, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ForkScheduleProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.ForkSchedule(ctx)
}

// Genesis fetches genesis information for the chain.
func (s *Sleepy) Genesis(ctx context.Context) (*apiv1.Genesis, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.GenesisProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Genesis(ctx)
}

// NodeSyncing provides the state of the node's synchronization with the chain.
func (s *Sleepy) NodeSyncing(ctx context.Context) (*apiv1.SyncState, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.NodeSyncingProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.NodeSyncing(ctx)
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
func (s *Sleepy) ProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*apiv1.ProposerDuty, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ProposerDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.ProposerDuties(ctx, epoch, validatorIndices)
}

// Spec provides the spec information of the chain.
func (s *Sleepy) Spec(ctx context.Context) (map[string]interface{}, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.SpecProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Spec(ctx)
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Sleepy) ValidatorBalances(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]phase0.Gwei, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ValidatorBalancesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.ValidatorBalances(ctx, stateID, validatorIndices)
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
// will be applied.
func (s *Sleepy) Validators(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex) (map[phase0.ValidatorIndex]*apiv1.Validator, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.Validators(ctx, stateID, validatorIndices)
}

// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
// supplied no filter will be applied.
func (s *Sleepy) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []phase0.BLSPubKey) (map[phase0.ValidatorIndex]*apiv1.Validator, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}
	return next.ValidatorsByPubKey(ctx, stateID, validatorPubKeys)
}

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Sleepy) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *phase0.SignedVoluntaryExit) error {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.VoluntaryExitSubmitter)
	if !isNext {
		return errors.New("next does not support this call")
	}
	return next.SubmitVoluntaryExit(ctx, voluntaryExit)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Sleepy) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return phase0.Domain{}, errors.New("next does not support this call")
	}
	return next.Domain(ctx, domainType, epoch)
}

// GenesisTime provides the genesis time of the chain.
func (s *Sleepy) GenesisTime(ctx context.Context) (time.Time, error) {
	s.sleep(ctx)
	next, isNext := s.next.(consensusclient.GenesisTimeProvider)
	if !isNext {
		return time.Time{}, errors.New("next does not support this call")
	}
	return next.GenesisTime(ctx)
}
