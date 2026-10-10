// Copyright © 2021, 2026 Attestant Limited.
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

package altair_test

import (
	"bytes"
	"encoding/json"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestSyncCommitteeContributionJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.syncCommitteeContributionJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BeaconBlockRootWrongMissing",
			input: []byte(`{"slot":"1","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "beacon block root missing",
		},
		{
			name:  "BeaconBlockRootWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":true,"subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.beacon_block_root of type string",
		},
		{
			name:  "BeaconBlockRootInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"invalid","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid value for beacon block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BeaconBlockRootShort",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d06448","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "BeaconBlockRootLong",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c1c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "SubcommitteeIndexMissing",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "subcommittee index missing",
		},
		{
			name:  "SubcommitteeIndexWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":true,"aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.subcommittee_index of type string",
		},
		{
			name:  "SubcommitteeIndexInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"-3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid value for subcommittee index: strconv.ParseUint: parsing \"-3\": invalid syntax",
		},
		{
			name:  "AggregationBitsMissing",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "aggregation bits missing",
		},
		{
			name:  "AggregationBitsWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":true,"signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.aggregation_bits of type string",
		},
		{
			name:  "AggregationBitsInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"invalid","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
			err:   "invalid value for aggregation bits: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureMissing",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001"}`),
			err:   "signature missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.signature of type string",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"invalid"}`),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c2"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c1c"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncCommitteeContribution
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSyncCommitteeContributionYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{slot: 1, beacon_block_root: '0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c', subcommittee_index: 3, aggregation_bits: '0x0004000000000000000000000000000001', signature: '0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncCommitteeContribution
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(rt), res.String())
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSyncCommitteeContributionSSZ(t *testing.T) {
	t.Run("Static", func(t *testing.T) {
		contribution := &altair.SyncCommitteeContribution{
			Slot:              phase0.Slot(1024),
			BeaconBlockRoot:   phase0.Root{0x01, 0x02, 0x03, 0x04},
			SubcommitteeIndex: 3,
			AggregationBits:   bitfield.Bitvector128{0xaa, 0xbb, 0xcc, 0xdd, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c},
			Signature:         phase0.BLSSignature{0x99},
		}

		encoded, err := contribution.MarshalSSZ()
		require.NoError(t, err)
		require.Equal(t, contribution.SizeSSZ(), len(encoded))
		require.Len(t, encoded, 160)

		var decoded altair.SyncCommitteeContribution
		require.NoError(t, decoded.UnmarshalSSZ(encoded))
		require.Equal(t, contribution.Slot, decoded.Slot)
		require.Equal(t, contribution.BeaconBlockRoot, decoded.BeaconBlockRoot)
		require.Equal(t, contribution.SubcommitteeIndex, decoded.SubcommitteeIndex)
		require.Equal(t, contribution.AggregationBits, decoded.AggregationBits)
		require.Equal(t, contribution.Signature, decoded.Signature)

		root, err := contribution.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)

		decodedRoot, err := decoded.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, root, decodedRoot)

		// Short buffer.
		require.Error(t, decoded.UnmarshalSSZ(encoded[:len(encoded)-1]))

		// Trailing data.
		longBuf := append(append([]byte{}, encoded...), 0x00)
		require.Error(t, decoded.UnmarshalSSZ(longBuf))
	})

	t.Run("DynamicStandardSpec", func(t *testing.T) {
		dynamicSSZ := dynssz.NewDynSsz(map[string]any{
			"SYNC_COMMITTEE_SIZE":         uint64(512),
			"SYNC_COMMITTEE_SUBNET_COUNT": uint64(4),
		})

		contribution := &altair.SyncCommitteeContribution{
			Slot:              phase0.Slot(2048),
			BeaconBlockRoot:   phase0.Root{0x10, 0x20, 0x30},
			SubcommitteeIndex: 2,
			AggregationBits:   bitfield.Bitvector128{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
			Signature:         phase0.BLSSignature{0x42},
		}

		encoded, err := dynamicSSZ.MarshalSSZ(contribution)
		require.NoError(t, err)
		require.Len(t, encoded, 160)

		var decoded altair.SyncCommitteeContribution
		require.NoError(t, dynamicSSZ.UnmarshalSSZ(&decoded, encoded))
		require.Equal(t, contribution.AggregationBits, decoded.AggregationBits)

		staticRoot, err := contribution.HashTreeRoot()
		require.NoError(t, err)

		dynRoot, err := dynamicSSZ.HashTreeRoot(contribution)
		require.NoError(t, err)
		require.Equal(t, staticRoot, dynRoot)
	})

	t.Run("DynamicCustomSpecSubnetCount8", func(t *testing.T) {
		// When SYNC_COMMITTEE_SUBNET_COUNT is 8, AggregationBits is 512 / 8 / 8 = 8 bytes.
		// Total SSZ size is 8 + 32 + 8 + 8 + 96 = 152 bytes.
		dynamicSSZ := dynssz.NewDynSsz(map[string]any{
			"SYNC_COMMITTEE_SIZE":         uint64(512),
			"SYNC_COMMITTEE_SUBNET_COUNT": uint64(8),
		})

		contribution := &altair.SyncCommitteeContribution{
			Slot:              phase0.Slot(4096),
			BeaconBlockRoot:   phase0.Root{0xde, 0xad, 0xbe, 0xef},
			SubcommitteeIndex: 1,
			AggregationBits:   bitfield.Bitvector128{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			Signature:         phase0.BLSSignature{0x55},
		}

		encoded, err := dynamicSSZ.MarshalSSZ(contribution)
		require.NoError(t, err)
		require.Len(t, encoded, 152)

		var decoded altair.SyncCommitteeContribution
		require.NoError(t, dynamicSSZ.UnmarshalSSZ(&decoded, encoded))
		require.Equal(t, contribution.Slot, decoded.Slot)
		require.Equal(t, contribution.BeaconBlockRoot, decoded.BeaconBlockRoot)
		require.Equal(t, contribution.SubcommitteeIndex, decoded.SubcommitteeIndex)
		require.Equal(t, []byte(contribution.AggregationBits), []byte(decoded.AggregationBits))
		require.Equal(t, contribution.Signature, decoded.Signature)

		root, err := dynamicSSZ.HashTreeRoot(contribution)
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)

		decodedRoot, err := dynamicSSZ.HashTreeRoot(&decoded)
		require.NoError(t, err)
		require.Equal(t, root, decodedRoot)

		// Static unmarshaler expects standard 160 bytes and fails on 152 bytes.
		var staticDecoded altair.SyncCommitteeContribution
		require.Error(t, staticDecoded.UnmarshalSSZ(encoded))
	})

	t.Run("DynamicCustomSpecSubnetCount2", func(t *testing.T) {
		// When SYNC_COMMITTEE_SUBNET_COUNT is 2, AggregationBits is 512 / 2 / 8 = 32 bytes.
		// Total SSZ size is 8 + 32 + 8 + 32 + 96 = 176 bytes.
		dynamicSSZ := dynssz.NewDynSsz(map[string]any{
			"SYNC_COMMITTEE_SIZE":         uint64(512),
			"SYNC_COMMITTEE_SUBNET_COUNT": uint64(2),
		})

		aggBits := make(bitfield.Bitvector128, 32)
		for i := range aggBits {
			aggBits[i] = byte(i + 1)
		}

		contribution := &altair.SyncCommitteeContribution{
			Slot:              phase0.Slot(8192),
			BeaconBlockRoot:   phase0.Root{0xaa, 0xbb},
			SubcommitteeIndex: 0,
			AggregationBits:   aggBits,
			Signature:         phase0.BLSSignature{0x77},
		}

		encoded, err := dynamicSSZ.MarshalSSZ(contribution)
		require.NoError(t, err)
		require.Len(t, encoded, 176)

		var decoded altair.SyncCommitteeContribution
		require.NoError(t, dynamicSSZ.UnmarshalSSZ(&decoded, encoded))
		require.Equal(t, contribution.Slot, decoded.Slot)
		require.Equal(t, contribution.BeaconBlockRoot, decoded.BeaconBlockRoot)
		require.Equal(t, contribution.SubcommitteeIndex, decoded.SubcommitteeIndex)
		require.Equal(t, []byte(contribution.AggregationBits), []byte(decoded.AggregationBits))
		require.Equal(t, contribution.Signature, decoded.Signature)

		root, err := dynamicSSZ.HashTreeRoot(contribution)
		require.NoError(t, err)
		require.NotEqual(t, [32]byte{}, root)

		decodedRoot, err := dynamicSSZ.HashTreeRoot(&decoded)
		require.NoError(t, err)
		require.Equal(t, root, decodedRoot)

		// Static unmarshaler expects standard 160 bytes and fails on 176 bytes.
		var staticDecoded altair.SyncCommitteeContribution
		require.Error(t, staticDecoded.UnmarshalSSZ(encoded))
	})
}

