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

// KZGCommitment is an KZG commitment.
type KZGCommitment [48]byte

// KZGCommitmentLength is the number of bytes in a KZG commitment.
const KZGCommitmentLength = 48

// String returns a string version of the structure.
func (k KZGCommitment) String() string {
	return fmt.Sprintf("%#x", k)
}

// Format formats the KZG commitment.
func (k KZGCommitment) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, k.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, k[:])
	default:
		fmt.Fprintf(state, "%"+format, k[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (k *KZGCommitment) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+KZGCommitmentLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(k[:], input[3:3+KZGCommitmentLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+KZGCommitmentLength*2]))
	}

	if length != KZGCommitmentLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (k KZGCommitment) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, k)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (k *KZGCommitment) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+KZGCommitmentLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(k[:], input[3:3+KZGCommitmentLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+KZGCommitmentLength*2]))
	}

	if length != KZGCommitmentLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (k KZGCommitment) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, k)), nil
}
