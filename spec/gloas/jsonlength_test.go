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
	"strings"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	require "github.com/stretchr/testify/require"
)

// hexOfLen returns a 0x-prefixed hex string that decodes to exactly n bytes.
func hexOfLen(n int) string {
	return "0x" + strings.Repeat("ab", n)
}

// mustMarshal renders v as JSON, which the tests use to build a valid document
// for the types whose JSON is too large to spell out by hand.
func mustMarshal(t *testing.T, v any) string {
	t.Helper()

	data, err := json.Marshal(v)
	require.NoError(t, err)

	return string(data)
}

// withHexField returns input with the hex value of the given JSON key replaced,
// making a single field wrong-length while the rest of the document stays valid.
// The key must appear exactly once, so the mutation is unambiguous. Input comes
// from mustMarshal, so the pattern can assume compact JSON and lowercase hex.
func withHexField(t *testing.T, input, key, value string) string {
	t.Helper()

	re := regexp.MustCompile(`"` + key + `":"0x[0-9a-f]*"`)
	require.Len(t, re.FindAllString(input, -1), 1, "key %q must appear exactly once", key)

	return re.ReplaceAllString(input, `"`+key+`":"`+value+`"`)
}

// validExecutionPayloadEnvelope returns an envelope whose nested payload and
// execution requests satisfy the checks preceding the fields under test.
func validExecutionPayloadEnvelope() *gloas.ExecutionPayloadEnvelope {
	return &gloas.ExecutionPayloadEnvelope{
		Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(0)},
		ExecutionRequests: &gloas.ExecutionRequests{},
	}
}

// validBeaconState returns a state populated just far enough for the
// unmarshaler to reach latest_block_hash: every field checked before it must be
// present and well-formed. Fields after it are left at their zero value, which
// is harmless because a wrong-length latest_block_hash short-circuits first.
func validBeaconState() *gloas.BeaconState {
	// One pubkey per committee suffices: altair.SyncCommittee's unpack only
	// requires a non-empty list, and gloas.BeaconState adds no size check. A
	// spec-sized 512 would inflate this fixture from ~2.6 KB to ~106 KB.
	pubkeys := make([]phase0.BLSPubKey, 1)

	return &gloas.BeaconState{
		Fork:                        &phase0.Fork{},
		LatestBlockHeader:           &phase0.BeaconBlockHeader{},
		ETH1Data:                    &phase0.ETH1Data{BlockHash: make([]byte, phase0.Hash32Length)},
		JustificationBits:           bitfield.NewBitvector4(),
		PreviousJustifiedCheckpoint: &phase0.Checkpoint{},
		CurrentJustifiedCheckpoint:  &phase0.Checkpoint{},
		FinalizedCheckpoint:         &phase0.Checkpoint{},
		CurrentSyncCommittee:        &altair.SyncCommittee{Pubkeys: pubkeys},
		NextSyncCommittee:           &altair.SyncCommittee{Pubkeys: pubkeys},
	}
}

// validBeaconBlock returns a block whose body satisfies the nested checks that
// precede a signed block's own signature.
func validBeaconBlock() *gloas.BeaconBlock {
	return &gloas.BeaconBlock{
		Body: &gloas.BeaconBlockBody{
			ETH1Data:                  &phase0.ETH1Data{BlockHash: make([]byte, phase0.Hash32Length)},
			SyncAggregate:             &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			ParentExecutionRequests:   &gloas.ExecutionRequests{},
		},
	}
}

// TestJSONFixedLengthGuards verifies that the hand-written JSON unmarshalers
// reject a hex value whose decoded length does not match the fixed-size array it
// is copied into. hex.DecodeString accepts any even-length input, so without a
// length check a short value silently zero-pads the array's tail and a long one
// is silently truncated, with UnmarshalJSON still returning nil — handing the
// caller a corrupted root, hash, address or signature and surfacing no error.
func TestJSONFixedLengthGuards(t *testing.T) {
	bidJSON := mustMarshal(t, &gloas.ExecutionPayloadBid{})
	envelopeJSON := mustMarshal(t, validExecutionPayloadEnvelope())
	blockJSON := mustMarshal(t, validBeaconBlock())
	stateJSON := mustMarshal(t, validBeaconState())
	dataJSON := mustMarshal(t, &gloas.PayloadAttestationData{})
	prefsJSON := mustMarshal(t, &gloas.ProposerPreferences{})

	tests := []struct {
		name   string
		target json.Unmarshaler
		input  string
		err    string
	}{
		{
			name:   "PayloadAttestationDataBeaconBlockRootShort",
			target: &gloas.PayloadAttestationData{},
			input:  `{"beacon_block_root":"` + hexOfLen(31) + `","slot":"1","payload_present":true,"blob_data_available":true}`,
			err:    "incorrect length for beacon block root",
		},
		{
			name:   "PayloadAttestationDataBeaconBlockRootLong",
			target: &gloas.PayloadAttestationData{},
			input:  `{"beacon_block_root":"` + hexOfLen(33) + `","slot":"1","payload_present":true,"blob_data_available":true}`,
			err:    "incorrect length for beacon block root",
		},
		{
			name:   "PayloadAttestationMessageSignatureShort",
			target: &gloas.PayloadAttestationMessage{},
			input:  `{"validator_index":"1","data":` + dataJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "PayloadAttestationMessageSignatureLong",
			target: &gloas.PayloadAttestationMessage{},
			input:  `{"validator_index":"1","data":` + dataJSON + `,"signature":"` + hexOfLen(97) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "ProposerPreferencesFeeRecipientShort",
			target: &gloas.ProposerPreferences{},
			input:  withHexField(t, prefsJSON, "fee_recipient", hexOfLen(19)),
			err:    "incorrect length for fee recipient",
		},
		{
			name:   "ProposerPreferencesFeeRecipientLong",
			target: &gloas.ProposerPreferences{},
			input:  withHexField(t, prefsJSON, "fee_recipient", hexOfLen(21)),
			err:    "incorrect length for fee recipient",
		},
		{
			name:   "PayloadAttestationSignatureShort",
			target: &gloas.PayloadAttestation{},
			input: `{"aggregation_bits":"` + hexOfLen(64) + `","data":` + dataJSON +
				`,"signature":"` + hexOfLen(95) + `"}`,
			err: "incorrect length for signature",
		},
		{
			name:   "IndexedPayloadAttestationSignatureShort",
			target: &gloas.IndexedPayloadAttestation{},
			input:  `{"attesting_indices":["1"],"data":` + dataJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "SignedProposerPreferencesSignatureShort",
			target: &gloas.SignedProposerPreferences{},
			input:  `{"message":` + prefsJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "SignedExecutionPayloadBidSignatureShort",
			target: &gloas.SignedExecutionPayloadBid{},
			input:  `{"message":` + bidJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "SignedExecutionPayloadEnvelopeSignatureShort",
			target: &gloas.SignedExecutionPayloadEnvelope{},
			input:  `{"message":` + envelopeJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "SignedBeaconBlockSignatureShort",
			target: &gloas.SignedBeaconBlock{},
			input:  `{"message":` + blockJSON + `,"signature":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for signature",
		},
		{
			name:   "ExecutionPayloadEnvelopeBeaconBlockRootShort",
			target: &gloas.ExecutionPayloadEnvelope{},
			input:  withHexField(t, envelopeJSON, "beacon_block_root", hexOfLen(31)),
			err:    "incorrect length for beacon block root",
		},
		{
			name:   "ExecutionPayloadEnvelopeParentBeaconBlockRootShort",
			target: &gloas.ExecutionPayloadEnvelope{},
			input:  withHexField(t, envelopeJSON, "parent_beacon_block_root", hexOfLen(31)),
			err:    "incorrect length for parent beacon block root",
		},
		{
			name:   "ExecutionPayloadBidParentBlockHashShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "parent_block_hash", hexOfLen(31)),
			err:    "incorrect length for parent block hash",
		},
		{
			name:   "ExecutionPayloadBidParentBlockRootShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "parent_block_root", hexOfLen(31)),
			err:    "incorrect length for parent block root",
		},
		{
			name:   "ExecutionPayloadBidBlockHashShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "block_hash", hexOfLen(31)),
			err:    "incorrect length for block hash",
		},
		{
			name:   "ExecutionPayloadBidPrevRandaoShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "prev_randao", hexOfLen(31)),
			err:    "incorrect length for prev randao",
		},
		{
			name:   "ExecutionPayloadBidFeeRecipientShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "fee_recipient", hexOfLen(19)),
			err:    "incorrect length for fee recipient",
		},
		{
			name:   "ExecutionPayloadBidExecutionRequestsRootShort",
			target: &gloas.ExecutionPayloadBid{},
			input:  withHexField(t, bidJSON, "execution_requests_root", hexOfLen(31)),
			err:    "incorrect length for execution requests root",
		},
		{
			// randao_reveal is the body's first check, so no other field is needed.
			name:   "BeaconBlockBodyRANDAORevealShort",
			target: &gloas.BeaconBlockBody{},
			input:  `{"randao_reveal":"` + hexOfLen(95) + `"}`,
			err:    "incorrect length for randao reveal",
		},
		{
			name:   "BeaconBlockBodyGraffitiShort",
			target: &gloas.BeaconBlockBody{},
			input:  `{"randao_reveal":"` + hexOfLen(96) + `","graffiti":"` + hexOfLen(31) + `"}`,
			err:    "incorrect length for graffiti",
		},
		{
			name:   "BeaconStateLatestBlockHashShort",
			target: &gloas.BeaconState{},
			input:  withHexField(t, stateJSON, "latest_block_hash", hexOfLen(31)),
			err:    "incorrect length for latest block hash",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(test.input), test.target)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestJSONFixedLengthRoundTrip is the counterpart to TestJSONFixedLengthGuards:
// it verifies the guards reject only genuinely wrong-length values. A
// correctly-sized document must still unmarshal and re-marshal unchanged, which
// would fail if a guard used the wrong length constant or a one-sided
// comparison.
//
// The BeaconState row does double duty: its fixture populates all 46 fields rather
// than only those preceding a guard, so it also catches any field whose marshal and
// unmarshal disagree, length-guarded or not — which is how it caught gloas.Builder
// dropping its version.
func TestJSONFixedLengthRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		target json.Unmarshaler
	}{
		{
			name:   "PayloadAttestationData",
			value:  &gloas.PayloadAttestationData{},
			target: &gloas.PayloadAttestationData{},
		},
		{
			name:   "PayloadAttestationMessage",
			value:  &gloas.PayloadAttestationMessage{Data: &gloas.PayloadAttestationData{}},
			target: &gloas.PayloadAttestationMessage{},
		},
		{
			name:   "ProposerPreferences",
			value:  &gloas.ProposerPreferences{},
			target: &gloas.ProposerPreferences{},
		},
		{
			name:   "SignedProposerPreferences",
			value:  &gloas.SignedProposerPreferences{Message: &gloas.ProposerPreferences{}},
			target: &gloas.SignedProposerPreferences{},
		},
		{
			name:   "ExecutionPayloadBid",
			value:  &gloas.ExecutionPayloadBid{},
			target: &gloas.ExecutionPayloadBid{},
		},
		{
			name:   "SignedExecutionPayloadBid",
			value:  &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			target: &gloas.SignedExecutionPayloadBid{},
		},
		{
			name:   "ExecutionPayloadEnvelope",
			value:  validExecutionPayloadEnvelope(),
			target: &gloas.ExecutionPayloadEnvelope{},
		},
		{
			name:   "SignedExecutionPayloadEnvelope",
			value:  &gloas.SignedExecutionPayloadEnvelope{Message: validExecutionPayloadEnvelope()},
			target: &gloas.SignedExecutionPayloadEnvelope{},
		},
		{
			// Also exercises the body's randao_reveal and graffiti guards.
			name:   "SignedBeaconBlock",
			value:  &gloas.SignedBeaconBlock{Message: validBeaconBlock()},
			target: &gloas.SignedBeaconBlock{},
		},
		{
			// All 46 fields populated, per the note above. A field wrong the same
			// way in both directions still passes here; TestBeaconStateJSONWireShapes
			// is what catches those.
			name:   "BeaconState",
			value:  populatedBeaconState(),
			target: &gloas.BeaconState{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := mustMarshal(t, test.value)
			require.NoError(t, json.Unmarshal([]byte(input), test.target))
			require.Equal(t, input, mustMarshal(t, test.target))
		})
	}
}
