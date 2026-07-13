// Copyright © 2025 Attestant Limited.
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
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// MarshalYAML implements yaml.Marshaler.
func (p *ProposerPreferences) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&proposerPreferencesJSON{
		ProposalSlot:   fmt.Sprintf("%d", p.ProposalSlot),
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
		FeeRecipient:   fmt.Sprintf("%#x", p.FeeRecipient),
		GasLimit:       fmt.Sprintf("%d", p.GasLimit),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *ProposerPreferences) UnmarshalYAML(input []byte) error {
	var data proposerPreferencesJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return p.UnmarshalJSON(marshaled)
}
