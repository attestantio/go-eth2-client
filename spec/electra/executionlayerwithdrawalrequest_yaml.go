// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"bytes"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// executionLayerWithdrawalRequestYAML is the spec representation of the struct.
type executionLayerWithdrawalRequestYAML struct {
	SourceAddress   bellatrix.ExecutionAddress `yaml:"source_address"`
	ValidatorPubkey phase0.BLSPubKey           `yaml:"validator_pubkey"`
	Amount          phase0.Gwei                `yaml:"amount"`
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionLayerWithdrawalRequest) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&executionLayerWithdrawalRequestYAML{
		SourceAddress:   e.SourceAddress,
		ValidatorPubkey: e.ValidatorPubkey,
		Amount:          e.Amount,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionLayerWithdrawalRequest) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var data executionLayerWithdrawalRequestJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	bytes, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(bytes)
}
