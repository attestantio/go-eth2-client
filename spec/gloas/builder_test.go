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
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
)

// testBuilder returns a builder whose every field carries a distinctive value, so
// a field dropped by one of the four hand-written codec literals shows up as a
// difference rather than as a coincidentally-matching zero.
func testBuilder() *gloas.Builder {
	return &gloas.Builder{
		PublicKey:         phase0.BLSPubKey{0x01},
		Version:           1,
		ExecutionAddress:  bellatrix.ExecutionAddress{0x02},
		Balance:           3,
		DepositEpoch:      4,
		WithdrawableEpoch: 5,
	}
}

// TestBuilderVersionUnmarshal verifies that a builder's version survives a JSON
// decode. beacon-APIs types/gloas/builder.yaml lists version among Builder's
// required properties, and the SSZ codec carries it, so a JSON decode that drops
// it hands the caller a builder that silently claims version 0.
//
// The field is decoded into the intermediate JSON struct but was never copied
// across to the target, so every builder read over JSON lost it.
func TestBuilderVersionUnmarshal(t *testing.T) {
	builder := testBuilder()

	var decoded gloas.Builder
	require.NoError(t, json.Unmarshal([]byte(mustMarshal(t, builder)), &decoded))
	require.Equal(t, builder, &decoded)
}

// TestBuilderVersionYAML verifies that a builder's version survives a YAML
// round-trip. The read direction is shared — UnmarshalYAML routes through the
// same unpack helper UnmarshalJSON uses — but MarshalJSON and MarshalYAML are two
// independent composite literals, so only the JSON one was ever populated. The
// consensus-spec vectors would catch it, but they are absent, so this is the only
// guard.
//
// YAML writes version as a bare integer. JSON quotes it, as beacon-APIs types
// it Uint8 (see TestBuilderJSONVersionIsString).
func TestBuilderVersionYAML(t *testing.T) {
	builder := testBuilder()

	data, err := builder.MarshalYAML()
	require.NoError(t, err)
	require.Contains(t, string(data), "version: 1")

	var decoded gloas.Builder
	require.NoError(t, decoded.UnmarshalYAML(data))
	require.Equal(t, builder, &decoded)
}

// TestBuilderJSONVersionIsString verifies the JSON form of version. beacon-APIs
// types/gloas/builder.yaml types it as primitive Uint8, a quoted decimal string,
// and Prysm, Lighthouse and Teku all send "version":"0".
func TestBuilderJSONVersionIsString(t *testing.T) {
	input := []byte(`{"pubkey":"0x863917c482adda67fba8305a8ef2cce39cfcd0d13311c86c0cb5b648102b6a481d0b25931d7cfbb5fa3e21c114ec0951","version":"0","execution_address":"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91","balance":"39550000000","deposit_epoch":"1434","withdrawable_epoch":"18446744073709551615"}`)

	var builder gloas.Builder
	require.NoError(t, json.Unmarshal(input, &builder))
	require.Equal(t, uint8(0), builder.Version)
	require.Equal(t, phase0.Epoch(18446744073709551615), builder.WithdrawableEpoch)

	rt, err := json.Marshal(&builder)
	require.NoError(t, err)
	require.Contains(t, string(rt), `"version":"0"`)

	var decoded gloas.Builder
	require.NoError(t, json.Unmarshal(rt, &decoded))
	require.Equal(t, builder, decoded)
}

// builderJSONInput and builderYAMLInput are complete builders that the error
// tests below vary one field at a time.
const (
	builderJSONInput = `{"pubkey":"0x863917c482adda67fba8305a8ef2cce39cfcd0d13311c86c0cb5b648102b6a481d0b25931d7cfbb5fa3e21c114ec0951","version":"0","execution_address":"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91","balance":"1","deposit_epoch":"1","withdrawable_epoch":"1"}`
	builderYAMLInput = `{pubkey: '0x863917c482adda67fba8305a8ef2cce39cfcd0d13311c86c0cb5b648102b6a481d0b25931d7cfbb5fa3e21c114ec0951', version: 7, execution_address: '0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91', balance: 1, deposit_epoch: 2, withdrawable_epoch: 3}`
)

func TestBuilderJSONErrors(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		err  string
	}{
		{name: "VersionMissing", old: `"version":"0",`, err: "version missing"},
		{name: "VersionNumber", old: `"version":"0"`, new: `"version":0`, err: "invalid JSON: json: cannot unmarshal number into Go struct field builderJSON.version of type string"},
		{name: "VersionMax", old: `"version":"0"`, new: `"version":"255"`},
		{name: "VersionOverflow", old: `"version":"0"`, new: `"version":"256"`, err: "invalid value for version: strconv.ParseUint: parsing \"256\": value out of range"},
		{name: "ExecutionAddressShort", old: `"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91"`, new: `"0x6756"`, err: "incorrect length 2 for execution address"},
		{name: "ExecutionAddressLong", old: `"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91"`, new: `"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91aa"`, err: "incorrect length 21 for execution address"},
		{name: "BalanceInvalid", old: `"balance":"1"`, new: `"balance":"abc"`, err: `invalid value for balance: strconv.ParseUint: parsing "abc": invalid syntax`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := strings.Replace(builderJSONInput, test.old, test.new, 1)
			require.NotEqual(t, builderJSONInput, input)
			var builder gloas.Builder
			err := json.Unmarshal([]byte(input), &builder)
			if test.err == "" {
				require.NoError(t, err)
				require.Equal(t, uint8(255), builder.Version)
				out, err := json.Marshal(&builder)
				require.NoError(t, err)
				require.Contains(t, string(out), `"version":"255"`)
				return
			}
			require.EqualError(t, err, test.err)
		})
	}
}

// TestBuilderYAMLUnquotedVersion verifies that YAML with version as a bare
// integer, the form MarshalYAML writes, still decodes.
func TestBuilderYAMLUnquotedVersion(t *testing.T) {
	var builder gloas.Builder
	require.NoError(t, builder.UnmarshalYAML([]byte(builderYAMLInput)))
	require.Equal(t, uint8(7), builder.Version)
}

func TestBuilderYAMLMissingVersion(t *testing.T) {
	input := strings.Replace(builderYAMLInput, "version: 7, ", "", 1)
	require.NotEqual(t, builderYAMLInput, input)
	var builder gloas.Builder
	require.EqualError(t, yaml.Unmarshal([]byte(input), &builder), "version missing")
}
