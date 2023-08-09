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
	"encoding/json"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// blindedBlockContentsYAML is the spec representation of the struct.
type blindedBlockContentsYAML struct {
	BlindedBlock        *BlindedBeaconBlock   `yaml:"blinded_block"`
	BlindedBlobSidecars []*BlindedBlobSidecar `yaml:"blinded_blob_sidecars"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBlockContents) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blindedBlockContentsYAML{
		BlindedBlock:        b.BlindedBlock,
		BlindedBlobSidecars: b.BlindedBlobSidecars,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBlockContents) UnmarshalYAML(input []byte) error {
	var data blindedBlockContentsJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(bytes)
}
