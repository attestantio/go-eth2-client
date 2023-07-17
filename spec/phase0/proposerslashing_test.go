// Copyright © 2020, 2021 Attestant Limited.
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

package phase0_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProposerSlashingJSON(t *testing.T) {
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
			name:  "BadJSON",
			input: []byte(`[]`),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.proposerSlashingJSON",
		},
		{
			name:  "SignedHeader1Missing",
			input: []byte(`{"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
			err:   "signed header 1 missing",
		},
		{
			name:  "SignedHeader1WrongType",
			input: []byte(`{"signed_header_1":true,"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.signedBeaconBlockHeaderJSON",
		},
		{
			name:  "SignedHeader1Invalid",
			input: []byte(`{"signed_header_1":{},"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
			err:   "invalid JSON: message missing",
		},
		{
			name:  "SignedHeader2Missing",
			input: []byte(`{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
			err:   "signed header 2 missing",
		},
		{
			name:  "SignedHeader2WrongType",
			input: []byte(`{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"signed_header_2":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.signedBeaconBlockHeaderJSON",
		},
		{
			name:  "SignedHeader2Invalid",
			input: []byte(`{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"signed_header_2":{}}`),
			err:   "invalid JSON: message missing",
		},
		{
			name:  "Good",
			input: []byte(`{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ProposerSlashing
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

func TestProposerSlashingYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{signed_header_1: {message: {slot: 1, proposer_index: 2, parent_root: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f', state_root: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f', body_root: '0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f'}, signature: '0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf'}, signed_header_2: {message: {slot: 1, proposer_index: 2, parent_root: '0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f', state_root: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f', body_root: '0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f'}, signature: '0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf'}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ProposerSlashing
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
