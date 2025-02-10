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

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// pendingPartialWithdrawalYAML is the spec representation of the struct.
type pendingPartialWithdrawalYAML struct {
	ValidatorIndex    phase0.ValidatorIndex `yaml:"validator_index"`
	Amount            phase0.Gwei           `yaml:"amount"`
	WithdrawableEpoch phase0.Epoch          `yaml:"withdrawable_epoch"`
}

// MarshalYAML implements yaml.Marshaler.
func (p *PendingPartialWithdrawal) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&pendingPartialWithdrawalYAML{
		ValidatorIndex:    p.ValidatorIndex,
		Amount:            p.Amount,
		WithdrawableEpoch: p.WithdrawableEpoch,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *PendingPartialWithdrawal) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var unmarshaled pendingPartialWithdrawalJSON
	if err := yaml.Unmarshal(input, &unmarshaled); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&unmarshaled)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return p.UnmarshalJSON(marshaled)
}
