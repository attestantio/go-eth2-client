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

package deneb

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// kzgCommitmentProofElementLength is the length of each element in the proof.
const kzgCommitmentProofElementLength = 32

// kzgCommitmentProofElements is the number of element in the proof.
const kzgCommitmentProofElements = 17

// KZGCommitmentInclusionProof is the proof of inclusion for a KZG commitment.
type KZGCommitmentInclusionProof [kzgCommitmentProofElements][kzgCommitmentProofElementLength]byte

// UnmarshalJSON implements json.Unmarshaler.
func (k *KZGCommitmentInclusionProof) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if input[0] != '[' {
		return errors.New("invalid prefix")
	}

	values := bytes.Split(input[1:len(input)-1], []byte(","))
	if len(values) != kzgCommitmentProofElements {
		return errors.New("incorrect number of elements")
	}

	for i := range values {
		if err := k.unmarshalElementJSON(i, bytes.TrimSpace(values[i])); err != nil {
			return err
		}
	}

	return nil
}

func (k *KZGCommitmentInclusionProof) unmarshalElementJSON(element int, input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid element prefix")
	}
	if len(input) != 1+2+kzgCommitmentProofElementLength*2+1 {
		return errors.New("incorrect element length")
	}

	_, err := hex.Decode(k[element][:], input[3:3+kzgCommitmentProofElementLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+kzgCommitmentProofElementLength*2]))
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (k *KZGCommitmentInclusionProof) MarshalJSON() ([]byte, error) {
	if k == nil {
		return nil, errors.New("value nil")
	}

	return []byte(fmt.Sprintf(`["%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x"]`,
		k[0],
		k[1],
		k[2],
		k[3],
		k[4],
		k[5],
		k[6],
		k[7],
		k[8],
		k[9],
		k[10],
		k[11],
		k[12],
		k[13],
		k[14],
		k[15],
		k[16],
	)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (k *KZGCommitmentInclusionProof) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if input[0] != '[' {
		return errors.New("invalid prefix")
	}

	values := bytes.Split(input[1:len(input)-1], []byte(","))
	if len(values) != kzgCommitmentProofElements {
		return errors.New("incorrect number of elements")
	}

	for i := range values {
		if err := k.unmarshalElementYAML(i, bytes.TrimSpace(values[i])); err != nil {
			return err
		}
	}

	return nil
}

func (k *KZGCommitmentInclusionProof) unmarshalElementYAML(element int, input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid element prefix")
	}
	if len(input) != 1+2+kzgCommitmentProofElementLength*2+1 {
		return errors.New("incorrect element length")
	}

	_, err := hex.Decode(k[element][:], input[3:3+kzgCommitmentProofElementLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+kzgCommitmentProofElementLength*2]))
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (k KZGCommitmentInclusionProof) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x","%#x"]`,
		k[0],
		k[1],
		k[2],
		k[3],
		k[4],
		k[5],
		k[6],
		k[7],
		k[8],
		k[9],
		k[10],
		k[11],
		k[12],
		k[13],
		k[14],
		k[15],
		k[16],
	)), nil
}
