// Copyright © 2022 - 2024 Attestant Limited.
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

package bellatrix

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// ExecutionAddress is a execution address.
type ExecutionAddress [20]byte

var emptyExecutionAddress = ExecutionAddress{}

// IsZero returns true if the execution address is zero.
func (a ExecutionAddress) IsZero() bool {
	return bytes.Equal(a[:], emptyExecutionAddress[:])
}

// String returns an EIP-55 string version of the address.
func (a ExecutionAddress) String() string {
	bytes := []byte(hex.EncodeToString(a[:]))

	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(bytes)
	hash := keccak.Sum(nil)

	for i := 0; i < len(bytes); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte >>= 4
		} else {
			hashByte &= 0xf
		}
		if bytes[i] > '9' && hashByte > 7 {
			bytes[i] -= 32
		}
	}

	return fmt.Sprintf("0x%s", string(bytes))
}

// Format formats the execution address.
func (a ExecutionAddress) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, a.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, a[:])
	default:
		fmt.Fprintf(state, "%"+format, a[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *ExecutionAddress) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+ExecutionAddressLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(a[:], input[3:3+ExecutionAddressLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+ExecutionAddressLength*2]))
	}

	if length != ExecutionAddressLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (a ExecutionAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, a.String())), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *ExecutionAddress) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+ExecutionAddressLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(a[:], input[3:3+ExecutionAddressLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+ExecutionAddressLength*2]))
	}

	if length != ExecutionAddressLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a ExecutionAddress) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%s'`, a.String())), nil
}
