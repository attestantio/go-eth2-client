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

package gloas

import (
	"bytes"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// executionRequestsYAML is the spec representation of the struct.
type executionRequestsYAML struct {
	Deposits        []*electra.DepositRequest       `yaml:"deposits"`
	Withdrawals     []*electra.WithdrawalRequest    `yaml:"withdrawals"`
	Consolidations  []*electra.ConsolidationRequest `yaml:"consolidations"`
	BuilderDeposits []*BuilderDepositRequest        `yaml:"builder_deposits"`
	BuilderExits    []*BuilderExitRequest           `yaml:"builder_exits"`
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionRequests) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&executionRequestsYAML{
		Deposits:        e.Deposits,
		Withdrawals:     e.Withdrawals,
		Consolidations:  e.Consolidations,
		BuilderDeposits: e.BuilderDeposits,
		BuilderExits:    e.BuilderExits,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionRequests) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var unmarshaled executionRequestsJSON
	if err := yaml.Unmarshal(input, &unmarshaled); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}

	marshaled, err := json.Marshal(unmarshaled)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(marshaled)
}
