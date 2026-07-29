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
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// EPBSProposalOpts are the options for obtaining ePBS block proposals.
//
// Post-Gloas there is no blinded proposal: proposers commit to an execution
// payload bid rather than to a payload, so the axis that replaced blinded versus
// unblinded is IncludePayload — whether the execution payload envelope travels
// with the block or stays cached on the producing node.
type EPBSProposalOpts struct {
	Common CommonOpts

	// Slot is the slot for which the proposal should be fetched.
	Slot phase0.Slot
	// RandaoReveal is the RANDAO reveal for the proposal.
	RandaoReveal phase0.BLSSignature
	// Graffiti is the graffiti to be included in the beacon block body.
	Graffiti [32]byte
	// SkipRandaoVerification is true if we do not want the server to verify our RANDAO reveal.
	// If this is set then the RANDAO reveal should be passed as the point at infinity (0xc0…00)
	SkipRandaoVerification bool
	// IncludePayload selects whether a self-built proposal carries its execution
	// payload envelope and blobs.
	//
	// True is stateless: everything needed to publish travels in the response, so
	// any beacon node can publish the block.  This is what multiple-beacon-node,
	// distributed-validator and failover setups require.
	//
	// False is stateful: the producing node caches the envelope and blobs, so the
	// block must be published via that same node, and the envelope retrieved from
	// it with ExecutionPayloadEnvelope().
	//
	// The spec marks this parameter required with no default, and the two choices
	// carry materially different operational constraints, so it is a pointer:
	// there is no safe value to assume on the caller's behalf, and leaving it
	// unset is rejected rather than silently resolved to the constraining mode.
	IncludePayload *bool
	// BuilderBoostFactor is the relative weight of the builder payload versus a locally-produced
	// payload, as per https://ethereum.github.io/beacon-APIs/#/Validator/produceBlockV4
	// This is optional; if not supplied it will use the default value of 100.
	BuilderBoostFactor *uint64
}
