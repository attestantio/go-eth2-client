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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

// validBuilderPreferencesEntry uses the url and max_execution_payment examples
// from beacon-APIs types/gloas/builder_preferences_entry.yaml.
func validBuilderPreferencesEntry() *gloas.BuilderPreferencesEntry {
	return &gloas.BuilderPreferencesEntry{
		ProposerPubkey: phase0.BLSPubKey{0x01},
		URL:            []byte("https://builder.example.com"),
		Auth: &gloas.SignedBuilderRequestAuth{
			Message:   &gloas.BuilderRequestAuth{Data: []byte("builder.example.com"), Slot: 123},
			Signature: phase0.BLSSignature{0x02},
		},
		MaxExecutionPayment: 1000000000,
	}
}

func TestBuilderPreferencesEntryJSON(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := json.Marshal(entry)
	require.NoError(t, err)
	require.JSONEq(t, fmt.Sprintf(`{
		"proposer_pubkey":"%#x",
		"url":"https://builder.example.com",
		"auth":{"message":{"data":"%#x","slot":"123"},"signature":"%#x"},
		"max_execution_payment":"1000000000"
	}`, entry.ProposerPubkey, entry.Auth.Message.Data, entry.Auth.Signature), string(data))

	var decoded gloas.BuilderPreferencesEntry
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, entry, &decoded)
}

// TestBuilderPreferencesEntryYAML round-trips an entry through YAML.  The URL's
// "://" must survive as a quoted scalar.  The YAML and JSON marshalers each
// build their own representation, so the test checks they carry the same values.
func TestBuilderPreferencesEntryYAML(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := yaml.Marshal(entry)
	require.NoError(t, err)
	require.Contains(t, string(data), "url: 'https://builder.example.com'")
	require.Contains(t, string(data), "max_execution_payment: '1000000000'")
	require.Equal(t, entry.String(), string(data))

	var decoded gloas.BuilderPreferencesEntry
	require.NoError(t, yaml.Unmarshal(data, &decoded))
	require.Equal(t, entry, &decoded)

	jsonData, err := json.Marshal(entry)
	require.NoError(t, err)
	var fromYAML, fromJSON map[string]any
	require.NoError(t, yaml.Unmarshal(data, &fromYAML))
	require.NoError(t, json.Unmarshal(jsonData, &fromJSON))
	require.Equal(t, fromJSON, fromYAML)
}

// TestBuilderPreferencesEntryYAMLQuotes round-trips URLs holding the characters
// that YAML quoting has to escape.
func TestBuilderPreferencesEntryYAMLQuotes(t *testing.T) {
	for _, url := range []string{`https://b.example/?q="x"`, `https://b.example/?q='x'`} {
		entry := validBuilderPreferencesEntry()
		entry.URL = []byte(url)

		data, err := yaml.Marshal(entry)
		require.NoError(t, err)
		var decoded gloas.BuilderPreferencesEntry
		require.NoError(t, yaml.Unmarshal(data, &decoded), url)
		require.Equal(t, entry, &decoded, url)
	}
}

func TestBuilderPreferencesEntryUnmarshalJSONRejectsBadInput(t *testing.T) {
	valid, err := json.Marshal(validBuilderPreferencesEntry())
	require.NoError(t, err)
	without := func(key string) string {
		var fields map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(valid, &fields))
		delete(fields, key)
		data, err := json.Marshal(fields)
		require.NoError(t, err)

		return string(data)
	}
	replace := func(old, replacement string) string {
		require.Contains(t, string(valid), old)

		return strings.Replace(string(valid), old, replacement, 1)
	}

	tests := []struct {
		name  string
		input string
		err   string
	}{
		{name: "MissingPubkey", input: without("proposer_pubkey"), err: "proposer public key missing"},
		{name: "UnprefixedPubkey", input: replace(`"proposer_pubkey":"0x`, `"proposer_pubkey":"`), err: "proposer public key missing 0x prefix"},
		{name: "ShortPubkey", input: replace(`"proposer_pubkey":"0x01`, `"proposer_pubkey":"0x`), err: "incorrect proposer public key length 47"},
		{name: "MissingURL", input: without("url"), err: "builder URL missing"},
		{name: "OversizedURL", input: replace(`"url":"https://`, `"url":"https://`+strings.Repeat("a", 2048)), err: "builder URL exceeds 2048 bytes"},
		{name: "MissingAuth", input: without("auth"), err: "builder authorization missing"},
		{name: "MissingMaxExecutionPayment", input: without("max_execution_payment"), err: "invalid max execution payment: strconv.ParseUint: parsing \"\": invalid syntax"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var entry gloas.BuilderPreferencesEntry
			require.EqualError(t, json.Unmarshal([]byte(test.input), &entry), test.err)
		})
	}
}

func TestBuilderPreferencesEntrySSZ(t *testing.T) {
	entry := validBuilderPreferencesEntry()

	data, err := entry.MarshalSSZ()
	require.NoError(t, err)

	// Field order per the spec container: proposer_pubkey (48), offset of url,
	// offset of auth, max_execution_payment (8), then url and auth.
	urlLen := uint32(len(entry.URL))
	require.Equal(t, entry.ProposerPubkey[:], data[0:48])
	require.Equal(t, uint32(64), binary.LittleEndian.Uint32(data[48:52]))
	require.Equal(t, 64+urlLen, binary.LittleEndian.Uint32(data[52:56]))
	require.Equal(t, uint64(entry.MaxExecutionPayment), binary.LittleEndian.Uint64(data[56:64]))
	require.Equal(t, entry.URL, data[64:64+urlLen])

	var decoded gloas.BuilderPreferencesEntry
	require.NoError(t, decoded.UnmarshalSSZ(data))
	require.Equal(t, entry, &decoded)
}
