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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// blindedBlobSidecarYAML is the spec representation of the struct.
type blindedBlobSidecarYAML struct {
	BlockRoot       phase0.Root         `yaml:"block_root"`
	Index           uint64              `yaml:"index"`
	Slot            uint64              `yaml:"slot"`
	BlockParentRoot phase0.Root         `yaml:"block_parent_root"`
	ProposerIndex   uint64              `yaml:"proposer_index"`
	BlobRoot        phase0.Root         `yaml:"blob_root"`
	KZGCommitment   deneb.KZGCommitment `yaml:"kzg_commitment"`
	KZGProof        deneb.KZGProof      `yaml:"kzg_proof"`
}

// MarshalYAML implements json.Marshaler.
func (b *BlindedBlobSidecar) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blindedBlobSidecarYAML{
		BlockRoot:       b.BlockRoot,
		Index:           uint64(b.Index),
		Slot:            uint64(b.Slot),
		BlockParentRoot: b.BlockParentRoot,
		ProposerIndex:   uint64(b.ProposerIndex),
		BlobRoot:        b.BlobRoot,
		KZGCommitment:   b.KZGCommitment,
		KZGProof:        b.KZGProof,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements json.Unmarshaler.
func (b *BlindedBlobSidecar) UnmarshalYAML(input []byte) error {
	var data blindedBlobSidecarJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(bytes)
}
