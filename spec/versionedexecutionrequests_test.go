// Copyright © 2026 Attestant Limited.
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

package spec_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestVersionedExecutionRequestsAccessors(t *testing.T) {
	deposits := []*electra.DepositRequest{{Index: 1}}
	withdrawals := []*electra.WithdrawalRequest{{}}
	consolidations := []*electra.ConsolidationRequest{{}}
	builderDeposits := []*gloas.BuilderDepositRequest{{}}
	builderExits := []*gloas.BuilderExitRequest{{}}
	populated := &spec.VersionedExecutionRequests{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.ExecutionRequests{
			Deposits:        deposits,
			Withdrawals:     withdrawals,
			Consolidations:  consolidations,
			BuilderDeposits: builderDeposits,
			BuilderExits:    builderExits,
		},
	}
	empty := &spec.VersionedExecutionRequests{Version: spec.DataVersionGloas}
	tests := []struct {
		name     string
		accessor func(*spec.VersionedExecutionRequests) (any, error)
		want     any
	}{
		{
			name: "Deposits",
			accessor: func(v *spec.VersionedExecutionRequests) (any, error) {
				return v.Deposits()
			},
			want: deposits,
		},
		{
			name: "Withdrawals",
			accessor: func(v *spec.VersionedExecutionRequests) (any, error) {
				return v.Withdrawals()
			},
			want: withdrawals,
		},
		{
			name: "Consolidations",
			accessor: func(v *spec.VersionedExecutionRequests) (any, error) {
				return v.Consolidations()
			},
			want: consolidations,
		},
		{
			name: "BuilderDeposits",
			accessor: func(v *spec.VersionedExecutionRequests) (any, error) {
				return v.BuilderDeposits()
			},
			want: builderDeposits,
		},
		{
			name: "BuilderExits",
			accessor: func(v *spec.VersionedExecutionRequests) (any, error) {
				return v.BuilderExits()
			},
			want: builderExits,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.accessor(populated)
			require.NoError(t, err)
			require.Equal(t, test.want, got)

			_, err = test.accessor(empty)
			require.EqualError(t, err, "no gloas execution requests")
		})
	}
}

func TestVersionedExecutionRequestsState(t *testing.T) {
	tests := []struct {
		name     string
		requests *spec.VersionedExecutionRequests
		empty    bool
		expected string
	}{
		{
			name: "Gloas",
			requests: &spec.VersionedExecutionRequests{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.ExecutionRequests{},
			},
			expected: "deposits:",
		},
		{
			name: "GloasNil",
			requests: &spec.VersionedExecutionRequests{
				Version: spec.DataVersionGloas,
			},
			empty: true,
		},
		{
			name:     "UnknownVersion",
			requests: &spec.VersionedExecutionRequests{},
			empty:    true,
			expected: "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.empty, test.requests.IsEmpty())
			if test.expected == "" {
				require.Empty(t, test.requests.String())
			} else {
				require.Contains(t, test.requests.String(), test.expected)
			}
		})
	}
}
