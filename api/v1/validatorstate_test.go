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

package v1_test

import (
	"encoding/json"
	"strings"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func gweiPtr(input phase0.Gwei) *phase0.Gwei {
	return &input
}

func TestValidatorStateJSON(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		isPending    bool
		isActive     bool
		hasActivated bool
		isAttesting  bool
		isExited     bool
		hasExited    bool
		hasBalance   bool
		err          string
	}{
		{
			name:       "PendingQueued",
			input:      []byte(`"pending_queued"`),
			isPending:  true,
			hasBalance: true,
		},
		{
			name:       "PendingInitialized",
			input:      []byte(`"pending_initialized"`),
			isPending:  true,
			hasBalance: true,
		},
		{
			name:         "ActiveOngoing",
			input:        []byte(`"active_ongoing"`),
			isActive:     true,
			hasActivated: true,
			isAttesting:  true,
			hasBalance:   true,
		},
		{
			name:         "ActiveExiting",
			input:        []byte(`"active_exiting"`),
			isActive:     true,
			hasActivated: true,
			isAttesting:  true,
			hasBalance:   true,
		},
		{
			name:         "ActiveSlashed",
			input:        []byte(`"active_slashed"`),
			isActive:     true,
			hasActivated: true,
			hasBalance:   true,
		},
		{
			name:         "ExitedUnslashed",
			input:        []byte(`"exited_unslashed"`),
			hasActivated: true,
			isExited:     true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "ExitedSlashed",
			input:        []byte(`"exited_slashed"`),
			hasActivated: true,
			isExited:     true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "WithdrawalPossible",
			input:        []byte(`"withdrawal_possible"`),
			hasActivated: true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "WithdrawalDone",
			input:        []byte(`"withdrawal_done"`),
			hasActivated: true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:  "Unknown",
			input: []byte(`"unknown"`),
		},
		{
			name:  "Invalid",
			input: []byte(`"Invalid"`),
			err:   "unrecognised validator state \"Invalid\"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ValidatorState
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, test.isPending, res.IsPending())
				assert.Equal(t, test.isActive, res.IsActive())
				assert.Equal(t, test.hasActivated, res.HasActivated())
				assert.Equal(t, test.isAttesting, res.IsAttesting())
				assert.Equal(t, test.isExited, res.IsExited())
				assert.Equal(t, test.hasExited, res.HasExited())
				assert.Equal(t, test.hasBalance, res.HasBalance())
				assert.Equal(t, strings.Trim(string(rt), `"`), res.String())
			}
		})
	}
}

func TestValidatorToState(t *testing.T) {
	farFutureEpoch := phase0.Epoch(99999)
	currentEpoch := phase0.Epoch(100)
	tests := []struct {
		name      string
		validator *phase0.Validator
		balance   *phase0.Gwei
		state     api.ValidatorState
	}{
		{
			name:  "Nil",
			state: api.ValidatorStateUnknown,
		},
		{
			name: "PendingInitialized",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: farFutureEpoch,
				ActivationEpoch:            farFutureEpoch,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingInitialized,
		},
		{
			name: "PendingQueued",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch + 10,
				ActivationEpoch:            farFutureEpoch,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingQueued,
		},
		{
			name: "PendingQueued",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch + 10,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingQueued,
		},
		{
			name: "ActiveOngoingNext",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch + 1,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingQueued,
		},
		{
			name: "ActiveOngoingThis",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveOngoing,
		},
		{
			name: "ActiveOngoing",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveOngoing,
		},
		{
			name: "ActiveExiting",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 10,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveExiting,
		},
		{
			name: "ActiveSlashed",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 10,
				WithdrawableEpoch:          farFutureEpoch,
				Slashed:                    true,
			},
			state: api.ValidatorStateActiveSlashed,
		},
		{
			name: "ExitedNext",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 1,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveExiting,
		},
		{
			name: "ExitedThis",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateExitedUnslashed,
		},
		{
			name: "ExitedUnslashed",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch + 50,
			},
			state: api.ValidatorStateExitedUnslashed,
		},
		{
			name: "ExitedSlashedNext",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 1,
				WithdrawableEpoch:          farFutureEpoch,
				Slashed:                    true,
			},
			state: api.ValidatorStateActiveSlashed,
		},
		{
			name: "ExitedSlashedThis",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch,
				WithdrawableEpoch:          farFutureEpoch,
				Slashed:                    true,
			},
			state: api.ValidatorStateExitedSlashed,
		},
		{
			name: "ExitedSlashed",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch + 50,
				Slashed:                    true,
			},
			state: api.ValidatorStateExitedSlashed,
		},
		{
			name: "WithdrawalPossible",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
			},
			state: api.ValidatorStateWithdrawalPossible,
		},
		{
			name: "WithdrawalPossibleSlashed",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
				Slashed:                    true,
			},
			state: api.ValidatorStateWithdrawalPossible,
		},
		{
			name: "WithdrawalPossibleBalance",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
			},
			balance: gweiPtr(5),
			state:   api.ValidatorStateWithdrawalPossible,
		},
		{
			name: "WithdrawalDone",
			validator: &phase0.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
			},
			balance: gweiPtr(0),
			state:   api.ValidatorStateWithdrawalDone,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := api.ValidatorToState(test.validator, test.balance, currentEpoch, farFutureEpoch)
			assert.Equal(t, test.state, state)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		state    api.ValidatorState
		expected string
	}{
		{
			name:     "valid state",
			state:    api.ValidatorStateActiveOngoing,
			expected: "active_ongoing",
		},
		{
			name:     "negative index",
			state:    -1,
			expected: "unknown",
		},
		{
			name:     "edge bound index",
			state:    10,
			expected: "unknown",
		},
		{
			name:     "high out of bound index",
			state:    250,
			expected: "unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := test.state.String()
			require.Equal(t, test.expected, resp)
		})
	}
}
