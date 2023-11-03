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

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

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

type SuffixStateDiffJSON struct {
	Suffix       uint8   `json:"suffix"`
	CurrentValue *string `json:"currentValue"`
	NewValue     *string `json:"newValue"`
}

func (s *SuffixStateDiff) UnmarshalJSON(input []byte) error {
	var (
		ssd SuffixStateDiffJSON
		err error
	)

	if err := json.Unmarshal(input, &ssd); err != nil {
		return fmt.Errorf("error unmarshalling JSON SuffixStateDiff: %w", err)
	}
	s.Suffix = ssd.Suffix
	if ssd.CurrentValue != nil {
		s.CurrentValue, err = hex.DecodeString(strings.TrimPrefix(*ssd.CurrentValue, "0x"))
		if err != nil {
			return fmt.Errorf("error decoding currentValue string: %w", err)
		}
	}
	if ssd.NewValue != nil {
		s.NewValue, err = hex.DecodeString(strings.TrimPrefix(*ssd.NewValue, "0x"))
		if err != nil {
			return fmt.Errorf("error decoding newValue string: %w", err)
		}
	}

	return nil
}

type StemStateDiff struct {
	Stem        [31]byte           `ssz-size:"31"`
	SuffixDiffs []*SuffixStateDiff `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

type stemStateDiffJSON struct {
	Stem        string             `json:"stem"`
	SuffixDiffs []*SuffixStateDiff `json:"suffixDiffs"`
}

func (s *StemStateDiff) UnmarshalJSON(input []byte) error {
	var ssd stemStateDiffJSON
	if err := json.Unmarshal(input, &ssd); err != nil {
		return fmt.Errorf("error unmarshalling JSON StemStateDiff: %w", err)
	}
	stem, err := hex.DecodeString(strings.TrimPrefix(ssd.Stem, "0x"))
	if err != nil {
		return fmt.Errorf("error decoding stem string: %w", err)
	}

	copy(s.Stem[:], stem)
	s.SuffixDiffs = ssd.SuffixDiffs
	return nil
}

func (s *VerkleProof) UnmarshalJSON(input []byte) error {
	// TODO: parse VerkleProof

	return nil
}

type ExecutionWitness struct {
	StateDiff   []*StemStateDiff `ssz-max:"1048576,1073741824" ssz-size:"?,?" json:"stateDiff"`
	VerkleProof *VerkleProof     `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

type executionWitnessJSON struct {
	StateDiff []*StemStateDiff `json"stateDiff""`
}

func (ew *ExecutionWitness) UnmarshalJSON(input []byte) error {
	var res executionWitnessJSON
	err := json.Unmarshal(input, &res)
	if err != nil {
		return err
	}
	ew.StateDiff = make([]*StemStateDiff, len(res.StateDiff))
	for i, sd := range res.StateDiff {
		ew.StateDiff[i] = sd
	}
	return nil
}
