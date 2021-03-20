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

package v2_test

import (
	"encoding/json"
	"strings"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v2"
	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

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
			input:      []byte(`"Pending_queued"`),
			isPending:  true,
			hasBalance: true,
		},
		{
			name:       "PendingInitialized",
			input:      []byte(`"Pending_initialized"`),
			isPending:  true,
			hasBalance: true,
		},
		{
			name:         "ActiveOngoing",
			input:        []byte(`"Active_ongoing"`),
			isActive:     true,
			hasActivated: true,
			isAttesting:  true,
			hasBalance:   true,
		},
		{
			name:         "ActiveExiting",
			input:        []byte(`"Active_exiting"`),
			isActive:     true,
			hasActivated: true,
			isAttesting:  true,
			hasBalance:   true,
		},
		{
			name:         "ActiveSlashed",
			input:        []byte(`"Active_slashed"`),
			isActive:     true,
			hasActivated: true,
			hasBalance:   true,
		},
		{
			name:         "ExitedUnslashed",
			input:        []byte(`"Exited_unslashed"`),
			hasActivated: true,
			isExited:     true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "ExitedSlashed",
			input:        []byte(`"Exited_slashed"`),
			hasActivated: true,
			isExited:     true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "WithdrawalPossible",
			input:        []byte(`"Withdrawal_possible"`),
			hasActivated: true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:         "WithdrawalDone",
			input:        []byte(`"Withdrawal_done"`),
			hasActivated: true,
			hasExited:    true,
			hasBalance:   true,
		},
		{
			name:  "Unknown",
			input: []byte(`"Unknown"`),
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
	farFutureEpoch := spec.Epoch(99999)
	currentEpoch := spec.Epoch(100)
	tests := []struct {
		name      string
		validator *spec.Validator
		state     api.ValidatorState
	}{
		{
			name:  "Nil",
			state: api.ValidatorStateUnknown,
		},
		{
			name: "PendingInitialized",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: farFutureEpoch,
				ActivationEpoch:            farFutureEpoch,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingInitialized,
		},
		{
			name: "PendingQueued",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch + 10,
				ActivationEpoch:            farFutureEpoch,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingQueued,
		},
		{
			name: "PendingQueued",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch + 10,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStatePendingQueued,
		},
		{
			name: "ActiveOngoing",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  farFutureEpoch,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveOngoing,
		},
		{
			name: "ActiveExiting",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 10,
				WithdrawableEpoch:          farFutureEpoch,
			},
			state: api.ValidatorStateActiveExiting,
		},
		{
			name: "ActiveSlashed",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch + 10,
				WithdrawableEpoch:          farFutureEpoch,
				Slashed:                    true,
			},
			state: api.ValidatorStateActiveSlashed,
		},
		{
			name: "ExitedUnslashed",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch + 50,
			},
			state: api.ValidatorStateExitedUnslashed,
		},
		{
			name: "ExitedUnslashed",
			validator: &spec.Validator{
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
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
			},
			state: api.ValidatorStateWithdrawalPossible,
		},
		{
			name: "WithdrawalPossibleSlashed",
			validator: &spec.Validator{
				ActivationEligibilityEpoch: currentEpoch - 50,
				ActivationEpoch:            currentEpoch - 40,
				ExitEpoch:                  currentEpoch - 30,
				WithdrawableEpoch:          currentEpoch - 20,
				Slashed:                    true,
			},
			state: api.ValidatorStateWithdrawalPossible,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := api.ValidatorToState(test.validator, currentEpoch, farFutureEpoch)
			assert.Equal(t, test.state, state)
		})
	}
}
