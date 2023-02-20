// Copyright © 2022 Attestant Limited.
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

import "fmt"

// Transactions provides information about transactions
type Transactions struct {
	Transactions []Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

// Transaction is an opaque execution layer transaction.
type Transaction []byte

// ExecutionAddress is a execution address.
type ExecutionAddress [20]byte

// String returns a string version of the structure.
func (a ExecutionAddress) String() string {
	return fmt.Sprintf("%#x", a)
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
