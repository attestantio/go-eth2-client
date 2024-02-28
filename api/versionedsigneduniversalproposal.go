// Copyright © 2024 Attestant Limited.
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

// VersionedSignedUniversalProposal contains a full or blinded signed proposal.
type VersionedSignedUniversalProposal struct {
	Proposal        *VersionedSignedProposal
	BlindedProposal *VersionedSignedBlindedProposal
}

// Slot returns the slot of the signed proposal.
func (v *VersionedSignedUniversalProposal) Slot() (phase0.Slot, error) {
	if v.Proposal != nil {
		return v.Proposal.Slot()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Slot()
	}
	return 0, ErrDataMissing
}

// ExecutionBlockHash returns the hash of the execution payload.
func (v *VersionedSignedUniversalProposal) ExecutionBlockHash() (phase0.Hash32, error) {
	if v.Proposal != nil {
		return v.Proposal.ExecutionBlockHash()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.ExecutionBlockHash()
	}
	return phase0.Hash32{}, ErrDataMissing
}

// String returns a string version of the structure.
func (v *VersionedSignedUniversalProposal) String() string {
	if v.Proposal != nil {
		return v.Proposal.String()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.String()
	}
	return "unsupported version"
}
