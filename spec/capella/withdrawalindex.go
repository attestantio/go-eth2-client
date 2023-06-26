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

package capella

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// WithdrawalIndex is the index of a withdrawal operation.
type WithdrawalIndex uint64

// UnmarshalJSON implements json.Unmarshaler.
func (w *WithdrawalIndex) UnmarshalJSON(input []byte) error {
	if len(input) == 0 {
		return errors.New("input missing")
	}

	if !bytes.HasPrefix(input, []byte{'"'}) {
		return errors.New("invalid prefix")
	}
	if !bytes.HasSuffix(input, []byte{'"'}) {
		return errors.New("invalid suffix")
	}

	val, err := strconv.ParseUint(string(input[1:len(input)-1]), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid value %s", string(input[1:len(input)-1]))
	}
	*w = WithdrawalIndex(val)

	return nil
}

// MarshalJSON implements json.Marshaler.
func (w WithdrawalIndex) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, w)), nil
}
