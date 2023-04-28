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

import "fmt"

// VersionedHash is a hash with version information.
type VersionedHash [32]byte

// BlobIndex is the index of a blob in a block.
type BlobIndex uint64

// KzgCommitment is an KZG commitment.
type KzgCommitment [48]byte

// String returns a string version of the structure.
func (k KzgCommitment) String() string {
	return fmt.Sprintf("%#x", k)
}

// Format formats the KZG commitment.
func (k KzgCommitment) Format(state fmt.State, v rune) {
	format := string(v)
	switch v {
	case 's':
		fmt.Fprint(state, k.String())
	case 'x', 'X':
		if state.Flag('#') {
			format = "#" + format
		}
		fmt.Fprintf(state, "%"+format, k[:])
	default:
		fmt.Fprintf(state, "%"+format, k[:])
	}
}
