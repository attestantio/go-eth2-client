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

// consolidationYAML is the spec representation of the struct.
type consolidationYAML struct {
	SourceIndex phase0.ValidatorIndex `yaml:"source_index"`
	TargetIndex phase0.ValidatorIndex `yaml:"target_index"`
	Epoch       phase0.Epoch          `yaml:"epoch"`
}

// MarshalYAML implements yaml.Marshaler.
func (c *Consolidation) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&consolidationYAML{
		SourceIndex: c.SourceIndex,
		TargetIndex: c.TargetIndex,
		Epoch:       c.Epoch,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (c *Consolidation) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var unmarshaled consolidationJSON
	if err := yaml.Unmarshal(input, &unmarshaled); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&unmarshaled)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return c.UnmarshalJSON(marshaled)
}