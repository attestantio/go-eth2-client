// Copyright © 2022, 2023 Attestant Limited.
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
)

// Transaction is an opaque execution layer transaction.
type Transaction []byte

// UnmarshalJSON implements json.Unmarshaler.
func (t Transaction) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}

	_, err := hex.Decode(t, input[3:len(input)-1])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:len(input)-1]))
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (t Transaction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%#x"`, t)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (t *Transaction) UnmarshalYAML(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'\'', '0', 'x'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'\''}) {
		return errors.New("invalid suffix")
	}

	_, err := hex.Decode(*t, input[3:len(input)-1])
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[3:len(input)-1]))
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (t Transaction) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf(`'%#x'`, t)), nil
}
