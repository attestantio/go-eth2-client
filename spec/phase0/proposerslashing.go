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

package phase0

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ProposerSlashing provides information about a proposer slashing.
type ProposerSlashing struct {
	SignedHeader1 *SignedBeaconBlockHeader
	SignedHeader2 *SignedBeaconBlockHeader
}

// proposerSlashingJSON is the spec representation of the struct.
type proposerSlashingJSON struct {
	SignedHeader1 *SignedBeaconBlockHeader `json:"signed_header_1"`
	SignedHeader2 *SignedBeaconBlockHeader `json:"signed_header_2"`
}

// proposerSlashingYAML is the spec representation of the struct.
type proposerSlashingYAML struct {
	SignedHeader1 *SignedBeaconBlockHeader `yaml:"signed_header_1"`
	SignedHeader2 *SignedBeaconBlockHeader `yaml:"signed_header_2"`
}

// MarshalJSON implements json.Marshaler.
func (p *ProposerSlashing) MarshalJSON() ([]byte, error) {
	return json.Marshal(&proposerSlashingJSON{
		SignedHeader1: p.SignedHeader1,
		SignedHeader2: p.SignedHeader2,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *ProposerSlashing) UnmarshalJSON(input []byte) error {
	var proposerSlashingJSON proposerSlashingJSON
	if err := json.Unmarshal(input, &proposerSlashingJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&proposerSlashingJSON)
}

func (p *ProposerSlashing) unpack(proposerSlashingJSON *proposerSlashingJSON) error {
	if proposerSlashingJSON.SignedHeader1 == nil {
		return errors.New("signed header 1 missing")
	}
	p.SignedHeader1 = proposerSlashingJSON.SignedHeader1
	if proposerSlashingJSON.SignedHeader2 == nil {
		return errors.New("signed header 2 missing")
	}
	p.SignedHeader2 = proposerSlashingJSON.SignedHeader2

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (p *ProposerSlashing) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&proposerSlashingYAML{
		SignedHeader1: p.SignedHeader1,
		SignedHeader2: p.SignedHeader2,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *ProposerSlashing) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var proposerSlashingJSON proposerSlashingJSON
	if err := yaml.Unmarshal(input, &proposerSlashingJSON); err != nil {
		return err
	}

	return p.unpack(&proposerSlashingJSON)
}

// String returns a string version of the structure.
func (p *ProposerSlashing) String() string {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
