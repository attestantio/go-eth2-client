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

package http

import (
	"testing"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestVerifyPTCDuties covers the epoch-range check applied to duties returned
// by the node.  It is a pure function over plain Go values, so it needs no
// beacon node.
func TestVerifyPTCDuties(t *testing.T) {
	duty := func(slot phase0.Slot) *apiv1.PTCDuty {
		return &apiv1.PTCDuty{ValidatorIndex: 8, Slot: slot}
	}

	tests := []struct {
		name          string
		epoch         phase0.Epoch
		slotsPerEpoch uint64
		duties        []*apiv1.PTCDuty
		err           string
	}{
		{
			name:          "Empty",
			epoch:         7537,
			slotsPerEpoch: 8,
		},
		{
			name:          "InRange",
			epoch:         7537,
			slotsPerEpoch: 8,
			// Epoch 7537 at 8 slots per epoch spans [60296,60303].
			duties: []*apiv1.PTCDuty{duty(60296), duty(60300), duty(60303)},
		},
		{
			name:          "BeforeEpoch",
			epoch:         7537,
			slotsPerEpoch: 8,
			duties:        []*apiv1.PTCDuty{duty(60295)},
			err:           "received PTC duty for slot 60295 outside of range [60296,60303]",
		},
		{
			name:          "AfterEpoch",
			epoch:         7537,
			slotsPerEpoch: 8,
			duties:        []*apiv1.PTCDuty{duty(60300), duty(60304)},
			err:           "received PTC duty for slot 60304 outside of range [60296,60303]",
		},
		{
			// A node reporting SLOTS_PER_EPOCH of 0 would otherwise underflow
			// the end of the range to MaxUint64 and accept every duty, so the
			// check would silently pass anything.
			name:          "ZeroSlotsPerEpoch",
			epoch:         7537,
			slotsPerEpoch: 0,
			duties:        []*apiv1.PTCDuty{duty(60300)},
			err:           "invalid slots per epoch 0",
		},
		{
			name:          "NilDuty",
			epoch:         7537,
			slotsPerEpoch: 8,
			duties:        []*apiv1.PTCDuty{nil},
			err:           "received nil PTC duty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := verifyPTCDuties(test.epoch, test.slotsPerEpoch, test.duties)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
