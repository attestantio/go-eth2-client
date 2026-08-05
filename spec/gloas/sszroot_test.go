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

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

// The expected roots in this file were independently calculated with the
// Ethereum consensus-specs v1.7.0-alpha.12 PySpec at commit
// 84454a9d57f4f49f74c78e9f375347baf30474f2.  That implementation uses
// eth-remerkleable rather than this repository's generated dynamic-ssz codec.
// This independent coverage complements, but does not replace, official
// v1.7.0-alpha.12 consensus-spec vector validation.
//
// Each fixture below mirrors the named PySpec container's non-zero fields.  The
// literals deliberately make a generated-code regression fail even if the same
// regression affects both this package's MarshalSSZ and HashTreeRoot methods.
type sszRootable interface {
	MarshalSSZ() ([]byte, error)
	UnmarshalSSZ([]byte) error
	HashTreeRoot() ([32]byte, error)
}

func TestInheritedContainerSSZRoots(t *testing.T) {
	tests := []struct {
		name     string
		value    sszRootable
		newValue func() sszRootable
		root     string
	}{
		{
			name:     "Attestation",
			value:    testAttestation(),
			newValue: func() sszRootable { return &gloas.Attestation{} },
			root:     "7f14183aef36f97fd2e90104f3afb1e8eeef426d9c2613fae882131db3bb820b",
		},
		{
			name:     "IndexedAttestation",
			value:    testIndexedAttestation(1),
			newValue: func() sszRootable { return &gloas.IndexedAttestation{} },
			root:     "683e085a809f52ba767d8e50033c25684a71708fe063224ad8ef293ecc8b55b6",
		},
		{
			name:     "ExecutionPayload",
			value:    testExecutionPayload(),
			newValue: func() sszRootable { return &gloas.ExecutionPayload{} },
			root:     "f1a113a895e90bce01951250d9858cf9d71ace6c4355f92dd409be321975d8bc",
		},
		{
			name:     "ExecutionRequests",
			value:    testExecutionRequests(),
			newValue: func() sszRootable { return &gloas.ExecutionRequests{} },
			root:     "b43a532aa4dbf7e7b6b5cfa4b6cc1c54941be8e5f4eda56d965a534a8daec3b4",
		},
		{
			name:     "BeaconBlockBody",
			value:    testBeaconBlockBody(),
			newValue: func() sszRootable { return &gloas.BeaconBlockBody{} },
			root:     "85ce69a9ad9f3abfca93126bd85b98a005e0405112b70f0bfb3e245affcf3598",
		},
		{
			name:     "BeaconState",
			value:    testBeaconState(),
			newValue: func() sszRootable { return &gloas.BeaconState{} },
			root:     "4e2b7050ea9041ac17e2308892b23f24a895302f73ae662c12641fdfb3b38431",
		},
		{
			name:     "SignedAggregateAndProof",
			value:    testSignedAggregateAndProof(),
			newValue: func() sszRootable { return &gloas.SignedAggregateAndProof{} },
			root:     "71473a74b9b3273a783edc7ebdeb412ccdbfc034dcd860dd7f77b4d506996f38",
		},
		{
			name:     "AttesterSlashing",
			value:    testAttesterSlashing(),
			newValue: func() sszRootable { return &gloas.AttesterSlashing{} },
			root:     "cbb503af73492f1ab7dbbd53ae3aa83f781e99a12015eb8b072871be5630c60c",
		},
		{
			name:     "SignedBeaconBlock",
			value:    testSignedBeaconBlock(),
			newValue: func() sszRootable { return &gloas.SignedBeaconBlock{} },
			root:     "99fff2530b1b3e26ab697f29dde7ae9f8ee0bbfdd07255861be1b33f5864a79c",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.value.MarshalSSZ()
			require.NoError(t, err)

			decoded := test.newValue()
			require.NoError(t, decoded.UnmarshalSSZ(encoded))

			roundTripped, err := decoded.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, encoded, roundTripped)

			root, err := decoded.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, rootFromHex(t, test.root), root)
		})
	}
}

// TestAttestationSSZCustomPreset exercises the dynamic path used by callers
// with custom spec support.  The minimal preset has four committees per slot,
// so its committee bitvector is one byte rather than mainnet's eight.
func TestAttestationSSZCustomPreset(t *testing.T) {
	attestation := &gloas.Attestation{
		AggregationBits: bitfield.Bitlist{0x19},
		CommitteeBits:   bitfield.Bitvector64{0},
	}
	dynamicSSZ := dynssz.NewDynSsz(map[string]any{
		"MAX_COMMITTEES_PER_SLOT": uint64(4),
	})

	encoded, err := dynamicSSZ.MarshalSSZ(attestation)
	require.NoError(t, err)
	// Independently encoded by the v1.7.0-alpha.12 minimal PySpec.
	require.Len(t, encoded, 230)

	var decoded gloas.Attestation
	require.NoError(t, dynamicSSZ.UnmarshalSSZ(&decoded, encoded))
	require.Equal(t, []byte{0x19}, []byte(decoded.AggregationBits))
	require.Equal(t, []byte{0}, []byte(decoded.CommitteeBits))

	root, err := dynamicSSZ.HashTreeRoot(&decoded)
	require.NoError(t, err)
	require.Equal(t, rootFromHex(t, "2c107e8672069142b8bfa7924cdb90be0fd9f7ebb61c6157f8f40fc608df5bf5"), root)
}

func rootFromHex(t *testing.T, input string) [32]byte {
	t.Helper()

	decoded, err := hex.DecodeString(input)
	require.NoError(t, err)
	require.Len(t, decoded, 32)

	var root [32]byte
	copy(root[:], decoded)

	return root
}

func testAttestation() *gloas.Attestation {
	aggregationBits := bitfield.NewBitlist(8)
	aggregationBits.SetBitAt(0, true)
	aggregationBits.SetBitAt(3, true)

	return &gloas.Attestation{AggregationBits: aggregationBits}
}

func testIndexedAttestation(index uint64) *gloas.IndexedAttestation {
	return &gloas.IndexedAttestation{AttestingIndices: []uint64{index}}
}

func testExecutionPayload() *gloas.ExecutionPayload {
	return &gloas.ExecutionPayload{
		BlockNumber:   1,
		BaseFeePerGas: uint256.NewInt(1),
		SlotNumber:    2,
	}
}

func testExecutionRequests() *gloas.ExecutionRequests {
	return &gloas.ExecutionRequests{
		BuilderExits: []*gloas.BuilderExitRequest{{
			SourceAddress: executionAddress(1),
			Pubkey:        publicKey(2),
		}},
	}
}

func testBeaconBlockBody() *gloas.BeaconBlockBody {
	return &gloas.BeaconBlockBody{
		AttesterSlashings: []*gloas.AttesterSlashing{testAttesterSlashing()},
		Attestations:      []*gloas.Attestation{testAttestation()},
	}
}

func testBeaconState() *gloas.BeaconState {
	return &gloas.BeaconState{
		GenesisTime: 1,
		Builders: []*gloas.Builder{{
			PublicKey: publicKey(3),
			Version:   1,
		}},
		NextWithdrawalBuilderIndex: 2,
	}
}

func testSignedAggregateAndProof() *gloas.SignedAggregateAndProof {
	return &gloas.SignedAggregateAndProof{
		Message: &gloas.AggregateAndProof{
			AggregatorIndex: 7,
			Aggregate:       testAttestation(),
		},
	}
}

func testAttesterSlashing() *gloas.AttesterSlashing {
	return &gloas.AttesterSlashing{
		Attestation1: testIndexedAttestation(1),
		Attestation2: testIndexedAttestation(2),
	}
}

func testSignedBeaconBlock() *gloas.SignedBeaconBlock {
	return &gloas.SignedBeaconBlock{
		Message: &gloas.BeaconBlock{
			Slot: 3,
			Body: testBeaconBlockBody(),
		},
	}
}

func executionAddress(value byte) (address bellatrix.ExecutionAddress) {
	for i := range address {
		address[i] = value
	}

	return address
}

func publicKey(value byte) (pubkey phase0.BLSPubKey) {
	for i := range pubkey {
		pubkey[i] = value
	}

	return pubkey
}
