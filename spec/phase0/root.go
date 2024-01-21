// Copyright © 2020, 2024 Attestant Limited.
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

// Root is a merkle root.
type Root [32]byte

var zeroRoot = Root{}

// IsZero returns true if the root is zero.
func (r Root) IsZero() bool {
	return bytes.Equal(r[:], zeroRoot[:])
}

// String returns a string version of the structure.
func (r Root) String() string {
	return fmt.Sprintf("%#x", r)
}

// Format formats the root.
func (r Root) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, r.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, r[:])
	default:
		fmt.Fprintf(state, "%"+format, r[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Root) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+RootLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(r[:], input[3:3+RootLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+RootLength*2]))
	}

	if length != RootLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (r Root) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, r)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (r *Root) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+RootLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(r[:], input[3:3+RootLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+RootLength*2]))
	}

	if length != RootLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (r Root) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, r)), nil
}
