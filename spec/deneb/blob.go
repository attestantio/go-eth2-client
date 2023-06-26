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

// Blob is a data blob.
type Blob [131072]byte

// BlobLength is the number of bytes in a data blob.
const BlobLength = 131072

// String returns a string version of the structure.
func (b Blob) String() string {
	return fmt.Sprintf("%#x", b)
}

// Format formats the blob.
func (b Blob) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, b.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, b[:])
	default:
		fmt.Fprintf(state, "%"+format, b[:])
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Blob) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("Invalid suffix")
	}
	if len(input) != 1+2+BlobLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(b[:], input[3:3+BlobLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+BlobLength*2]))
	}

	if length != BlobLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Blob) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, b)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *Blob) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("Invalid suffix")
	}
	if len(input) != 1+2+BlobLength*2+1 {
		return errors.New("incorrect length")
	}

	length, err := hex.Decode(b[:], input[3:3+BlobLength*2])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:3+BlobLength*2]))
	}

	if length != BlobLength {
		return errors.New("incorrect length")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b Blob) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, b)), nil
}
