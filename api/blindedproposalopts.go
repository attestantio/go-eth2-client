// Copyright © 2023 Attestant Limited.
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

// BlindedProposalOpts are the options for obtaining blinded proposals.
type BlindedProposalOpts struct {
	Common CommonOpts

	// Slot is the slot for which the proposal should be fetched.
	Slot phase0.Slot
	// RandaoReveal is the RANDAO reveal for the proposal.
	RandaoReveal phase0.BLSSignature
	// Graffit is the graffiti to be included in the beacon block body.
	Graffiti [32]byte
	// SkipRandaoVerification is true if we do not want the server to verify our RANDAO reveal.
	// If this is set then the RANDAO reveal should be passed as the point at infinity (0xc0…00)
	SkipRandaoVerification bool
}
