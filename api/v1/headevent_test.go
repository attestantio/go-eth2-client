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

func TestHeadEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.headEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field headEventJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"slot":"525277","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "block missing",
		},
		{
			name:  "BlockWrongType",
			input: []byte(`{"slot":"525277","block":true,"state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field headEventJSON.block of type string",
		},
		{
			name:  "BlockInvalid",
			input: []byte(`{"slot":"525277","block":"invalid","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid value for block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockShort",
			input: []byte(`{"slot":"525277","block":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "BlockLong",
			input: []byte(`{"slot":"525277","block":"0x9999e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 33 for block",
		},
		{
			name:  "StateMissing",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "state missing",
		},
		{
			name:  "StateWrongType",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":true,"epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field headEventJSON.state of type string",
		},
		{
			name:  "StateInvalid",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"invalid","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid value for state: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "StateShort",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x9a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 31 for state",
		},
		{
			name:  "StateLong",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x74749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 33 for state",
		},
		{
			name:  "EpochTransitionWrongType",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":2,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field headEventJSON.epoch_transition of type bool",
		},
		{
			name:  "CurrentDutyDependentRootMissing",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
		},
		{
			name:  "CurrentDutyDependentRootWrongType",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":true,"previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field headEventJSON.current_duty_dependent_root of type string",
		},
		{
			name:  "CurrentDutyDependentRootInvalid",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"invalid","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "invalid value for current duty dependent root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "CurrentDutyDependentRootShort",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x7a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 31 for current duty dependent root",
		},
		{
			name:  "CurrentDutyDependentRootLong",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x90907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 33 for current duty dependent root",
		},
		{
			name:  "PreviousDutyDependentRootMissing",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61"}`),
		},
		{
			name:  "PreviousDutyDependentRootWrongType",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field headEventJSON.previous_duty_dependent_root of type string",
		},
		{
			name:  "PreviousDutyDependentRootInvalid",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"invalid"}`),
			err:   "invalid value for previous duty dependent root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PreviousDutyDependentRootShort",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x5569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 31 for previous duty dependent root",
		},
		{
			name:  "PreviousDutyDependentRootLong",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x93935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
			err:   "incorrect length 33 for previous duty dependent root",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch_transition":false,"current_duty_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","previous_duty_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.HeadEvent
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
