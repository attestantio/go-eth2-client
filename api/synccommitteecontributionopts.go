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

// SyncCommitteeContributionOpts are the options for obtaining sync committee contributions.
type SyncCommitteeContributionOpts struct {
	Common CommonOpts

	// Slot is the slot for which the data is obtained.
	Slot phase0.Slot
	// SubcommitteeIndex is the index of the sync subcommittee for which the data is obtained.
	SubcommitteeIndex uint64
	// BeaconBlockRoot is the root of the beacon block for which the data is obtained.
	BeaconBlockRoot phase0.Root
}
