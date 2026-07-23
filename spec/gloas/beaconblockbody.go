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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BeaconBlockBody represents the body of a beacon block for EIP-7732.
type BeaconBlockBody struct {
	RANDAOReveal              phase0.BLSSignature                   `ssz-index:"0"`
	ETH1Data                  *phase0.ETH1Data                      `ssz-index:"1"`
	Graffiti                  [32]byte                              `ssz-index:"2"`
	ProposerSlashings         []*phase0.ProposerSlashing            `ssz-index:"3"  ssz-type:"progressive-list"`
	AttesterSlashings         []*AttesterSlashing                   `ssz-index:"4"  ssz-type:"progressive-list"`
	Attestations              []*Attestation                        `ssz-index:"5"  ssz-type:"progressive-list"`
	Deposits                  []*phase0.Deposit                     `ssz-index:"6"  ssz-type:"progressive-list"`
	VoluntaryExits            []*phase0.SignedVoluntaryExit         `ssz-index:"7"  ssz-type:"progressive-list"`
	SyncAggregate             *altair.SyncAggregate                 `ssz-index:"8"`
	BLSToExecutionChanges     []*capella.SignedBLSToExecutionChange `ssz-index:"9"  ssz-type:"progressive-list"`
	SignedExecutionPayloadBid *SignedExecutionPayloadBid            `ssz-index:"10"`
	PayloadAttestations       []*PayloadAttestation                 `ssz-index:"11" ssz-type:"progressive-list"`
	ParentExecutionRequests   *ExecutionRequests                    `ssz-index:"12"`
}

// String returns a string version of the structure.
func (b *BeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
