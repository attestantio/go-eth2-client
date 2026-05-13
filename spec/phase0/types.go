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
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Epoch is an epoch number.
type Epoch uint64

// UnmarshalJSON implements json.Unmarshaler. The spec encodes Epoch as a
// quoted decimal string, but some clients (notably Erigon's Caplin) emit
// uint64 fields in state-tree types as bare JSON numbers. Accept both.
func (e *Epoch) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	s := string(input)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", s)
	}

	*e = Epoch(val)

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e Epoch) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `"%d"`, e), nil
}

// CommitteeIndex is a committee index at a slot.
type CommitteeIndex uint64

// Version is a fork version.
type Version [4]byte

// DomainType is a domain type.
type DomainType [4]byte

// ForkDigest is a digest of fork data.
type ForkDigest [4]byte

// Domain is a signature domain.
type Domain [32]byte
