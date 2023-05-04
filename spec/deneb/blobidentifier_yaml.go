// Copyright © 2023 Attestant Limited.
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

package deneb

import (
	"bytes"

	"github.com/goccy/go-yaml"
)

// blobIdentifierYAML is the spec representation of the struct.
type blobIdentifierYAML struct {
	BlockRoot string `yaml:"block_root"`
	Index     uint64 `yaml:"index"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BlobIdentifier) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blobIdentifierYAML{
		BlockRoot: b.BlockRoot.String(),
		Index:     uint64(b.Index),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BlobIdentifier) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data blobIdentifierJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}
	return b.unpack(&data)
}
