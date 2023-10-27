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

// BeaconCommitteesOpts are the options for obtaining proposer duties.
type BeaconCommitteesOpts struct {
	Common CommonOpts

	// State is the state at which the data is obtained.
	// It can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
	State string
	// Epoch is the epoch for which the data is obtained.
	// This is optional; if not supplied it will obtain the data at the epoch relating to the state.
	Epoch *phase0.Epoch
}
