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

package spec_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestVersionedSignedExecutionPayloadBidAccessors(t *testing.T) {
	parentBlockHash := phase0.Hash32{0x01}
	parentBlockRoot := phase0.Root{0x02}
	blockHash := phase0.Hash32{0x03}
	feeRecipient := bellatrix.ExecutionAddress{0x04}
	blobKZGCommitments := []deneb.KZGCommitment{{0x05}}
	signature := phase0.BLSSignature{0x06}
	populated := &spec.VersionedSignedExecutionPayloadBid{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.SignedExecutionPayloadBid{
			Message: &gloas.ExecutionPayloadBid{
				ParentBlockHash:    parentBlockHash,
				ParentBlockRoot:    parentBlockRoot,
				BlockHash:          blockHash,
				FeeRecipient:       feeRecipient,
				GasLimit:           7,
				BuilderIndex:       8,
				Slot:               9,
				Value:              10,
				ExecutionPayment:   11,
				BlobKZGCommitments: blobKZGCommitments,
			},
			Signature: signature,
		},
	}
	gloasNoMessage := &spec.VersionedSignedExecutionPayloadBid{
		Version: spec.DataVersionGloas,
		Gloas:   &gloas.SignedExecutionPayloadBid{},
	}
	gloasNil := &spec.VersionedSignedExecutionPayloadBid{Version: spec.DataVersionGloas}
	tests := []struct {
		name     string
		accessor func(*spec.VersionedSignedExecutionPayloadBid) (any, error)
		want     any
		missing  *spec.VersionedSignedExecutionPayloadBid
	}{
		{
			name: "ParentBlockHash",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.ParentBlockHash()
			},
			want:    parentBlockHash,
			missing: gloasNoMessage,
		},
		{
			name: "ParentBlockRoot",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.ParentBlockRoot()
			},
			want:    parentBlockRoot,
			missing: gloasNoMessage,
		},
		{
			name: "BlockHash",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.BlockHash()
			},
			want:    blockHash,
			missing: gloasNoMessage,
		},
		{
			name: "FeeRecipient",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.FeeRecipient()
			},
			want:    feeRecipient,
			missing: gloasNoMessage,
		},
		{
			name: "GasLimit",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.GasLimit()
			},
			want:    uint64(7),
			missing: gloasNoMessage,
		},
		{
			name: "BuilderIndex",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.BuilderIndex()
			},
			want:    gloas.BuilderIndex(8),
			missing: gloasNoMessage,
		},
		{
			name: "Slot",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.Slot()
			},
			want:    phase0.Slot(9),
			missing: gloasNoMessage,
		},
		{
			name: "Value",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.Value()
			},
			want:    phase0.Gwei(10),
			missing: gloasNoMessage,
		},
		{
			name: "ExecutionPayment",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.ExecutionPayment()
			},
			want:    phase0.Gwei(11),
			missing: gloasNoMessage,
		},
		{
			name: "BlobKZGCommitments",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.BlobKZGCommitments()
			},
			want:    blobKZGCommitments,
			missing: gloasNoMessage,
		},
		{
			name: "Signature",
			accessor: func(v *spec.VersionedSignedExecutionPayloadBid) (any, error) {
				return v.Signature()
			},
			want:    signature,
			missing: gloasNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.accessor(populated)
			require.NoError(t, err)
			require.Equal(t, test.want, got)

			_, err = test.accessor(test.missing)
			require.EqualError(t, err, "no gloas execution payload bid")
		})
	}
}

func TestVersionedSignedExecutionPayloadBidString(t *testing.T) {
	tests := []struct {
		name     string
		bid      *spec.VersionedSignedExecutionPayloadBid
		expected string
	}{
		{
			name: "Gloas",
			bid: &spec.VersionedSignedExecutionPayloadBid{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedExecutionPayloadBid{
					Message: &gloas.ExecutionPayloadBid{},
				},
			},
			expected: "parent_block_hash",
		},
		{
			name: "GloasNil",
			bid: &spec.VersionedSignedExecutionPayloadBid{
				Version: spec.DataVersionGloas,
			},
		},
		{
			name:     "UnknownVersion",
			bid:      &spec.VersionedSignedExecutionPayloadBid{},
			expected: "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expected == "" {
				require.Empty(t, test.bid.String())
			} else {
				require.Contains(t, test.bid.String(), test.expected)
			}
		})
	}
}
