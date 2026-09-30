// Copyright © 2026 Attestant Limited.
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

package gloas_test

import (
	"encoding/hex"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestBuilderConfigHashTreeRoots pins the merkleization of the builder
// configuration types.
//
// UNVERIFIED AGAINST THE CONTRACT. All four carry ssz-index tags, which makes
// them progressive containers: the generated code ends in
// MerkleizeProgressiveWithActiveFields rather than Merkleize, unlike the plain
// Builder container in this package, which has no ssz-index tags. That choice
// changes hash tree roots but not serialisation, so nothing else in the suite
// would notice if it were wrong.
//
// It matters most for BuilderRequestAuth, whose HTR *is* the signing root:
// SignedBuilderRequestAuth.Signature is the proposer's signature over it. If
// the pinned contract defines it as an ordinary SSZ container, every
// authorization this library helps sign has the wrong root, the builder
// declines the bid request, and -- since the beacon node does not verify the
// authorization -- the proposal quietly falls back to a p2p or local build
// with no error raised anywhere.
//
// These values are what the generated code produces today, recorded so that
// changing the tags is a deliberate act with a visible diff rather than a
// silent change to a signing root. Someone with the contract to hand should
// confirm them and replace this comment with the reference.
func TestBuilderConfigHashTreeRoots(t *testing.T) {
	auth := &gloas.BuilderRequestAuth{Data: []byte{0x01, 0x02, 0x03}, Slot: 60300}
	signed := &gloas.SignedBuilderRequestAuth{Message: auth, Signature: phase0.BLSSignature{0x0a, 0x0b}}
	entry := &gloas.BuilderEntry{
		URL:                 []byte("https://builder.example"),
		Auth:                signed,
		BuilderPubkeys:      []phase0.BLSPubKey{{0x11}},
		MaxExecutionPayment: 5,
		MinBid:              7,
		BuilderBoostFactor:  9,
	}
	config := &gloas.BuilderConfig{MinBid: 3, BuilderBoostFactor: 11, Builders: []*gloas.BuilderEntry{entry}}

	tests := []struct {
		name string
		root func() ([32]byte, error)
		want string
	}{
		{
			name: "BuilderRequestAuth",
			root: auth.HashTreeRoot,
			want: "5c080ef2a5aa78fe3a3875bf3d57c2c570be8fde89d3492cc4a880ac3593cf7e",
		},
		{
			name: "SignedBuilderRequestAuth",
			root: signed.HashTreeRoot,
			want: "6e125f65556a0aae894121a310cd2bc02297c0126e07eb8e2490ea58123bb1dd",
		},
		{
			name: "BuilderEntry",
			root: entry.HashTreeRoot,
			want: "2f5754ea0ec026795b0df916181f73c36772b02bbd806c4f0811e4eb13bf7cf1",
		},
		{
			name: "BuilderConfig",
			root: config.HashTreeRoot,
			want: "f3683e36e2091aa8bbe573098aae17e9cec932995d8a3c6bdd2f02d9130fa9be",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := test.root()
			require.NoError(t, err)

			want, err := hex.DecodeString(test.want)
			require.NoError(t, err)
			require.Equal(t, want, root[:])
		})
	}
}
