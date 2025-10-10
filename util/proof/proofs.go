// Copyright © 2025 Attestant Limited.
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

/*
Package proof provides helper functions for merkle tree operations.
# Merkle Trees and Proof Generation

We demonstrate the Merkle tree functionality in the context of the beacon state.
There are two key concepts for understanding how fields are accessed in the tree:

**Field Index vs Generalized Index:**

 1. Field Index: Position of a field in the BeaconState struct (0-based)

    BeaconState struct {
    GenesisTime           [0]
    GenesisValidatorsRoot [1]  <-- Field index is 1
    Slot                  [2]
    Fork                  [3]
    ...
    }

 2. Generalized Index: Position in the binary Merkle tree

    Root [1]
    /              \
    [2]                [3]
    /    \             /    \
    [4]     [5]        [6]     [7]
    /  \    /  \      /  \     /  \
    [8] [9] [10][11] [12][13] [14][15]  <-- Leaf level
    |   |   |   |    |   |    |   |
    GT GVR  S   F   ... ...  ... ...    <-- Fields

For example, GenesisValidatorsRoot (field index 1):

	depth = ceil(log2(num_fields))  // = 4 for this example
	general_index = 2^depth + field_index = 9

**Proof Generation:**

A Merkle proof consists of the sibling hashes needed to reconstruct the root:

	                Root [R]
	            /              \
	        [A]                [B]
	      /    \             /    \
	   [C]     [D]        [E]     [F]
	  /  \    /  \      /  \     /  \
	[G] [H] [J] [K]   [L] [M]  [N] [O]
	    ^
	    GenesisValidatorsRoot

To prove GenesisValidatorsRoot exists:
1. Start at H (GenesisValidatorsRoot)
2. Collect siblings: [G, D, B]
3. To verify:
  - Hash(G,H) = C
  - Hash(C,D) = A
  - Hash(A,B) = R

Since the proof process can be applied to any leaf (A-O), we need to always
use the generalized index for proof generation.
*/
package proof

import (
	"errors"
	"math"
	"reflect"
)

// NumFields returns the number of fields in a struct.
func NumFields(o any) int {
	if o == nil {
		return 0
	}

	t := reflect.TypeOf(o)

	// If it's a pointer, get the underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// If it's not a struct, return 0
	if t.Kind() != reflect.Struct {
		return 0
	}

	return t.NumField()
}

// FieldIndex returns the index of a field in a struct.
// The index represents the field's position in the struct's memory layout.
func FieldIndex(o any, fieldName string) (int, error) {
	if o == nil {
		return 0, errors.New("nil object")
	}

	t := reflect.TypeOf(o)

	// If it's a pointer, get the underlying type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// If it's not a struct, return error
	if t.Kind() != reflect.Struct {
		return 0, errors.New("not a struct")
	}

	field, ok := t.FieldByName(fieldName)
	if !ok {
		return 0, errors.New("field not found")
	}

	return field.Index[0], nil
}

// FieldGeneralizedIndex obtains the generalized index of a field leaf using the field name.
func FieldGeneralizedIndex(o any, fieldName string) (int, error) {
	index, err := FieldIndex(o, fieldName)
	if err != nil {
		return 0, err
	}

	return LeafGeneralizedIndex(index, NumFields(o)), nil
}

// LeafGeneralizedIndex calculates the generalized index of a leaf in a binary Merkle tree.
// The generalized index is the absolute position of the leaf in the tree's array representation.
func LeafGeneralizedIndex(leafIdx int, nleaves int) int {
	depth := TreeDepth(uint64(nleaves))
	generalizedIndex := math.Pow(2, float64(depth)) + float64(leafIdx)

	return int(generalizedIndex)
}

// TreeDepth returns the depth of a binary tree with given number of leaves.
func TreeDepth(nleaves uint64) uint8 {
	depth := math.Ceil(math.Log2(float64(nleaves)))

	return uint8(depth)
}
