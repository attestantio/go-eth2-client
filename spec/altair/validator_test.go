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

package altair_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestValidatorJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.validatorJSON",
		},
		{
			name:  "PublicKeyMissing",
			input: []byte(`{"withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"pubkey":true,"withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"pubkey":"invalid","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(`{"pubkey":"0x9bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "incorrect length 47 for public key",
		},
		{
			name:  "PublicKeyLong",
			input: []byte(`{"pubkey":"0xb8b89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "incorrect length 49 for public key",
		},
		{
			name:  "WithdrawalCredentialsMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "withdrawal credentials missing",
		},
		{
			name:  "WithdrawalCredentialsWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":true,"effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.withdrawal_credentials of type string",
		},
		{
			name:  "WithdrawalCredentialsInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"invalid","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for withdrawal credentials: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "WithdrawalCredentialsShort",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0xec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "incorrect length 31 for withdrawal credentials",
		},
		{
			name:  "WithdrawalCredentialsLong",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x0000ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "incorrect length 33 for withdrawal credentials",
		},
		{
			name:  "EffectiveBalanceMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "effective balance missing",
		},
		{
			name:  "EffectiveBalanceWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":true,"slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.effective_balance of type string",
		},
		{
			name:  "EffectiveBalanceInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"-1","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for effective balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SlashedWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":"false","activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal string into Go struct field validatorJSON.slashed of type bool",
		},
		{
			name:  "ActivationEligibilityEpochMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "activation eligibility epoch missing",
		},
		{
			name:  "ActivationEligibilityEpochWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":true,"activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.activation_eligibility_epoch of type string",
		},
		{
			name:  "ActivationEligibilityInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"-1","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for activation eligibility epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ActivationEpochMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "activation epoch missing",
		},
		{
			name:  "ActivationEligibilityEpochWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":true,"exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.activation_epoch of type string",
		},
		{
			name:  "ActivationInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"-1","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for activation epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ExitEpochMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","withdrawable_epoch":"18446744073709551615"}`),
			err:   "exit epoch missing",
		},
		{
			name:  "ExitEligibilityEpochWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":true,"withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.exit_epoch of type string",
		},
		{
			name:  "ExitInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"-1","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for exit epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "WithdrawableEpochMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615"}`),
			err:   "withdrawable epoch missing",
		},
		{
			name:  "WithdrawableEligibilityEpochWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.withdrawable_epoch of type string",
		},
		{
			name:  "WithdrawableInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"-1"}`),
			err:   "invalid value for withdrawable epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.Validator
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

func TestValidatorYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{pubkey: '0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b', withdrawal_credentials: '0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594', effective_balance: 32000000000, slashed: false, activation_eligibility_epoch: 0, activation_epoch: 0, exit_epoch: 18446744073709551615, withdrawable_epoch: 18446744073709551615}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.Validator
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestValidatorSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "altair", "ssz_static", "Validator", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := ioutil.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res spec.Validator
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				specSSZ, err := ioutil.ReadFile(filepath.Join(path, "serialized.ssz"))
				require.NoError(t, err)

				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)

				root, err := res.HashTreeRoot()
				require.NoError(t, err)
				rootsYAML, err := ioutil.ReadFile(filepath.Join(path, "roots.yaml"))
				require.NoError(t, err)
				require.Equal(t, string(rootsYAML), fmt.Sprintf("{root: '%#x'}\n", root))
			})
		}
		return nil
	}))
}
