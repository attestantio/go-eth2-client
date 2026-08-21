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

package api

import (
	apiv2 "github.com/attestantio/go-eth2-client/api/v2"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// SubmitExecutionPayloadEnvelopeOpts are the options for submitting execution
// payload envelopes. The submission always uses the stateless request form
// (SignedExecutionPayloadEnvelopeContents): the signed envelope plus the
// blobs and KZG cell proofs the beacon node needs to build data column
// sidecars. The stateful/blinded form only works when the beacon node cached
// the full envelope from its own block production, so it is not offered here.
type SubmitExecutionPayloadEnvelopeOpts struct {
	Common CommonOpts

	// SignedExecutionPayloadEnvelope is the signed envelope to publish.
	// Its Version selects the wire schema and the consensus version header.
	SignedExecutionPayloadEnvelope *spec.VersionedSignedExecutionPayloadEnvelope

	// KZGProofs are the cell KZG proofs for the blobs committed to by the
	// envelope's payload.
	KZGProofs []deneb.KZGProof

	// Blobs are the blobs committed to by the envelope's payload.
	Blobs []deneb.Blob

	// BroadcastValidation is the validation required of the consensus node
	// before broadcasting the envelope.
	BroadcastValidation *apiv2.BroadcastValidation
}
