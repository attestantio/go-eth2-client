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

package spec

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

// Tree represents a Merkle tree structure for beacon state data.
type Tree struct {
	node *ssz.Node
}

// Root returns the root hash of the tree.
func (t *Tree) Root() (phase0.Hash32, error) {
	if t.node == nil {
		return phase0.Hash32{}, errors.New("nil tree")
	}
	var root phase0.Hash32
	copy(root[:], t.node.Hash())

	return root, nil
}

// Subtree returns a subtree at the given generalized index.
func (t *Tree) Subtree(index int) (*Tree, error) {
	if t.node == nil {
		return nil, errors.New("nil tree")
	}
	node, err := t.node.Get(index)
	if err != nil {
		return nil, err
	}

	return &Tree{node: node}, nil
}

// Prove generates a Merkle proof for a given generalized index.
func (t *Tree) Prove(index int) ([]phase0.Hash32, error) {
	if t.node == nil {
		return nil, errors.New("nil tree")
	}
	proof, err := t.node.Prove(index)
	if err != nil {
		return nil, err
	}

	proofBytes := make([]phase0.Hash32, len(proof.Hashes))
	for i, hash := range proof.Hashes {
		copy(proofBytes[i][:], hash)
	}

	return proofBytes, nil
}
