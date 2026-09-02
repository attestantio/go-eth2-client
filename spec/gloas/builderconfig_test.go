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
	"fmt"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func validBuilderEntry() *gloas.BuilderEntry {
	return &gloas.BuilderEntry{
		URL: []byte("https://builder.example"),
		Auth: &gloas.SignedBuilderRequestAuth{
			Message:   &gloas.BuilderRequestAuth{Data: []byte{0x01}, Slot: 123},
			Signature: phase0.BLSSignature{0x02},
		},
		BuilderPubkeys:      []phase0.BLSPubKey{},
		MaxExecutionPayment: 4,
		MinBid:              5,
		BuilderBoostFactor:  6,
	}
}

func TestBuilderConfigUnmarshalJSON(t *testing.T) {
	var config gloas.BuilderConfig
	err := json.Unmarshal([]byte(`{
		"min_bid":"0",
		"builder_boost_factor":"100",
		"builders":[{
			"url":"https://builder.example",
			"auth":{"message":{"data":"0x01","slot":"123"},"signature":"0x020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
			"builder_pubkeys":["0x030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"],
			"max_execution_payment":"4",
			"min_bid":"5",
			"builder_boost_factor":"6"
		}]
	}`), &config)
	require.NoError(t, err)
	require.Equal(t, phase0.Gwei(0), config.MinBid)
	require.Equal(t, uint64(100), config.BuilderBoostFactor)
	require.Len(t, config.Builders, 1)
	require.Equal(t, []byte("https://builder.example"), config.Builders[0].URL)
	require.Equal(t, []byte{0x01}, config.Builders[0].Auth.Message.Data)
	require.Equal(t, phase0.Slot(123), config.Builders[0].Auth.Message.Slot)
	require.Equal(t, phase0.Gwei(4), config.Builders[0].MaxExecutionPayment)
	require.Equal(t, phase0.Gwei(5), config.Builders[0].MinBid)
	require.Equal(t, uint64(6), config.Builders[0].BuilderBoostFactor)
}

func TestBuilderConfigUnmarshalJSONRejectsTooManyBuilders(t *testing.T) {
	builders := make([]*gloas.BuilderEntry, 65)
	for i := range builders {
		builders[i] = validBuilderEntry()
	}
	input, err := json.Marshal(&gloas.BuilderConfig{Builders: builders})
	require.NoError(t, err)

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "TooManyBuilders",
			err:  "too many builders",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var config gloas.BuilderConfig
			err := json.Unmarshal(input, &config)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderConfigUnmarshalJSONRequiresBuilders(t *testing.T) {
	var config gloas.BuilderConfig
	err := json.Unmarshal([]byte(`{"min_bid":"0","builder_boost_factor":"0"}`), &config)
	require.EqualError(t, err, "builders missing")
}

func TestSignedBuilderRequestAuthUnmarshalJSONRequiresHexPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name: "UnprefixedSignature",
			input: fmt.Sprintf(`{
				"message":{"data":"0x01","slot":"123"},
				"signature":"%x"
			}`, phase0.BLSSignature{0x02}),
			err: "authorization signature missing 0x prefix",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var auth gloas.SignedBuilderRequestAuth
			err := json.Unmarshal([]byte(test.input), &auth)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderRequestAuthUnmarshalJSONRejectsMissingData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "MissingData",
			input: `{"slot":"123"}`,
			err:   "authorization data missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var auth gloas.BuilderRequestAuth
			err := json.Unmarshal([]byte(test.input), &auth)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderRequestAuthUnmarshalJSONRequiresHexPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "UnprefixedData",
			input: `{"data":"01","slot":"123"}`,
			err:   "authorization data missing 0x prefix",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var auth gloas.BuilderRequestAuth
			err := json.Unmarshal([]byte(test.input), &auth)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderRequestAuthUnmarshalJSONRejectsEmptyData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "EmptyData",
			input: `{"data":"0x","slot":"123"}`,
			err:   "authorization data empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var auth gloas.BuilderRequestAuth
			err := json.Unmarshal([]byte(test.input), &auth)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderRequestAuthUnmarshalJSONRejectsOversizedData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "OversizedData",
			input: fmt.Sprintf(`{"data":"0x%s","slot":"123"}`, strings.Repeat("01", 4097)),
			err:   "authorization data exceeds 4096 bytes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var auth gloas.BuilderRequestAuth
			err := json.Unmarshal([]byte(test.input), &auth)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderEntryUnmarshalJSONRequiresPubkeyHexPrefix(t *testing.T) {
	entry := validBuilderEntry()
	entry.BuilderPubkeys = []phase0.BLSPubKey{{0x03}}
	input, err := json.Marshal(entry)
	require.NoError(t, err)
	input = []byte(strings.Replace(string(input), `"0x03`, `"03`, 1))

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "UnprefixedPubkey",
			err:  "builder public key 0 missing 0x prefix",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var decoded gloas.BuilderEntry
			err := json.Unmarshal(input, &decoded)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderEntryUnmarshalJSONRejectsTooManyPubkeys(t *testing.T) {
	entry := validBuilderEntry()
	entry.BuilderPubkeys = make([]phase0.BLSPubKey, 65)
	input, err := json.Marshal(entry)
	require.NoError(t, err)

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "TooManyPubkeys",
			err:  "too many builder public keys",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var decoded gloas.BuilderEntry
			err := json.Unmarshal(input, &decoded)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderEntryUnmarshalJSONRejectsOversizedURL(t *testing.T) {
	entry := validBuilderEntry()
	entry.URL = []byte(strings.Repeat("a", 2049))
	input, err := json.Marshal(entry)
	require.NoError(t, err)

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "OversizedURL",
			err:  "builder URL exceeds 2048 bytes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var decoded gloas.BuilderEntry
			err := json.Unmarshal(input, &decoded)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderEntryUnmarshalJSONRejectsMissingURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name: "MissingURL",
			input: fmt.Sprintf(`{
				"auth":{"message":{"data":"0x01","slot":"123"},"signature":"%#x"},
				"builder_pubkeys":[],
				"max_execution_payment":"4",
				"min_bid":"5",
				"builder_boost_factor":"6"
			}`, phase0.BLSSignature{0x02}),
			err: "builder URL missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var entry gloas.BuilderEntry
			err := json.Unmarshal([]byte(test.input), &entry)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderEntryUnmarshalJSONRejectsMissingAuth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name: "MissingAuth",
			input: `{
				"url":"https://builder.example",
				"builder_pubkeys":[],
				"max_execution_payment":"4",
				"min_bid":"5",
				"builder_boost_factor":"6"
			}`,
			err: "builder authorization missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var entry gloas.BuilderEntry
			err := json.Unmarshal([]byte(test.input), &entry)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestBuilderConfigYAML(t *testing.T) {
	config := &gloas.BuilderConfig{
		BuilderBoostFactor: 100,
		Builders:           []*gloas.BuilderEntry{},
	}

	data, err := yaml.Marshal(config)
	require.NoError(t, err)
	require.Contains(t, string(data), "min_bid: '0'")
	require.Contains(t, string(data), "builder_boost_factor: '100'")
	require.Contains(t, string(data), "builders: []")
}

func TestBuilderConfigSSZ(t *testing.T) {
	config := &gloas.BuilderConfig{
		MinBid:             0,
		BuilderBoostFactor: 100,
		Builders: []*gloas.BuilderEntry{{
			URL: []byte("https://builder.example"),
			Auth: &gloas.SignedBuilderRequestAuth{
				Message:   &gloas.BuilderRequestAuth{Data: []byte{0x01}, Slot: 123},
				Signature: phase0.BLSSignature{0x02},
			},
			BuilderPubkeys:      []phase0.BLSPubKey{{0x03}},
			MaxExecutionPayment: 4,
			MinBid:              5,
			BuilderBoostFactor:  6,
		}},
	}

	data, err := config.MarshalSSZ()
	require.NoError(t, err)

	var decoded gloas.BuilderConfig
	require.NoError(t, decoded.UnmarshalSSZ(data))
	require.Equal(t, config, &decoded)
}

func TestBuilderConfigJSON(t *testing.T) {
	config := &gloas.BuilderConfig{
		MinBid:             0,
		BuilderBoostFactor: 100,
		Builders: []*gloas.BuilderEntry{{
			URL: []byte("https://builder.example"),
			Auth: &gloas.SignedBuilderRequestAuth{
				Message:   &gloas.BuilderRequestAuth{Data: []byte{0x01}, Slot: 123},
				Signature: phase0.BLSSignature{0x02},
			},
			BuilderPubkeys:      []phase0.BLSPubKey{{0x03}},
			MaxExecutionPayment: 4,
			MinBid:              5,
			BuilderBoostFactor:  6,
		}},
	}

	data, err := json.Marshal(config)
	require.NoError(t, err)
	require.JSONEq(t, fmt.Sprintf(`{
		"min_bid":"0",
		"builder_boost_factor":"100",
		"builders":[{
			"url":"https://builder.example",
			"auth":{"message":{"data":"0x01","slot":"123"},"signature":"%#x"},
			"builder_pubkeys":["0x030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"],
			"max_execution_payment":"4",
			"min_bid":"5",
			"builder_boost_factor":"6"
		}]
	}`, config.Builders[0].Auth.Signature), string(data))
}
