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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blindedBeaconBlockBodyJSON is the spec representation of the struct.
type blindedBeaconBlockBodyJSON struct {
	RANDAOReveal           string                            `json:"randao_reveal"`
	ETH1Data               *phase0.ETH1Data                  `json:"eth1_data"`
	Graffiti               string                            `json:"graffiti"`
	ProposerSlashings      []*phase0.ProposerSlashing        `json:"proposer_slashings"`
	AttesterSlashings      []*phase0.AttesterSlashing        `json:"attester_slashings"`
	Attestations           []*phase0.Attestation             `json:"attestations"`
	Deposits               []*phase0.Deposit                 `json:"deposits"`
	VoluntaryExits         []*phase0.SignedVoluntaryExit     `json:"voluntary_exits"`
	SyncAggregate          *altair.SyncAggregate             `json:"sync_aggregate"`
	ExecutionPayloadHeader *bellatrix.ExecutionPayloadHeader `json:"execution_payload_header"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBeaconBlockBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blindedBeaconBlockBodyJSON{
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
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBeaconBlockBody) UnmarshalJSON(input []byte) error {
	var data blindedBeaconBlockBodyJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

func (b *BlindedBeaconBlockBody) unpack(data *blindedBeaconBlockBodyJSON) error {
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
	if data.ExecutionPayloadHeader == nil {
		return errors.New("execution payload header missing")
	}
	b.ExecutionPayloadHeader = data.ExecutionPayloadHeader

	return nil
}
