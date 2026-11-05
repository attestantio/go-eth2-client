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
	"bytes"
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

	var config gloas.BuilderConfig
	require.EqualError(t, json.Unmarshal(input, &config), "too many builders")
}

func TestBuilderConfigUnmarshalJSONRequiresBuilders(t *testing.T) {
	var config gloas.BuilderConfig
	err := json.Unmarshal([]byte(`{"min_bid":"0","builder_boost_factor":"0"}`), &config)
	require.EqualError(t, err, "builders missing")
}

func TestSignedBuilderRequestAuthUnmarshalJSONRequiresHexPrefix(t *testing.T) {
	input := fmt.Sprintf(`{
		"message":{"data":"0x01","slot":"123"},
		"signature":"%x"
	}`, phase0.BLSSignature{0x02})

	var auth gloas.SignedBuilderRequestAuth
	require.EqualError(t, json.Unmarshal([]byte(input), &auth), "authorization signature missing 0x prefix")
}

func TestBuilderRequestAuthUnmarshalJSONRejectsBadData(t *testing.T) {
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
		{
			name:  "UnprefixedData",
			input: `{"data":"01","slot":"123"}`,
			err:   "authorization data missing 0x prefix",
		},
		{
			name:  "EmptyData",
			input: `{"data":"0x","slot":"123"}`,
			err:   "authorization data empty",
		},
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

	var decoded gloas.BuilderEntry
	require.EqualError(t, json.Unmarshal(input, &decoded), "builder public key 0 missing 0x prefix")
}

func TestBuilderEntryUnmarshalJSONRejectsTooManyPubkeys(t *testing.T) {
	entry := validBuilderEntry()
	entry.BuilderPubkeys = make([]phase0.BLSPubKey, 65)
	input, err := json.Marshal(entry)
	require.NoError(t, err)

	var decoded gloas.BuilderEntry
	require.EqualError(t, json.Unmarshal(input, &decoded), "too many builder public keys")
}

func TestBuilderEntryUnmarshalJSONRejectsOversizedURL(t *testing.T) {
	entry := validBuilderEntry()
	entry.URL = []byte(strings.Repeat("a", 2049))
	input, err := json.Marshal(entry)
	require.NoError(t, err)

	var decoded gloas.BuilderEntry
	require.EqualError(t, json.Unmarshal(input, &decoded), "builder URL exceeds 2048 bytes")
}

func TestBuilderEntryUnmarshalJSONRejectsMissingFields(t *testing.T) {
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

// TestBuilderConfigYAMLWithBuilders round-trips a populated config.  An empty
// builders list exercises none of the entry encoding, so a nested entry has to
// be present for the wire names -- and the URL's "://", which a plain YAML
// scalar cannot carry in flow style -- to be checked at all.
func TestBuilderConfigYAMLWithBuilders(t *testing.T) {
	config := &gloas.BuilderConfig{
		BuilderBoostFactor: 100,
		Builders:           []*gloas.BuilderEntry{validBuilderEntry()},
	}
	config.Builders[0].BuilderPubkeys = []phase0.BLSPubKey{{0x03}}

	data, err := yaml.Marshal(config)
	require.NoError(t, err)
	require.Contains(t, string(data), "url: 'https://builder.example'")
	require.Contains(t, string(data), "max_execution_payment: '4'")
	require.Contains(t, string(data), "builder_pubkeys: ['0x03")

	var decoded gloas.BuilderConfig
	require.NoError(t, yaml.Unmarshal(data, &decoded))
	require.Equal(t, config, &decoded)
}

// TestBuilderConfigNilBuildersEncodeAsEmpty verifies a nil builders slice
// reaches the wire as [], not null: builders is a required array and an empty
// one is a meaningful request, so null would be rejected by the node and by
// this package's own decoder.
func TestBuilderConfigNilBuildersEncodeAsEmpty(t *testing.T) {
	config := &gloas.BuilderConfig{BuilderBoostFactor: 100}

	data, err := json.Marshal(config)
	require.NoError(t, err)
	require.JSONEq(t, `{"min_bid":"0","builder_boost_factor":"100","builders":[]}`, string(data))

	var decoded gloas.BuilderConfig
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Empty(t, decoded.Builders)
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

// TestBuilderAPIContainerHashTreeRoots pins the roots of the builder API
// types, which are regular SSZ containers rather than progressive ones:
// BuilderRequestAuth and SignedBuilderRequestAuth as defined in builder-specs
// at commit 61aeca42d9ba8eb467ae230d191ebc1c3f38107d, and BuilderEntry and
// BuilderConfig as defined in beacon-APIs at commit
// 01e99a6530b22941f619bde525cbf8bb61991d42, and BuilderPreferencesEntry as
// defined in beacon-APIs at commit 159622d983a703eb03a8a37bb1edeab7ffc3b6bc.
// The expected roots were independently calculated, rather than with this
// repository's generated dynamic-ssz codec: with eth-ssz-specs 0.1.0, the SSZ
// library of the consensus-specs PySpec, and for BuilderPreferencesEntry with
// remerkleable 0.1.28, which reproduces the BuilderRequestAuth and
// SignedBuilderRequestAuth roots below.
func TestBuilderAPIContainerHashTreeRoots(t *testing.T) {
	auth := &gloas.SignedBuilderRequestAuth{
		Message: &gloas.BuilderRequestAuth{
			Data: []byte("builder.example.com"),
			Slot: 12345,
		},
		Signature: phase0.BLSSignature{0x01, 0x02, 0x03},
	}
	var pubkeyA, pubkeyB phase0.BLSPubKey
	copy(pubkeyA[:], bytes.Repeat([]byte{0xaa}, len(pubkeyA)))
	copy(pubkeyB[:], bytes.Repeat([]byte{0xbb}, len(pubkeyB)))
	entry := &gloas.BuilderEntry{
		URL:                 []byte("https://builder.example.com"),
		Auth:                auth,
		BuilderPubkeys:      []phase0.BLSPubKey{pubkeyA, pubkeyB},
		MaxExecutionPayment: 7,
		MinBid:              5,
		BuilderBoostFactor:  100,
	}

	tests := []struct {
		name  string
		value interface{ HashTreeRoot() ([32]byte, error) }
		root  string
	}{
		{
			name:  "BuilderRequestAuth",
			value: auth.Message,
			root:  "6ad33c944510322cf6c8f8148a2e60a5b199367c6ef7ab4f0a98133f9c9fd6c9",
		},
		{
			name:  "SignedBuilderRequestAuth",
			value: auth,
			root:  "f15386db9a7a1961f761390f93b14eebbe9c81e29b32cc3307d6fade2fa5dbe0",
		},
		{
			name:  "BuilderEntry",
			value: entry,
			root:  "51483a4d4a92ea225fdb5bf8defc8b75d4fcc6647d7d7f3bafc53045b4154fa4",
		},
		{
			name: "BuilderConfig",
			value: &gloas.BuilderConfig{
				MinBid:             3,
				BuilderBoostFactor: 110,
				Builders:           []*gloas.BuilderEntry{entry, entry},
			},
			root: "5c505a24cf3775610a58ecdaa64a972abe7e1513a92c82d0ced248f4e972ae10",
		},
		{
			name:  "BuilderPreferencesEntry",
			value: validBuilderPreferencesEntry(),
			root:  "895c37789bba97bd876dbee401dbd053eb704e85cb7cbb7d3bc37afbeca8448d",
		},
		{
			name:  "EmptyBuilderConfig",
			value: &gloas.BuilderConfig{BuilderBoostFactor: 100},
			root:  "c5a7f486b2ebcc77b77b49d8714b489535e4435b7a4a3223ecfc480a83e5da91",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := test.value.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, rootFromHex(t, test.root), root)
		})
	}
}
