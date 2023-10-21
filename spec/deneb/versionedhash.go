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

// VersionedHash is a hash with version information.
type VersionedHash [32]byte

// VersionedHashLength is the number of bytes in a versioned hash.
const VersionedHashLength = 32

// String returns a string version of the structure.
func (h VersionedHash) String() string {
	return fmt.Sprintf("%#x", h)
}

// Format formats the root.
func (h VersionedHash) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, h.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, h[:])
	default:
		fmt.Fprintf(state, "%"+format, h[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *VersionedHash) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+VersionedHashLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(h[:], input[3:3+VersionedHashLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+VersionedHashLength*2]))
	}

	if length != VersionedHashLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (h *VersionedHash) MarshalJSON() ([]byte, error) {
	if h == nil {
		return nil, errors.New("value nil")
	}

	return []byte(fmt.Sprintf(`"%#x"`, h)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (h *VersionedHash) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}
	if len(input) != 1+2+VersionedHashLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(h[:], input[3:3+VersionedHashLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+VersionedHashLength*2]))
	}

	if length != VersionedHashLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (h *VersionedHash) MarshalYAML() ([]byte, error) {
	if h == nil {
		return nil, errors.New("value nil")
	}

	return []byte(fmt.Sprintf(`'%#x'`, h)), nil
}
