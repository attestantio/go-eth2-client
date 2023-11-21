// Copyright © 2022 Attestant Limited.
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

package bellatrix

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	RANDAOReveal      phase0.BLSSignature `ssz-size:"96"`
	ETH1Data          *phase0.ETH1Data
	Graffiti          [32]byte                      `ssz-size:"32"`
	ProposerSlashings []*phase0.ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings []*phase0.AttesterSlashing    `ssz-max:"2"`
	Attestations      []*phase0.Attestation         `ssz-max:"128"`
	Deposits          []*phase0.Deposit             `ssz-max:"16"`
	VoluntaryExits    []*phase0.SignedVoluntaryExit `ssz-max:"16"`
	SyncAggregate     *altair.SyncAggregate
	ExecutionPayload  *ExecutionPayload
}

// beaconBlockBodyJSON is the spec representation of the struct.
type beaconBlockBodyJSON struct {
	RANDAOReveal      string                        `json:"randao_reveal"`
	ETH1Data          *phase0.ETH1Data              `json:"eth1_data"`
	Graffiti          string                        `json:"graffiti"`
	ProposerSlashings []*phase0.ProposerSlashing    `json:"proposer_slashings"`
	AttesterSlashings []*phase0.AttesterSlashing    `json:"attester_slashings"`
	Attestations      []*phase0.Attestation         `json:"attestations"`
	Deposits          []*phase0.Deposit             `json:"deposits"`
	VoluntaryExits    []*phase0.SignedVoluntaryExit `json:"voluntary_exits"`
	SyncAggregate     *altair.SyncAggregate         `json:"sync_aggregate"`
	ExecutionPayload  *ExecutionPayload             `json:"execution_payload"`
}

// beaconBlockBodyYAML is the spec representation of the struct.
type beaconBlockBodyYAML struct {
	RANDAOReveal      string                        `yaml:"randao_reveal"`
	ETH1Data          *phase0.ETH1Data              `yaml:"eth1_data"`
	Graffiti          string                        `yaml:"graffiti"`
	ProposerSlashings []*phase0.ProposerSlashing    `yaml:"proposer_slashings"`
	AttesterSlashings []*phase0.AttesterSlashing    `yaml:"attester_slashings"`
	Attestations      []*phase0.Attestation         `yaml:"attestations"`
	Deposits          []*phase0.Deposit             `yaml:"deposits"`
	VoluntaryExits    []*phase0.SignedVoluntaryExit `yaml:"voluntary_exits"`
	SyncAggregate     *altair.SyncAggregate         `yaml:"sync_aggregate"`
	ExecutionPayload  *ExecutionPayload             `yaml:"execution_payload"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockBodyJSON{
		RANDAOReveal:      fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:          b.ETH1Data,
		Graffiti:          fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings: b.ProposerSlashings,
		AttesterSlashings: b.AttesterSlashings,
		Attestations:      b.Attestations,
		Deposits:          b.Deposits,
		VoluntaryExits:    b.VoluntaryExits,
		SyncAggregate:     b.SyncAggregate,
		ExecutionPayload:  b.ExecutionPayload,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	var data beaconBlockBodyJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

func (b *BeaconBlockBody) unpack(data *beaconBlockBodyJSON) error {
	if data.RANDAOReveal == "" {
		return errors.New("RANDAO reveal missing")
	}
	randaoReveal, err := hex.DecodeString(strings.TrimPrefix(data.RANDAOReveal, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for RANDAO reveal")
	}
	if len(randaoReveal) != phase0.SignatureLength {
		return errors.New("incorrect length for RANDAO reveal")
	}
	copy(b.RANDAOReveal[:], randaoReveal)
	if data.ETH1Data == nil {
		return errors.New("ETH1 data missing")
	}
	b.ETH1Data = data.ETH1Data
	if data.Graffiti == "" {
		return errors.New("graffiti missing")
	}
	graffiti, err := hex.DecodeString(strings.TrimPrefix(data.Graffiti, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for graffiti")
	}
	if len(graffiti) != phase0.GraffitiLength {
		return errors.New("incorrect length for graffiti")
	}
	copy(b.Graffiti[:], graffiti)
	if data.ProposerSlashings == nil {
		return errors.New("proposer slashings missing")
	}
	b.ProposerSlashings = data.ProposerSlashings
	if data.AttesterSlashings == nil {
		return errors.New("attester slashings missing")
	}
	b.AttesterSlashings = data.AttesterSlashings
	if data.Attestations == nil {
		return errors.New("attestations missing")
	}
	b.Attestations = data.Attestations
	if data.Deposits == nil {
		return errors.New("deposits missing")
	}
	b.Deposits = data.Deposits
	if data.VoluntaryExits == nil {
		return errors.New("voluntary exits missing")
	}
	b.VoluntaryExits = data.VoluntaryExits
	if data.SyncAggregate == nil {
		return errors.New("sync aggregate missing")
	}
	b.SyncAggregate = data.SyncAggregate
	if data.ExecutionPayload == nil {
		return errors.New("execution payload missing")
	}
	b.ExecutionPayload = data.ExecutionPayload

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconBlockBody) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconBlockBodyYAML{
		RANDAOReveal:      fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:          b.ETH1Data,
		Graffiti:          fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings: b.ProposerSlashings,
		AttesterSlashings: b.AttesterSlashings,
		Attestations:      b.Attestations,
		Deposits:          b.Deposits,
		VoluntaryExits:    b.VoluntaryExits,
		SyncAggregate:     b.SyncAggregate,
		ExecutionPayload:  b.ExecutionPayload,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data beaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return b.unpack(&data)
}

// String returns a string version of the structure.
func (b *BeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
