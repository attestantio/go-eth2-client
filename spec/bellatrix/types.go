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

// Transaction is an opaque execution layer transaction.
type Transaction []byte

// ExecutionAddress is a execution address.
type ExecutionAddress [20]byte

// String returns a string version of the structure.
func (a ExecutionAddress) String() string {
	return fmt.Sprintf("%#x", a)
}
