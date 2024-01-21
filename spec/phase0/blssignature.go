// Copyright © 2020 - 2023 Attestant Limited.
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

// BLSSignature is a BLS12-381 signature.
type BLSSignature [96]byte

// SignatureLength is the number of bytes in a signature.
const SignatureLength = 96

var (
	emptyBLSSignature    = BLSSignature{}
	infinityBLSSignature = BLSSignature{
		0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
)

// IsZero returns true if the signature is zero.
func (s BLSSignature) IsZero() bool {
	return bytes.Equal(s[:], emptyBLSSignature[:])
}

// IsInfinity returns true if the signature is infinity.
func (s BLSSignature) IsInfinity() bool {
	return bytes.Equal(s[:], infinityBLSSignature[:])
}

// String returns a string version of the structure.
func (s BLSSignature) String() string {
	return fmt.Sprintf("%#x", s)
}

// Format formats the signature.
func (s BLSSignature) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, s.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, s[:])
	default:
		fmt.Fprintf(state, "%"+format, s[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *BLSSignature) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+SignatureLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(s[:], input[3:3+SignatureLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+SignatureLength*2]))
	}

	if length != SignatureLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (s BLSSignature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, s)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *BLSSignature) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+SignatureLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(s[:], input[3:3+SignatureLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+SignatureLength*2]))
	}

	if length != SignatureLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s BLSSignature) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, s)), nil
}
