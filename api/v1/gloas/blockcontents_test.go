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
	"encoding/json"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	require "github.com/stretchr/testify/require"
)

// validBeaconBlock returns a block populated just far enough for its own JSON
// marshaler to run: every pointer the marshaler dereferences is non-nil.
func validBeaconBlock() *gloas.BeaconBlock {
	return &gloas.BeaconBlock{
		Body: &gloas.BeaconBlockBody{
			ETH1Data:                  &phase0.ETH1Data{BlockHash: make([]byte, phase0.Hash32Length)},
			SyncAggregate:             &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			ParentExecutionRequests:   &gloas.ExecutionRequests{},
		},
	}
}

// validExecutionPayloadEnvelope returns an envelope whose payload carries a
// non-nil BaseFeePerGas, which the payload's own marshaler dereferences.
func validExecutionPayloadEnvelope() *gloas.ExecutionPayloadEnvelope {
	return &gloas.ExecutionPayloadEnvelope{
		Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(0)},
		ExecutionRequests: &gloas.ExecutionRequests{},
	}
}

// validBlockContents returns block contents carrying one KZG proof and no
// blobs.  Blobs are left empty deliberately: deneb.Blob is a 128KiB array, so a
// populated one turns any require.Equal failure into a 256KiB hex diff.  The
// element path for blobs is covered separately by
// TestBlockContentsJSONPreservesBlobs.
func validBlockContents() *apiv1gloas.BlockContents {
	return &apiv1gloas.BlockContents{
		Block:                    validBeaconBlock(),
		ExecutionPayloadEnvelope: validExecutionPayloadEnvelope(),
		KZGProofs:                []deneb.KZGProof{{0x01, 0x02, 0x03}},
		Blobs:                    []deneb.Blob{},
	}
}

// TestBlockContentsJSON verifies that BlockContents marshals to the four field
// names the beacon-APIs Gloas.BlockContents schema requires, and that a
// populated value survives a JSON round trip.  The schema names all four in
// snake_case; Go's default marshaler would emit the exported field names
// instead, so the key assertions fail against an uncodec'd struct.
func TestBlockContentsJSON(t *testing.T) {
	contents := validBlockContents()

	data, err := json.Marshal(contents)
	require.NoError(t, err)

	for _, key := range []string{
		`"block":`,
		`"execution_payload_envelope":`,
		`"kzg_proofs":`,
		`"blobs":`,
	} {
		require.Contains(t, string(data), key)
	}

	var decoded apiv1gloas.BlockContents
	require.NoError(t, json.Unmarshal(data, &decoded))

	// Compare on the wire rather than with require.Equal on the structs: the
	// nested gloas.ExecutionPayload codec normalises a nil Transactions,
	// ExtraData or BlockAccessList to an empty slice on the way back in, so
	// struct equality would assert that package's conventions rather than this
	// codec's behaviour.  Byte equality still catches the failure that matters —
	// drop any one field from UnmarshalJSON and its key re-marshals as null.
	reencoded, err := json.Marshal(&decoded)
	require.NoError(t, err)
	require.Equal(t, string(data), string(reencoded))

	// The fields this type owns directly must survive intact.
	require.Equal(t, contents.KZGProofs, decoded.KZGProofs)
	require.Equal(t, contents.Blobs, decoded.Blobs)
}

// TestBlockContentsSSZ verifies that BlockContents round-trips through SSZ.  The
// generated codec is what the block-production endpoint uses when a node answers
// with application/octet-stream, so an absent or misordered codec would surface
// only as an undecodable proposal.
func TestBlockContentsSSZ(t *testing.T) {
	contents := validBlockContents()

	data, err := contents.MarshalSSZ()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var decoded apiv1gloas.BlockContents
	require.NoError(t, decoded.UnmarshalSSZ(data))

	// SSZ carries no field names, so a wrongly ordered or wrongly sized codec
	// shows up as a differing hash tree root rather than a decode error.
	want, err := contents.HashTreeRoot()
	require.NoError(t, err)
	got, err := decoded.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, want, got)

	require.Equal(t, contents.KZGProofs, decoded.KZGProofs)
}

// TestBlockContentsString verifies that stringifying block contents renders the
// spec's field names.  goccy/go-yaml would otherwise lower-case the Go field
// names, emitting executionpayloadenvelope rather than the spec's
// execution_payload_envelope, so this pins the hand-written YAML marshaler.
func TestBlockContentsString(t *testing.T) {
	var str string
	require.NotPanics(t, func() {
		str = validBlockContents().String()
	})

	require.Contains(t, str, "execution_payload_envelope:")
	require.Contains(t, str, "kzg_proofs:")
}

// TestBlockContentsJSONMissingField verifies that every one of the four fields
// the schema marks required is rejected when absent.  encoding/json ignores
// unknown keys and leaves absent ones at their zero value, so without the
// presence check a body missing its blobs would decode as success and hand the
// caller contents it cannot publish.
func TestBlockContentsJSONMissingField(t *testing.T) {
	tests := []struct {
		name    string
		omitKey string
		err     string
	}{
		{name: "Block", omitKey: "block", err: "block: missing"},
		{
			name:    "ExecutionPayloadEnvelope",
			omitKey: "execution_payload_envelope",
			err:     "execution_payload_envelope: missing",
		},
		{name: "KZGProofs", omitKey: "kzg_proofs", err: "kzg_proofs: missing"},
		{name: "Blobs", omitKey: "blobs", err: "blobs: missing"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(validBlockContents())
			require.NoError(t, err)

			var generic map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(data, &generic))
			delete(generic, test.omitKey)

			pruned, err := json.Marshal(generic)
			require.NoError(t, err)

			var decoded apiv1gloas.BlockContents
			require.EqualError(t, decoded.UnmarshalJSON(pruned), test.err)
		})
	}
}

// TestBlockContentsJSONPreservesBlobs verifies that a populated blob survives a
// JSON round trip end to end.  A blob is a 128KiB fixed-size array, so the
// assertions check its first and last byte rather than comparing whole values:
// that catches a truncating or zero-padding copy without producing a 256KiB
// diff on failure.
func TestBlockContentsJSONPreservesBlobs(t *testing.T) {
	contents := validBlockContents()

	var blob deneb.Blob
	blob[0] = 0xaa
	blob[len(blob)-1] = 0xbb
	contents.Blobs = []deneb.Blob{blob}

	data, err := json.Marshal(contents)
	require.NoError(t, err)

	var decoded apiv1gloas.BlockContents
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Len(t, decoded.Blobs, 1)
	require.Equal(t, byte(0xaa), decoded.Blobs[0][0])
	require.Equal(t, byte(0xbb), decoded.Blobs[0][len(decoded.Blobs[0])-1])
}
