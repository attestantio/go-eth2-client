// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BlindedBeaconBlockBody represents the body of a blinded beacon block.
type BlindedBeaconBlockBody struct {
	RANDAOReveal           phase0.BLSSignature `ssz-size:"96"`
	ETH1Data               *phase0.ETH1Data
	Graffiti               [32]byte                      `ssz-size:"32"`
	ProposerSlashings      []*phase0.ProposerSlashing    `dynssz-max:"MAX_PROPOSER_SLASHINGS" ssz-max:"16"`
	AttesterSlashings      []*electra.AttesterSlashing   `dynssz-max:"MAX_ATTESTER_SLASHINGS" ssz-max:"2"`
	Attestations           []*electra.Attestation        `dynssz-max:"MAX_ATTESTATIONS"       ssz-max:"128"`
	Deposits               []*phase0.Deposit             `dynssz-max:"MAX_DEPOSITS"           ssz-max:"16"`
	VoluntaryExits         []*phase0.SignedVoluntaryExit `dynssz-max:"MAX_VOLUNTARY_EXITS"    ssz-max:"16"`
	SyncAggregate          *altair.SyncAggregate
	ExecutionPayloadHeader *deneb.ExecutionPayloadHeader
	BLSToExecutionChanges  []*capella.SignedBLSToExecutionChange `dynssz-max:"MAX_BLS_TO_EXECUTION_CHANGES"   ssz-max:"16"`
	BlobKZGCommitments     []deneb.KZGCommitment                 `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK" ssz-max:"4096" ssz-size:"?,48"` //nolint:revive
	ExecutionRequests      *electra.ExecutionRequests
}

// String returns a string version of the structure.
func (b *BlindedBeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
