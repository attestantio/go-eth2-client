// Copyright © 2022, 2024 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedUniversalProposal contains a versioned universal proposal.
type VersionedUniversalProposal struct {
	Proposal        *VersionedProposal
	BlindedProposal *VersionedBlindedProposal
}

// IsEmpty returns true if there is no proposal.
func (v VersionedUniversalProposal) IsEmpty() bool {
	noFull := v.Proposal == nil || v.Proposal.IsEmpty()
	noBlinded := v.BlindedProposal == nil || v.BlindedProposal.IsEmpty()
	return noFull && noBlinded
}

// Slot returns the slot of the proposal.
func (v *VersionedUniversalProposal) Slot() (phase0.Slot, error) {
	if v.Proposal != nil {
		return v.Proposal.Slot()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Slot()
	}
	return 0, ErrDataMissing
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedUniversalProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	if v.Proposal != nil {
		return v.Proposal.ProposerIndex()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.ProposerIndex()
	}
	return 0, ErrDataMissing
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedUniversalProposal) RandaoReveal() (phase0.BLSSignature, error) {
	if v.Proposal != nil {
		return v.Proposal.RandaoReveal()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.RandaoReveal()
	}
	return phase0.BLSSignature{}, ErrDataMissing
}

// Graffiti returns the graffiti of the proposal.
func (v *VersionedUniversalProposal) Graffiti() ([32]byte, error) {
	if v.Proposal != nil {
		return v.Proposal.Graffiti()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Graffiti()
	}
	return [32]byte{}, ErrDataMissing
}

// Attestations returns the attestations of the proposal.
func (v *VersionedUniversalProposal) Attestations() ([]*phase0.Attestation, error) {
	if v.Proposal != nil {
		return v.Proposal.Attestations()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Attestations()
	}
	return nil, ErrDataMissing
}

// Root returns the root of the proposal.
func (v *VersionedUniversalProposal) Root() (phase0.Root, error) {
	if v.Proposal != nil {
		return v.Proposal.Root()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Root()
	}
	return phase0.Root{}, ErrDataMissing
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedUniversalProposal) BodyRoot() (phase0.Root, error) {
	if v.Proposal != nil {
		return v.Proposal.BodyRoot()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.BodyRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedUniversalProposal) ParentRoot() (phase0.Root, error) {
	if v.Proposal != nil {
		return v.Proposal.ParentRoot()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.ParentRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// StateRoot returns the state root of the proposal.
func (v *VersionedUniversalProposal) StateRoot() (phase0.Root, error) {
	if v.Proposal != nil {
		return v.Proposal.StateRoot()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.StateRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// Transactions returns the transactions of the proposal.
func (v *VersionedUniversalProposal) Transactions() ([]bellatrix.Transaction, error) {
	if v.Proposal != nil {
		return v.Proposal.Transactions()
	}
	return nil, ErrDataMissing
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedUniversalProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	if v.Proposal != nil {
		return v.Proposal.FeeRecipient()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.FeeRecipient()
	}
	return bellatrix.ExecutionAddress{}, ErrDataMissing
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedUniversalProposal) Timestamp() (uint64, error) {
	if v.Proposal != nil {
		return v.Proposal.Timestamp()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.Timestamp()
	}
	return 0, ErrDataMissing
}

// Blobs returns the blobs of the proposal.
func (v *VersionedUniversalProposal) Blobs() ([]deneb.Blob, error) {
	if v.Proposal != nil {
		return v.Proposal.Blobs()
	}
	return nil, ErrDataMissing
}

// KZGProofs returns the KZG proofs of the proposal.
func (v *VersionedUniversalProposal) KZGProofs() ([]deneb.KZGProof, error) {
	if v.Proposal != nil {
		return v.Proposal.KZGProofs()
	}
	return nil, ErrDataMissing
}

// String returns a string version of the structure.
func (v *VersionedUniversalProposal) String() string {
	if v.Proposal != nil {
		return v.Proposal.String()
	}
	if v.BlindedProposal != nil {
		return v.BlindedProposal.String()
	}
	return ""
}
