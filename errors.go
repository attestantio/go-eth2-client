// Copyright © 2024 - 2026 Attestant Limited.
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

package client

import "errors"

var (
	// ErrNotActive is returned when a client is not active.
	ErrNotActive = errors.New("client is not active")
	// ErrNotSynced is returned when a client is not synced.
	ErrNotSynced = errors.New("client is not synced")
	// ErrNoOptions is returned when a request is made without options.
	ErrNoOptions = errors.New("no options specified")
	// ErrInvalidOptions is returned when a request is made with invalid options.
	ErrInvalidOptions = errors.New("invalid options")
	// ErrInconsistentResult is returned when a request returns with data at odds to that requested.
	ErrInconsistentResult = errors.New("inconsistent result")
	// ErrNoPayloadAttestationData is returned when the node has seen no block for the requested
	// slot, and so the validator must not cast a payload attestation for it.  This is a normal
	// outcome rather than a failure, but it is reported as an error so that a caller cannot
	// mistake it for a datum that is safe to sign.
	ErrNoPayloadAttestationData = errors.New("no payload attestation data available")
	// ErrNoExecutionPayloadEnvelope is returned when the node holds no cached execution payload
	// envelope for the requested slot and block root.  A node caches only the envelope it built
	// for the slot it is proposing, so this is the normal answer for any other slot, for a block
	// the node did not build, and for a block that a re-org has moved off the chain.
	ErrNoExecutionPayloadEnvelope = errors.New("no execution payload envelope available")
)
