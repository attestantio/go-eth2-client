// Copyright © 2023, 2024 Attestant Limited.
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

import "github.com/attestantio/go-eth2-client/spec/phase0"

// ProposalOpts are the options for obtaining proposals.
type ProposalOpts struct {
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
	// BuilderBoostFactor is the relative weight of the builder payload versus a locally-produced
	// payload, as per https://ethereum.github.io/beacon-APIs/#/Validator/produceBlockV3
	// This is optional; if not supplied it will use the default value of 100.
	BuilderBoostFactor *uint64
}
