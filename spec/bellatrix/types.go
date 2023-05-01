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
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Transaction is an opaque execution layer transaction.
type Transaction []byte

// ExecutionAddress is a execution address.
type ExecutionAddress [20]byte

// String returns an EIP-55 string version of the address.
func (a ExecutionAddress) String() string {
	bytes := []byte(fmt.Sprintf("%x", a[:]))

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
