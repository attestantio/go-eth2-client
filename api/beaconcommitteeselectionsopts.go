// Copyright © 2025 Attestant Limited.
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

import apiv1 "github.com/attestantio/go-eth2-client/api/v1"

// BeaconCommitteeSelectionsOpts are the options for obtaining beacon committee selections.
type BeaconCommitteeSelectionsOpts struct {
	Common CommonOpts

	// Beacon Committee Selections are the selections which the DV should resolve.
	Selections []*apiv1.BeaconCommitteeSelection
}
