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

package phase0

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Slot is a slot number.
type Slot uint64

// UnmarshalJSON implements json.Unmarshaler. The spec encodes Slot as a
// quoted decimal string, but some clients (notably Erigon's Caplin) emit
// uint64 fields in state-tree types as bare JSON numbers. Accept both.
func (s *Slot) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	str := string(input)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", str)
	}

	*s = Slot(val)

	return nil
}

// MarshalJSON implements json.Marshaler.
func (s Slot) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%d"`, s), nil
}
