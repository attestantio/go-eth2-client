// Copyright © 2020 - 2026 Attestant Limited.
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

package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// attestationDataService returns a service whose every request is answered from
// memory: the two endpoints a service needs to come up, the configuration given by
// specEntries, and data as the answer to every attestation-data request.
//
// No live node can stand in for it. Which vote a conformant node returns is decided by
// its own fork choice, so a caller cannot ask for the payload-FULL one: eight
// consecutive requests to a Gloas devnet answered 0 every time, and a request that is
// answered 0 exercises the same branch as every pre-Gloas answer. The other half of the
// gate cannot be asked of that node at all, since a chain is on one side of a fork or
// the other. Synthesising the node makes both sides reachable on demand, and reachable
// on an endpoint that has never heard of Gloas.
func attestationDataService(ctx context.Context,
	t *testing.T,
	specEntries string,
	data *phase0.AttestationData,
) client.Service {
	t.Helper()

	// Marshalled here rather than in the handler because a failure must be reported
	// from the test's own goroutine.
	encoded, err := json.Marshal(data)
	require.NoError(t, err)

	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		// Both of these, or New leaves the service inactive and every call below
		// returns ErrNotActive without ever reaching the wire.
		case "/eth/v1/node/version":
			fmt.Fprint(w, `{"data":{"version":"stub"}}`)
		case "/eth/v1/node/syncing":
			fmt.Fprint(w, `{"data":{"head_slot":"1","sync_distance":"0","is_syncing":false,"is_optimistic":false,"el_offline":false}}`)
		case "/eth/v1/config/spec":
			fmt.Fprintf(w, `{"data":{%s}}`, specEntries)
		case "/eth/v1/validator/attestation_data":
			fmt.Fprintf(w, `{"data":%s}`, encoded)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	service, err := http.New(ctx, http.WithAddress(server.URL))
	require.NoError(t, err)

	return service
}

// TestAttestationDataGloasPayloadVote covers the fork gate on the index assertion.
//
// Gloas repurposes AttestationData.index as a payload-availability vote, so the
// Electra-era expectation that it is 0 turns a conformant node's payload-FULL answer
// into ErrInconsistentResult and leaves the validator unable to attest on those
// slots. What has to hold is that the two eras are told apart: bounded to {0, 1} at
// or after the fork, pinned to 0 before it.
func TestAttestationDataGloasPayloadVote(t *testing.T) {
	ctx := context.Background()

	const slotsPerEpoch = 32
	// Gloas at epoch 100, so slot 3200 is the fork's first slot and 3199 its last
	// pre-fork one.
	const gloasForkEpoch = "100"
	const gloasSlot = phase0.Slot(3200)

	tests := []struct {
		name string
		// Absent from the configuration when empty, as it is on every node that
		// does not know the fork.
		gloasForkEpoch string
		slot           phase0.Slot
		index          phase0.CommitteeIndex
		err            string
	}{
		{
			// The payload-EMPTY vote, which is also every pre-Gloas value, so this
			// is the case the electra-era expectation already admitted.
			name:           "PayloadEmptyVoteAccepted",
			gloasForkEpoch: gloasForkEpoch,
			slot:           gloasSlot,
			index:          0,
		},
		{
			// The payload-FULL vote: the case that was discarded, leaving a
			// validator client unable to attest on any slot where its node
			// legitimately answered 1.
			name:           "PayloadFullVoteAccepted",
			gloasForkEpoch: gloasForkEpoch,
			slot:           gloasSlot,
			index:          1,
		},
		{
			// Two is not a payload status.  The bound stays asserted, so widening
			// the vote does not amount to no longer checking the field.
			name:           "VoteAboveOneRejected",
			gloasForkEpoch: gloasForkEpoch,
			slot:           gloasSlot,
			index:          2,
			err:            "attestation data payload availability vote 2; expected 0 or 1",
		},
		{
			// One slot earlier the field is still a committee index that electra
			// zeroed, so the same answer is still inconsistent.
			name:           "LastPreGloasSlotStillClamped",
			gloasForkEpoch: gloasForkEpoch,
			slot:           gloasSlot - 1,
			index:          1,
			err:            "attestation data for committee index 1; expected 0",
		},
		{
			// A node that has never heard of the fork cannot be voting on payloads,
			// whatever the slot number, so its answer is held to electra's rule.
			name:  "NodeThatDoesNotKnowGloasStillClamped",
			slot:  gloasSlot,
			index: 1,
			err:   "attestation data for committee index 1; expected 0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Electra at epoch 0 so that the pre-Gloas expectation in force is
			// always the Electra-era clamp to 0, which is what the Gloas arm has
			// to displace.
			specEntries := fmt.Sprintf(`"SLOTS_PER_EPOCH":"%d","ELECTRA_FORK_EPOCH":"0"`, slotsPerEpoch)
			if test.gloasForkEpoch != "" {
				specEntries += fmt.Sprintf(`,"GLOAS_FORK_EPOCH":"%s"`, test.gloasForkEpoch)
			}

			// Slot and Index are the only fields the guard under test reads; the
			// checkpoints are present because the response would not decode
			// without them, not because their contents matter.
			service := attestationDataService(ctx, t, specEntries, &phase0.AttestationData{
				Slot:   test.slot,
				Index:  test.index,
				Source: &phase0.Checkpoint{},
				Target: &phase0.Checkpoint{},
			})

			response, err := service.(client.AttestationDataProvider).AttestationData(ctx, &api.AttestationDataOpts{
				Slot: test.slot,
			})

			if test.err != "" {
				require.ErrorIs(t, err, client.ErrInconsistentResult)
				require.ErrorContains(t, err, test.err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.index, response.Data.Index)
		})
	}
}

func TestAttestationData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	// Need to fetch current slot for attestation data.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)

	tests := []struct {
		name    string
		opts    *api.AttestationDataOpts
		err     []string
		errCode int
	}{
		{
			name: "Good",
			opts: &api.AttestationDataOpts{
				Slot: phase0.Slot(uint64(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / uint64(slotDuration.Seconds())),
			},
		},
		{
			name: "NilOpts",
			err:  []string{"no options specified"},
		},
		{
			name: "BadSlot",
			opts: &api.AttestationDataOpts{
				Slot: 999999999,
			},
			errCode: 400,
			err:     []string{"request slot 999999999 is more than one slot past the current slot", "slot 999999999 is not the current slot"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AttestationDataProvider).AttestationData(ctx, test.opts)
			switch {
			case len(test.err) > 0:
				found := false
				for _, errMsg := range test.err {
					if strings.Contains(err.Error(), errMsg) {
						require.ErrorContains(t, err, errMsg)
						found = true
						break
					}
				}
				require.True(t, found, "error message not found in error: %s", err.Error())
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
			}
		})
	}
}
