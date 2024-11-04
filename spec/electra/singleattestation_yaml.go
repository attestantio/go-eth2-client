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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// singleAttestationYAML is the spec representation of the struct.
type singleAttestationYAML struct {
	CommitteeIndex string                  `yaml:"committee_index"`
	AttesterIndex  string                  `yaml:"attester_index"`
	Data           *phase0.AttestationData `yaml:"data"`
	Signature      string                  `yaml:"signature"`
}

// MarshalYAML implements yaml.Marshaller.
func (a *SingleAttestation) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&singleAttestationYAML{
		CommitteeIndex: fmt.Sprintf("%d", a.CommitteeIndex),
		AttesterIndex:  fmt.Sprintf("%d", a.AttesterIndex),
		Data:           a.Data,
		Signature:      fmt.Sprintf("%#x", a.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *SingleAttestation) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var singleAttestationJSON singleAttestationJSON
	if err := yaml.Unmarshal(input, &singleAttestationJSON); err != nil {
		return err
	}

	return a.unpack(&singleAttestationJSON)
}
