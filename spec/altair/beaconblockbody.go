// Copyright © 2021 Attestant Limited.
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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

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
	SyncAggregate     *SyncAggregate
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
	SyncAggregate     *SyncAggregate                `json:"sync_aggregate"`
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
	SyncAggregate     *SyncAggregate                `yaml:"sync_aggregate"`
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
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	var beaconBlockBodyJSON beaconBlockBodyJSON
	if err := json.Unmarshal(input, &beaconBlockBodyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&beaconBlockBodyJSON)
}

func (b *BeaconBlockBody) unpack(beaconBlockBodyJSON *beaconBlockBodyJSON) error {
	if beaconBlockBodyJSON.RANDAOReveal == "" {
		return errors.New("RANDAO reveal missing")
	}
	randaoReveal, err := hex.DecodeString(strings.TrimPrefix(beaconBlockBodyJSON.RANDAOReveal, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for RANDAO reveal")
	}
	if len(randaoReveal) != phase0.SignatureLength {
		return errors.New("incorrect length for RANDAO reveal")
	}
	copy(b.RANDAOReveal[:], randaoReveal)
	if beaconBlockBodyJSON.ETH1Data == nil {
		return errors.New("ETH1 data missing")
	}
	b.ETH1Data = beaconBlockBodyJSON.ETH1Data
	if beaconBlockBodyJSON.Graffiti == "" {
		return errors.New("graffiti missing")
	}
	graffiti, err := hex.DecodeString(strings.TrimPrefix(beaconBlockBodyJSON.Graffiti, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for graffiti")
	}
	if len(graffiti) != phase0.GraffitiLength {
		return errors.New("incorrect length for graffiti")
	}
	copy(b.Graffiti[:], graffiti)
	if beaconBlockBodyJSON.ProposerSlashings == nil {
		return errors.New("proposer slashings missing")
	}
	b.ProposerSlashings = beaconBlockBodyJSON.ProposerSlashings
	if beaconBlockBodyJSON.AttesterSlashings == nil {
		return errors.New("attester slashings missing")
	}
	b.AttesterSlashings = beaconBlockBodyJSON.AttesterSlashings
	if beaconBlockBodyJSON.Attestations == nil {
		return errors.New("attestations missing")
	}
	b.Attestations = beaconBlockBodyJSON.Attestations
	if beaconBlockBodyJSON.Deposits == nil {
		return errors.New("deposits missing")
	}
	b.Deposits = beaconBlockBodyJSON.Deposits
	if beaconBlockBodyJSON.VoluntaryExits == nil {
		return errors.New("voluntary exits missing")
	}
	b.VoluntaryExits = beaconBlockBodyJSON.VoluntaryExits
	if beaconBlockBodyJSON.SyncAggregate == nil {
		return errors.New("sync aggregate missing")
	}
	b.SyncAggregate = beaconBlockBodyJSON.SyncAggregate

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
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var beaconBlockBodyJSON beaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &beaconBlockBodyJSON); err != nil {
		return err
	}

	return b.unpack(&beaconBlockBodyJSON)
}

// String returns a string version of the structure.
func (b *BeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
