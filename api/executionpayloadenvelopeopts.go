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

// ExecutionPayloadEnvelopeOpts are the options for obtaining the execution
// payload envelope a node cached while producing a block.
//
// This is the other half of asking for a proposal with the payload excluded: the
// node kept the envelope, and this is how the caller retrieves it to sign and
// publish.
type ExecutionPayloadEnvelopeOpts struct {
	Common CommonOpts

	// Slot is the slot for which the envelope was built.
	Slot phase0.Slot
	// BeaconBlockRoot is the root of the beacon block the envelope commits to.
	//
	// Passing it makes the lookup re-org resistant: if a re-org happens between
	// producing the block and retrieving the envelope, the caller gets nothing
	// rather than an envelope committed to a block that is no longer on the
	// chain.  The node is expected to check the cached envelope against it, and
	// the returned envelope's own root is checked against it again here, so a
	// node that does not is caught rather than trusted.
	BeaconBlockRoot phase0.Root
}
