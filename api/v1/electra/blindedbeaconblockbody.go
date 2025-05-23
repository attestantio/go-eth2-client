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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BlindedBeaconBlockBody represents the body of a blinded beacon block.
type BlindedBeaconBlockBody struct {
	RANDAOReveal           phase0.BLSSignature `ssz-size:"96"`
	ETH1Data               *phase0.ETH1Data
	Graffiti               [32]byte                      `ssz-size:"32"`
	ProposerSlashings      []*phase0.ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings      []*electra.AttesterSlashing   `ssz-max:"1"`
	Attestations           []*electra.Attestation        `ssz-max:"8"`
	Deposits               []*phase0.Deposit             `ssz-max:"16"`
	VoluntaryExits         []*phase0.SignedVoluntaryExit `ssz-max:"16"`
	SyncAggregate          *altair.SyncAggregate
	ExecutionPayloadHeader *deneb.ExecutionPayloadHeader
	BLSToExecutionChanges  []*capella.SignedBLSToExecutionChange `ssz-max:"16"`
	BlobKZGCommitments     []deneb.KZGCommitment                 `ssz-max:"4096" ssz-size:"?,48"`
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
