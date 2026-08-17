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

package gloas_test

import (
	"encoding/json"
	"regexp"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	require "github.com/stretchr/testify/require"
)

// populatedBeaconState returns a state with every one of BeaconState's 46 fields
// carrying a distinctive non-zero value, which is what makes a round-trip an audit
// rather than a formality: a field left at its zero value hides a codec that
// mishandles it, since an absent or empty value round-trips clean either way.
//
// Vectors carry two entries rather than their spec length. Nothing in the JSON
// codec checks a vector's width, so a spec-sized 8192 block_roots would inflate
// the fixture without testing anything more.
func populatedBeaconState() *gloas.BeaconState {
	state := &gloas.BeaconState{}

	state.GenesisTime = 1
	state.GenesisValidatorsRoot = phase0.Root{0x01}
	state.Slot = 2
	state.Fork = &phase0.Fork{
		PreviousVersion: phase0.Version{0x02},
		CurrentVersion:  phase0.Version{0x03},
		Epoch:           3,
	}
	state.LatestBlockHeader = &phase0.BeaconBlockHeader{
		Slot:          4,
		ProposerIndex: 5,
		ParentRoot:    phase0.Root{0x04},
		StateRoot:     phase0.Root{0x05},
		BodyRoot:      phase0.Root{0x06},
	}
	state.BlockRoots = []phase0.Root{{0x07}, {0x08}}
	state.StateRoots = []phase0.Root{{0x09}, {0x0a}}
	state.HistoricalRoots = []phase0.Root{{0x0b}, {0x0c}}
	state.ETH1Data = &phase0.ETH1Data{
		DepositRoot:  phase0.Root{0x0d},
		DepositCount: 6,
		BlockHash:    make([]byte, phase0.Hash32Length),
	}
	state.ETH1DataVotes = []*phase0.ETH1Data{{
		DepositRoot:  phase0.Root{0x0e},
		DepositCount: 7,
		BlockHash:    make([]byte, phase0.Hash32Length),
	}}
	state.ETH1DepositIndex = 8
	state.Validators = []*phase0.Validator{{
		PublicKey:                  phase0.BLSPubKey{0x0f},
		WithdrawalCredentials:      make([]byte, 32),
		EffectiveBalance:           9,
		ActivationEligibilityEpoch: 10,
		ActivationEpoch:            11,
		ExitEpoch:                  12,
		WithdrawableEpoch:          13,
	}}
	state.Balances = []phase0.Gwei{14, 15}
	state.RANDAOMixes = []phase0.Root{{0x10}, {0x11}}
	state.Slashings = []phase0.Gwei{16, 17}
	state.PreviousEpochParticipation = []altair.ParticipationFlags{1, 2}
	state.CurrentEpochParticipation = []altair.ParticipationFlags{3, 4}
	state.JustificationBits = bitfield.Bitvector4{0x0b}
	state.PreviousJustifiedCheckpoint = &phase0.Checkpoint{Epoch: 18, Root: phase0.Root{0x12}}
	state.CurrentJustifiedCheckpoint = &phase0.Checkpoint{Epoch: 19, Root: phase0.Root{0x13}}
	state.FinalizedCheckpoint = &phase0.Checkpoint{Epoch: 20, Root: phase0.Root{0x14}}
	state.InactivityScores = []uint64{21, 22}
	state.CurrentSyncCommittee = &altair.SyncCommittee{
		Pubkeys:         []phase0.BLSPubKey{{0x15}},
		AggregatePubkey: phase0.BLSPubKey{0x16},
	}
	state.NextSyncCommittee = &altair.SyncCommittee{
		Pubkeys:         []phase0.BLSPubKey{{0x17}},
		AggregatePubkey: phase0.BLSPubKey{0x18},
	}
	state.LatestBlockHash = phase0.Hash32{0x19}
	state.NextWithdrawalIndex = 23
	state.NextWithdrawalValidatorIndex = 24
	state.HistoricalSummaries = []*capella.HistoricalSummary{{
		BlockSummaryRoot: phase0.Root{0x1a},
		StateSummaryRoot: phase0.Root{0x1b},
	}}
	state.DepositRequestsStartIndex = 25
	state.DepositBalanceToConsume = 26
	state.ExitBalanceToConsume = 27
	state.EarliestExitEpoch = 28
	state.ConsolidationBalanceToConsume = 29
	state.EarliestConsolidationEpoch = 30
	state.PendingDeposits = []*electra.PendingDeposit{{
		Pubkey:                phase0.BLSPubKey{0x1c},
		WithdrawalCredentials: make([]byte, 32),
		Amount:                31,
		Signature:             phase0.BLSSignature{0x1d},
		Slot:                  32,
	}}
	state.PendingPartialWithdrawals = []*electra.PendingPartialWithdrawal{{
		ValidatorIndex:    33,
		Amount:            34,
		WithdrawableEpoch: 35,
	}}
	state.PendingConsolidations = []*electra.PendingConsolidation{{
		SourceIndex: 36,
		TargetIndex: 37,
	}}
	state.ProposerLookahead = []phase0.ValidatorIndex{38, 39}
	state.Builders = []*gloas.Builder{{
		PublicKey:         phase0.BLSPubKey{0x1e},
		Version:           1,
		ExecutionAddress:  bellatrix.ExecutionAddress{0x1f},
		Balance:           40,
		DepositEpoch:      41,
		WithdrawableEpoch: 42,
	}}
	state.NextWithdrawalBuilderIndex = 43
	state.ExecutionPayloadAvailability = []uint8{0xbf, 0xfb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	state.BuilderPendingPayments = []*gloas.BuilderPendingPayment{{
		Weight: 44,
		Withdrawal: &gloas.BuilderPendingWithdrawal{
			FeeRecipient: bellatrix.ExecutionAddress{0x20},
			Amount:       45,
			BuilderIndex: 46,
		},
		ProposerIndex: 47,
	}}
	state.BuilderPendingWithdrawals = []*gloas.BuilderPendingWithdrawal{{
		FeeRecipient: bellatrix.ExecutionAddress{0x21},
		Amount:       48,
		BuilderIndex: 49,
	}}
	state.LatestExecutionPayloadBid = &gloas.ExecutionPayloadBid{
		ParentBlockHash:       phase0.Hash32{0x22},
		ParentBlockRoot:       phase0.Root{0x23},
		BlockHash:             phase0.Hash32{0x24},
		FeeRecipient:          bellatrix.ExecutionAddress{0x25},
		GasLimit:              50,
		BuilderIndex:          51,
		Slot:                  52,
		Value:                 53,
		PrevRandao:            phase0.Root{0x26},
		ExecutionPayment:      57,
		BlobKZGCommitments:    []deneb.KZGCommitment{{0x28}},
		ExecutionRequestsRoot: phase0.Root{0x29},
	}
	state.PayloadExpectedWithdrawals = []*capella.Withdrawal{{
		Index:          54,
		ValidatorIndex: 55,
		Address:        bellatrix.ExecutionAddress{0x27},
		Amount:         56,
	}}
	state.PTCWindow = [][]phase0.ValidatorIndex{{57, 58}, {59, 60}}

	return state
}

// rawField returns the verbatim JSON value the given key carries in input,
// which is what lets a test assert a field's wire shape and not merely its value.
func rawField(t *testing.T, input, key string) string {
	t.Helper()

	var doc map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(input), &doc))
	require.Contains(t, doc, key)

	return string(doc[key])
}

// TestBeaconStateExecutionPayloadAvailabilityJSON verifies that
// execution_payload_availability is carried as a single 0x-prefixed hex string in
// both directions. beacon-APIs types/primitive.yaml defines Bitvector as
// {type: string, format: hex, pattern "^0x[a-fA-F0-9]{2,}$"}, and the consensus
// spec types this field Bitvector[SLOTS_PER_HISTORICAL_ROOT].
//
// Pinning the marshalled field against the literal a node actually sent is what
// makes decoding our own output equivalent to decoding the node's. Without that
// anchor a codec wrong the same way in both directions would satisfy the decode
// half of this test.
func TestBeaconStateExecutionPayloadAvailabilityJSON(t *testing.T) {
	// The value observed from a live Gloas node, whose minimal preset makes
	// SLOTS_PER_HISTORICAL_ROOT 64 and so the bitvector 8 bytes wide.
	const availability = `"0xbffbffffffffffff"`

	state := populatedBeaconState()
	input := mustMarshal(t, state)
	require.Equal(t, availability, rawField(t, input, "execution_payload_availability"))

	var decoded gloas.BeaconState
	require.NoError(t, json.Unmarshal([]byte(input), &decoded))
	require.Equal(t, state.ExecutionPayloadAvailability, decoded.ExecutionPayloadAvailability)
}

// TestBeaconStateExecutionPayloadAvailabilityYAML verifies that the bitvector
// survives a YAML round trip. UnmarshalYAML deliberately routes through
// beaconStateJSON rather than maintaining a second unmarshaler, so the JSON struct's
// field types govern the YAML read path too. That coupling is what made this defect
// break YAML as well: beaconstate_yaml.go already emitted the field as one hex
// string, which then met a []string and failed with "string was used where sequence
// is expected".
//
// Whole states are not compared because other fields have their own YAML gaps,
// tracked separately.
func TestBeaconStateExecutionPayloadAvailabilityYAML(t *testing.T) {
	state := populatedBeaconState()

	data, err := state.MarshalYAML()
	require.NoError(t, err)

	var decoded gloas.BeaconState
	require.NoError(t, decoded.UnmarshalYAML(data))
	require.Equal(t, state.ExecutionPayloadAvailability, decoded.ExecutionPayloadAvailability)
}

// TestBeaconStatePTCWindowJSONShapes verifies that both shapes a node may serve for
// ptc_window decode to the same window.
//
// The spec types the field Vector[Vector[ValidatorIndex, PTC_SIZE], (2 +
// MIN_SEED_LOOKAHEAD) * SLOTS_PER_EPOCH], whose JSON is nested arrays of decimal
// strings. Prysm instead wraps each inner vector in an object carrying
// validator_indices, because proto3 cannot express a nested repeated field without an
// intermediate message, so a Gloas state fetched from Prysm is rejected outright by a
// decoder that admits only the spec shape.
//
// Accepting the wrapper on the read path alone is deliberate: TestBeaconStateJSONWireShapes
// pins what we emit, so the deviation is tolerated where it arrives and never propagated.
func TestBeaconStatePTCWindowJSONShapes(t *testing.T) {
	tests := []struct {
		name      string
		ptcWindow string
		err       string
	}{
		{
			name:      "Spec",
			ptcWindow: `[["57","58"],["59","60"]]`,
		},
		{
			// Observed from prysm on a Gloas devnet.
			name:      "PrysmObjectWrapper",
			ptcWindow: `[{"validator_indices":["57","58"]},{"validator_indices":["59","60"]}]`,
		},
		{
			// Neither shape: the spec shape's error is the one reported, so accepting
			// a second shape does not cost the familiar diagnostic.
			name:      "NeitherShape",
			ptcWindow: `{"validator_indices":["57"]}`,
			err:       "ptc_window: json: cannot unmarshal object into Go value of type [][]string",
		},
		{
			// The wrapper is a shape, not an escape from validation.
			name:      "PrysmObjectWrapperBadIndex",
			ptcWindow: `[{"validator_indices":["fifty-seven"]}]`,
			err:       `ptc_window[0][0]: strconv.ParseUint: parsing "fifty-seven": invalid syntax`,
		},
	}

	state := populatedBeaconState()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var doc map[string]json.RawMessage
			require.NoError(t, json.Unmarshal([]byte(mustMarshal(t, state)), &doc))
			doc["ptc_window"] = json.RawMessage(test.ptcWindow)
			input, err := json.Marshal(doc)
			require.NoError(t, err)

			var decoded gloas.BeaconState
			err = json.Unmarshal(input, &decoded)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}
			require.NoError(t, err)
			require.Equal(t, [][]phase0.ValidatorIndex{{57, 58}, {59, 60}}, decoded.PTCWindow)
		})
	}
}

// TestBeaconStateJSONWireShapes pins the wire shape of the BeaconState fields
// whose SSZ type does not map onto JSON in the way their Go type suggests. These
// are the fields a round-trip cannot vouch for: a codec that is wrong the same
// way in both directions round-trips perfectly and still rejects a real node's
// JSON, so the shapes here are taken from the spec rather than from our output.
func TestBeaconStateJSONWireShapes(t *testing.T) {
	doc := mustMarshal(t, populatedBeaconState())

	// beacon-APIs types/primitive.yaml defines Bitvector as {type: string, format:
	// hex, pattern "^0x[a-fA-F0-9]{2,}$"}. Both of BeaconState's bitvectors are
	// therefore a single string: justification_bits is
	// Bitvector[JUSTIFICATION_BITS_LENGTH] and execution_payload_availability is
	// Bitvector[SLOTS_PER_HISTORICAL_ROOT]. Their Go types differ —
	// bitfield.Bitvector4 against a plain []uint8 — but their wire shape does not.
	bitvector := regexp.MustCompile(`^0x[a-fA-F0-9]{2,}$`)
	for _, field := range []string{"justification_bits", "execution_payload_availability"} {
		t.Run(field, func(t *testing.T) {
			var value string
			require.NoError(t, json.Unmarshal([]byte(rawField(t, doc, field)), &value))
			require.Regexp(t, bitvector, value)
		})
	}

	// ptc_window is Vector[Vector[ValidatorIndex, PTC_SIZE], (2 +
	// MIN_SEED_LOOKAHEAD) * SLOTS_PER_EPOCH] — nested arrays of decimal strings,
	// uint64 entries rather than packed bits. The decimal-per-entry shape that was
	// wrong for execution_payload_availability is the correct shape here, so this
	// pins it against being "fixed" to match its neighbour.
	var ptcWindow [][]string
	require.NoError(t, json.Unmarshal([]byte(rawField(t, doc, "ptc_window")), &ptcWindow))
	require.Equal(t, [][]string{{"57", "58"}, {"59", "60"}}, ptcWindow)
}
