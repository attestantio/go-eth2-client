// Copyright © 2020 Attestant Limited.
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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestGenesisJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.genesisJSON",
		},
		{
			name:  "GenesisTimeMissing",
			input: []byte(`{"genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
			err:   "genesis time missing",
		},
		{
			name:  "GenesisTimeWrongType",
			input: []byte(`{"genesis_time":true,"genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field genesisJSON.genesis_time of type string",
		},
		{
			name:  "GenesisTimeInvalid",
			input: []byte(`{"genesis_time":"invalid","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
			err:   "invalid value for genesis time: strconv.ParseInt: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "GenesisValidatorsRootMissing",
			input: []byte(`{"genesis_time":"1596546008","genesis_fork_version":"0x00000001"}`),
			err:   "genesis validators root missing",
		},
		{
			name:  "GenesisValidatorsRootWrongType",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":true,"genesis_fork_version":"0x00000001"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field genesisJSON.genesis_validators_root of type string",
		},
		{
			name:  "GenesisValidatorsRootInvalid",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"invalid","genesis_fork_version":"0x00000001"}`),
			err:   "invalid value for genesis validators root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "GenesisValidatorsRootShort",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
			err:   "incorrect length 31 for genesis validators root",
		},
		{
			name:  "GenesisValidatorsRootLong",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x0404700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
			err:   "incorrect length 33 for genesis validators root",
		},
		{
			name:  "GenesisForkVersionMissing",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673"}`),
			err:   "genesis fork version missing",
		},
		{
			name:  "GenesisForkVersionWrongType",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field genesisJSON.genesis_fork_version of type string",
		},
		{
			name:  "GenesisForkVersionInvalid",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"invalid"}`),
			err:   "invalid value for genesis fork version: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "GenesisForkVersionShort",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x000001"}`),
			err:   "incorrect length 3 for genesis fork version",
		},
		{
			name:  "GenesisForkVersionLong",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x0000000001"}`),
			err:   "incorrect length 5 for genesis fork version",
		},
		{
			name:  "Good",
			input: []byte(`{"genesis_time":"1596546008","genesis_validators_root":"0x04700007fabc8282644aed6d1c7c9e21d38a03a0c4ba193f3afe428824b3a673","genesis_fork_version":"0x00000001"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Genesis
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
