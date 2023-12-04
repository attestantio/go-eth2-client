// Copyright © 2020, 2021 Attestant Limited.
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

package v1

import (
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ValidatorState defines the state of the validator.
type ValidatorState int

const (
	// ValidatorStateUnknown means no information can be found about the validator.
	ValidatorStateUnknown ValidatorState = iota
	// ValidatorStatePendingInitialized means the validator is not yet in the queue to be activated.
	ValidatorStatePendingInitialized
	// ValidatorStatePendingQueued means the validator is in the queue to be activated.
	ValidatorStatePendingQueued
	// ValidatorStateActiveOngoing means the validator is active.
	ValidatorStateActiveOngoing
	// ValidatorStateActiveExiting means the validator is active but exiting.
	ValidatorStateActiveExiting
	// ValidatorStateActiveSlashed means the validator is active but exiting due to being slashed.
	ValidatorStateActiveSlashed
	// ValidatorStateExitedUnslashed means the validator has exited without being slashed.
	ValidatorStateExitedUnslashed
	// ValidatorStateExitedSlashed means the validator has exited due to being slashed.
	ValidatorStateExitedSlashed
	// ValidatorStateWithdrawalPossible means it is possible to withdraw funds from the validator.
	ValidatorStateWithdrawalPossible
	// ValidatorStateWithdrawalDone means funds have been withdrawn from the validator.
	ValidatorStateWithdrawalDone
)

var validatorStateStrings = [...]string{
	"unknown",
	"pending_initialized",
	"pending_queued",
	"active_ongoing",
	"active_exiting",
	"active_slashed",
	"exited_unslashed",
	"exited_slashed",
	"withdrawal_possible",
	"withdrawal_done",
}

// MarshalJSON implements json.Marshaler.
func (v *ValidatorState) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", validatorStateStrings[*v])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *ValidatorState) UnmarshalJSON(input []byte) error {
	var err error
	switch strings.ToLower(string(input)) {
	case `"unknown"`:
		*v = ValidatorStateUnknown
	case `"pending_initialized"`:
		*v = ValidatorStatePendingInitialized
	case `"pending_queued"`:
		*v = ValidatorStatePendingQueued
	case `"active_ongoing"`:
		*v = ValidatorStateActiveOngoing
	case `"active_exiting"`:
		*v = ValidatorStateActiveExiting
	case `"active_slashed"`:
		*v = ValidatorStateActiveSlashed
	case `"exited_unslashed"`:
		*v = ValidatorStateExitedUnslashed
	case `"exited_slashed"`:
		*v = ValidatorStateExitedSlashed
	case `"withdrawal_possible"`:
		*v = ValidatorStateWithdrawalPossible
	case `"withdrawal_done"`:
		*v = ValidatorStateWithdrawalDone
	default:
		err = fmt.Errorf("unrecognised validator state %s", string(input))
	}

	return err
}

func (v ValidatorState) String() string {
	if v < 0 || int(v) >= len(validatorStateStrings) {
		return validatorStateStrings[0] // unknown
	}

	return validatorStateStrings[v]
}

// IsPending returns true if the validator is pending.
func (v ValidatorState) IsPending() bool {
	return v == ValidatorStatePendingInitialized ||
		v == ValidatorStatePendingQueued
}

// IsActive returns true if the validator is active.
func (v ValidatorState) IsActive() bool {
	return v == ValidatorStateActiveOngoing ||
		v == ValidatorStateActiveExiting ||
		v == ValidatorStateActiveSlashed
}

// HasActivated returns true if the validator has activated.
func (v ValidatorState) HasActivated() bool {
	return v == ValidatorStateActiveOngoing ||
		v == ValidatorStateActiveExiting ||
		v == ValidatorStateActiveSlashed ||
		v == ValidatorStateExitedUnslashed ||
		v == ValidatorStateExitedSlashed ||
		v == ValidatorStateWithdrawalPossible ||
		v == ValidatorStateWithdrawalDone
}

// IsAttesting returns true if the validator should be attesting.
func (v ValidatorState) IsAttesting() bool {
	return v == ValidatorStateActiveOngoing || v == ValidatorStateActiveExiting
}

// IsExited returns true if the validator is exited.
func (v ValidatorState) IsExited() bool {
	return v == ValidatorStateExitedUnslashed ||
		v == ValidatorStateExitedSlashed
}

// HasExited returns true if the validator has exited.
func (v ValidatorState) HasExited() bool {
	return v == ValidatorStateExitedUnslashed ||
		v == ValidatorStateExitedSlashed ||
		v == ValidatorStateWithdrawalPossible ||
		v == ValidatorStateWithdrawalDone
}

// HasBalance returns true if the validator has a balance.
func (v ValidatorState) HasBalance() bool {
	return v != ValidatorStateUnknown
}

// ValidatorToState is a helper that calculates the validator status given a validator struct.
func ValidatorToState(validator *phase0.Validator,
	balance *phase0.Gwei,
	currentEpoch phase0.Epoch,
	farFutureEpoch phase0.Epoch,
) ValidatorState {
	if validator == nil {
		return ValidatorStateUnknown
	}

	switch {
	case validator.ActivationEpoch > currentEpoch:
		// Pending.
		if validator.ActivationEligibilityEpoch == farFutureEpoch {
			return ValidatorStatePendingInitialized
		}

		return ValidatorStatePendingQueued
	case validator.ExitEpoch == farFutureEpoch:
		// Active ongoing.
		return ValidatorStateActiveOngoing
	case validator.ExitEpoch > currentEpoch:
		// Active exiting.
		if validator.Slashed {
			return ValidatorStateActiveSlashed
		}

		return ValidatorStateActiveExiting
	case validator.WithdrawableEpoch > currentEpoch:
		// Exited.
		if validator.Slashed {
			return ValidatorStateExitedSlashed
		}

		return ValidatorStateExitedUnslashed
	case balance != nil && *balance == 0:
		return ValidatorStateWithdrawalDone
	default:
		return ValidatorStateWithdrawalPossible
	}
}
