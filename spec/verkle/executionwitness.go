// Copyright © 2023 Guillaume Ballet.
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

package verkle

// import "github.com/attestantio/go-eth2-client/spec/phase0"

const IPA_PROOF_DEPTH = 8

type IPAProof struct {
	CL              [IPA_PROOF_DEPTH][32]byte `ssz-size:"8,32"`
	CR              [IPA_PROOF_DEPTH][32]byte `ssz-size:"8,32"`
	FinalEvaluation [32]byte                  `ssz-size:"32"`
}

type VerkleProof struct {
	OtherStems            [][]byte  `ssz-max:"65536,31"`
	DepthExtensionPresent []byte    `ssz-max:"65536"`
	CommitmentsByPath     [][]byte  `ssz-max:"65536,32"`
	D                     [32]byte  `ssz-size:"32"`
	IPAProof              *IPAProof `ssz-size:"544"`
}

type SuffixStateDiff struct {
	Suffix       uint8  `ssz-size:"1"`
	CurrentValue []byte `ssz-max:"32"`
	NewValue     []byte `ssz-max:"32"`
}

type StemStateDiff struct {
	Stem        [31]byte          `ssz-size:"31"`
	SuffixDiffs []SuffixStateDiff `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}
