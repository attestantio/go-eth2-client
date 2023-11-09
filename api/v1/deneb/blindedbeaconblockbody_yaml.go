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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// blindedBeaconBlockBodyYAML is the spec representation of the struct.
type blindedBeaconBlockBodyYAML struct {
	RANDAOReveal           string                                `yaml:"randao_reveal"`
	ETH1Data               *phase0.ETH1Data                      `yaml:"eth1_data"`
	Graffiti               string                                `yaml:"graffiti"`
	ProposerSlashings      []*phase0.ProposerSlashing            `yaml:"proposer_slashings"`
	AttesterSlashings      []*phase0.AttesterSlashing            `yaml:"attester_slashings"`
	Attestations           []*phase0.Attestation                 `yaml:"attestations"`
	Deposits               []*phase0.Deposit                     `yaml:"deposits"`
	VoluntaryExits         []*phase0.SignedVoluntaryExit         `yaml:"voluntary_exits"`
	SyncAggregate          *altair.SyncAggregate                 `yaml:"sync_aggregate"`
	ExecutionPayloadHeader *deneb.ExecutionPayloadHeader         `yaml:"execution_payload_header"`
	BLSToExecutionChanges  []*capella.SignedBLSToExecutionChange `yaml:"bls_to_execution_changes"`
	BlobKZGCommitments     []string                              `yaml:"blob_kzg_commitments"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BlindedBeaconBlockBody) MarshalYAML() ([]byte, error) {
	blobKZGCommitments := make([]string, len(b.BlobKZGCommitments))
	for i := range b.BlobKZGCommitments {
		blobKZGCommitments[i] = b.BlobKZGCommitments[i].String()
	}

	yamlBytes, err := yaml.MarshalWithOptions(&blindedBeaconBlockBodyYAML{
		RANDAOReveal:           fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:               b.ETH1Data,
		Graffiti:               fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings:      b.ProposerSlashings,
		AttesterSlashings:      b.AttesterSlashings,
		Attestations:           b.Attestations,
		Deposits:               b.Deposits,
		VoluntaryExits:         b.VoluntaryExits,
		SyncAggregate:          b.SyncAggregate,
		ExecutionPayloadHeader: b.ExecutionPayloadHeader,
		BLSToExecutionChanges:  b.BLSToExecutionChanges,
		BlobKZGCommitments:     blobKZGCommitments,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BlindedBeaconBlockBody) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data blindedBeaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return b.unpack(&data)
}
