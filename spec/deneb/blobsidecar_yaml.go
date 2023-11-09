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

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// blobSidecarYAML is the spec representation of the struct.
type blobSidecarYAML struct {
	Index                       uint64                          `yaml:"index"`
	Blob                        Blob                            `yaml:"blob"`
	KZGCommitment               KZGCommitment                   `yaml:"kzg_commitment"`
	KZGProof                    KZGProof                        `yaml:"kzg_proof"`
	SignedBlockHeader           *phase0.SignedBeaconBlockHeader `yaml:"signed_block_header"`
	KZGCommitmentInclusionProof KZGCommitmentInclusionProof     `yaml:"kzg_commitment_inclusion_proof"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BlobSidecar) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blobSidecarYAML{
		Index:                       uint64(b.Index),
		Blob:                        b.Blob,
		KZGCommitment:               b.KZGCommitment,
		KZGProof:                    b.KZGProof,
		SignedBlockHeader:           b.SignedBlockHeader,
		KZGCommitmentInclusionProof: b.KZGCommitmentInclusionProof,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BlobSidecar) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var data blobSidecarJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	bytes, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(bytes)
}
