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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	RANDAOReveal      []byte `ssz-size:"96"`
	ETH1Data          *ETH1Data
	Graffiti          []byte                 `ssz-size:"32"`
	ProposerSlashings []*ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings []*AttesterSlashing    `ssz-max:"2"`
	Attestations      []*Attestation         `ssz-max:"128"`
	Deposits          []*Deposit             `ssz-max:"16"`
	VoluntaryExits    []*SignedVoluntaryExit `ssz-max:"16"`
}

// beaconBlockBodyJSON is the spec representation of the struct.
type beaconBlockBodyJSON struct {
	RANDAOReveal      string                 `json:"randao_reveal"`
	ETH1Data          *ETH1Data              `json:"eth1_data" yaml:"eth1_data"`
	Graffiti          string                 `json:"graffiti"`
	ProposerSlashings []*ProposerSlashing    `json:"proposer_slashings"`
	AttesterSlashings []*AttesterSlashing    `json:"attester_slashings"`
	Attestations      []*Attestation         `json:"attestations"`
	Deposits          []*Deposit             `json:"deposits"`
	VoluntaryExits    []*SignedVoluntaryExit `json:"voluntary_exits"`
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
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	var err error

	var beaconBlockBodyJSON beaconBlockBodyJSON
	if err = json.Unmarshal(input, &beaconBlockBodyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconBlockBodyJSON.RANDAOReveal == "" {
		return errors.New("RANDAO reveal missing")
	}
	if b.RANDAOReveal, err = hex.DecodeString(strings.TrimPrefix(beaconBlockBodyJSON.RANDAOReveal, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for RANDAO reveal")
	}
	if len(b.RANDAOReveal) != signatureLength {
		return errors.New("incorrect length for RANDAO reveal")
	}
	if beaconBlockBodyJSON.ETH1Data == nil {
		return errors.New("ETH1 data missing")
	}
	b.ETH1Data = beaconBlockBodyJSON.ETH1Data
	if beaconBlockBodyJSON.Graffiti == "" {
		return errors.New("graffiti missing")
	}
	if b.Graffiti, err = hex.DecodeString(strings.TrimPrefix(beaconBlockBodyJSON.Graffiti, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for graffiti")
	}
	if len(b.Graffiti) != graffitiLength {
		return errors.New("incorrect length for graffiti")
	}
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

	return nil
}

// String returns a string version of the structure.
func (b *BeaconBlockBody) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
