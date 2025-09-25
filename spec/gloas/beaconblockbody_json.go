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

package gloas

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// beaconBlockBodyJSON is the spec representation of the struct.
type beaconBlockBodyJSON struct {
	RANDAOReveal              string                                `json:"randao_reveal"`
	ETH1Data                  *phase0.ETH1Data                      `json:"eth1_data"`
	Graffiti                  string                                `json:"graffiti"`
	ProposerSlashings         []*phase0.ProposerSlashing            `json:"proposer_slashings"`
	AttesterSlashings         []*electra.AttesterSlashing           `json:"attester_slashings"`
	Attestations              []*electra.Attestation                `json:"attestations"`
	Deposits                  []*phase0.Deposit                     `json:"deposits"`
	VoluntaryExits            []*phase0.SignedVoluntaryExit         `json:"voluntary_exits"`
	SyncAggregate             *altair.SyncAggregate                 `json:"sync_aggregate"`
	BLSToExecutionChanges     []*capella.SignedBLSToExecutionChange `json:"bls_to_execution_changes"`
	SignedExecutionPayloadBid *SignedExecutionPayloadBid            `json:"signed_execution_payload_bid"`
	PayloadAttestations       []*PayloadAttestation                 `json:"payload_attestations"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockBodyJSON{
		RANDAOReveal:              fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:                  b.ETH1Data,
		Graffiti:                  fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings:         b.ProposerSlashings,
		AttesterSlashings:         b.AttesterSlashings,
		Attestations:              b.Attestations,
		Deposits:                  b.Deposits,
		VoluntaryExits:            b.VoluntaryExits,
		SyncAggregate:             b.SyncAggregate,
		BLSToExecutionChanges:     b.BLSToExecutionChanges,
		SignedExecutionPayloadBid: b.SignedExecutionPayloadBid,
		PayloadAttestations:       b.PayloadAttestations,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	var data beaconBlockBodyJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// RANDAO reveal
	if data.RANDAOReveal == "" {
		return errors.New("randao reveal missing")
	}
	randaoReveal, err := hex.DecodeString(strings.TrimPrefix(data.RANDAOReveal, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid randao reveal")
	}
	copy(b.RANDAOReveal[:], randaoReveal)

	b.ETH1Data = data.ETH1Data

	// Graffiti
	if data.Graffiti == "" {
		return errors.New("graffiti missing")
	}
	graffiti, err := hex.DecodeString(strings.TrimPrefix(data.Graffiti, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid graffiti")
	}
	copy(b.Graffiti[:], graffiti)

	b.ProposerSlashings = data.ProposerSlashings
	b.AttesterSlashings = data.AttesterSlashings
	b.Attestations = data.Attestations
	b.Deposits = data.Deposits
	b.VoluntaryExits = data.VoluntaryExits
	b.SyncAggregate = data.SyncAggregate
	b.BLSToExecutionChanges = data.BLSToExecutionChanges
	b.SignedExecutionPayloadBid = data.SignedExecutionPayloadBid
	b.PayloadAttestations = data.PayloadAttestations

	return nil
}
