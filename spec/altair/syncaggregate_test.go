// Copyright © 2021 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncAggregateJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.syncAggregateJSON",
		},
		{
			name:  "SyncCommitteeBitsMissing",
			input: []byte(`{"sync_committee_signature":"0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "sync committee bits missing",
		},
		{
			name:  "SyncCommitteeBitsInvalid",
			input: []byte(`{"sync_committee_bits":"invalid","sync_committee_signature":"0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "invalid value for sync committee bits: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SyncCommitteeBitsShort",
			input: []byte(`{"sync_committee_bits":"0xfcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "sync committee bits too short",
		},
		{
			name:  "SyncCommitteeBitsLong",
			input: []byte(`{"sync_committee_bits":"0xe7e7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "sync committee bits too long",
		},
		{
			name:  "SyncCommitteeSignatureMissing",
			input: []byte(`{"sync_committee_bits":"0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec"}`),
			err:   "sync committee signature missing",
		},
		{
			name:  "SyncCommitteeSignatureInvalid",
			input: []byte(`{"sync_committee_bits":"0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"invalid"}`),
			err:   "invalid value for sync committee signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(`{"sync_committee_bits":"0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"0x3b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "sync committee signature short",
		},
		{
			name:  "SignatureLong",
			input: []byte(`{"sync_committee_bits":"0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"0xe6e63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
			err:   "sync committee signature long",
		},
		{
			name:  "Good",
			input: []byte(`{"sync_committee_bits":"0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec","sync_committee_signature":"0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncAggregate
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

func TestSyncAggregateYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{sync_committee_bits: '0xe7fcbc21f184b9b89bfc57cc07232a4fce8e12efee3a8c4967932491267a215cd0aff3e79f19645d6f832592f93d91271071a4e911d3f64447e1f6f68247fdec', sync_committee_signature: '0xe63b8ab602266593dbfe7f714891c5fed225e09c214bda8281c86ceddb6ee10727a854f213d33be1f032399e0044db6fa30368b6dc857fa8f12f61fc3bf4113a6e9cefeb11758fb01a9939950e127d71dc9c54a26aec63ef024b6620e6d32e44'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncAggregate
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
