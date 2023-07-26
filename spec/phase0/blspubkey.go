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
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// BLSPubKey is a BLS12-381 public key.
type BLSPubKey [48]byte

// String returns a string version of the structure.
func (pk BLSPubKey) String() string {
	return fmt.Sprintf("%#x", pk)
}

// Format formats the public key.
func (pk BLSPubKey) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, pk.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, pk[:])
	default:
		fmt.Fprintf(state, "%"+format, pk[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (pk *BLSPubKey) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+PublicKeyLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(pk[:], input[3:3+PublicKeyLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+PublicKeyLength*2]))
	}

	if length != PublicKeyLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (pk BLSPubKey) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, pk)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (pk *BLSPubKey) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+PublicKeyLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(pk[:], input[3:3+PublicKeyLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+PublicKeyLength*2]))
	}

	if length != PublicKeyLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (pk BLSPubKey) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, pk)), nil
}
