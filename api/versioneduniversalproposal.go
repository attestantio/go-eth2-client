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
	Full    *VersionedProposal
	Blinded *VersionedBlindedProposal
}

// IsEmpty returns true if there is no proposal.
func (v VersionedUniversalProposal) IsEmpty() bool {
	noFull := v.Full == nil || v.Full.IsEmpty()
	noBlinded := v.Blinded == nil || v.Blinded.IsEmpty()
	return noFull && noBlinded
}

// Slot returns the slot of the proposal.
func (v *VersionedUniversalProposal) Slot() (phase0.Slot, error) {
	if v.Full != nil {
		return v.Full.Slot()
	}
	if v.Blinded != nil {
		return v.Blinded.Slot()
	}
	return 0, ErrDataMissing
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedUniversalProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	if v.Full != nil {
		return v.Full.ProposerIndex()
	}
	if v.Blinded != nil {
		return v.Blinded.ProposerIndex()
	}
	return 0, ErrDataMissing
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedUniversalProposal) RandaoReveal() (phase0.BLSSignature, error) {
	if v.Full != nil {
		return v.Full.RandaoReveal()
	}
	if v.Blinded != nil {
		return v.Blinded.RandaoReveal()
	}
	return phase0.BLSSignature{}, ErrDataMissing
}

// Graffiti returns the graffiti of the proposal.
func (v *VersionedUniversalProposal) Graffiti() ([32]byte, error) {
	if v.Full != nil {
		return v.Full.Graffiti()
	}
	if v.Blinded != nil {
		return v.Blinded.Graffiti()
	}
	return [32]byte{}, ErrDataMissing
}

// Attestations returns the attestations of the proposal.
func (v *VersionedUniversalProposal) Attestations() ([]*phase0.Attestation, error) {
	if v.Full != nil {
		return v.Full.Attestations()
	}
	if v.Blinded != nil {
		return v.Blinded.Attestations()
	}
	return nil, ErrDataMissing
}

// Root returns the root of the proposal.
func (v *VersionedUniversalProposal) Root() (phase0.Root, error) {
	if v.Full != nil {
		return v.Full.Root()
	}
	if v.Blinded != nil {
		return v.Blinded.Root()
	}
	return phase0.Root{}, ErrDataMissing
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedUniversalProposal) BodyRoot() (phase0.Root, error) {
	if v.Full != nil {
		return v.Full.BodyRoot()
	}
	if v.Blinded != nil {
		return v.Blinded.BodyRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedUniversalProposal) ParentRoot() (phase0.Root, error) {
	if v.Full != nil {
		return v.Full.ParentRoot()
	}
	if v.Blinded != nil {
		return v.Blinded.ParentRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// StateRoot returns the state root of the proposal.
func (v *VersionedUniversalProposal) StateRoot() (phase0.Root, error) {
	if v.Full != nil {
		return v.Full.StateRoot()
	}
	if v.Blinded != nil {
		return v.Blinded.StateRoot()
	}
	return phase0.Root{}, ErrDataMissing
}

// Transactions returns the transactions of the proposal.
func (v *VersionedUniversalProposal) Transactions() ([]bellatrix.Transaction, error) {
	if v.Full != nil {
		return v.Full.Transactions()
	}
	return nil, ErrDataMissing
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedUniversalProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	if v.Full != nil {
		return v.Full.FeeRecipient()
	}
	if v.Blinded != nil {
		return v.Blinded.FeeRecipient()
	}
	return bellatrix.ExecutionAddress{}, ErrDataMissing
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedUniversalProposal) Timestamp() (uint64, error) {
	if v.Full != nil {
		return v.Full.Timestamp()
	}
	if v.Blinded != nil {
		return v.Blinded.Timestamp()
	}
	return 0, ErrDataMissing
}

// Blobs returns the blobs of the proposal.
func (v *VersionedUniversalProposal) Blobs() ([]deneb.Blob, error) {
	if v.Full != nil {
		return v.Full.Blobs()
	}
	return nil, ErrDataMissing
}

// KZGProofs returns the KZG proofs of the proposal.
func (v *VersionedUniversalProposal) KZGProofs() ([]deneb.KZGProof, error) {
	if v.Full != nil {
		return v.Full.KZGProofs()
	}
	return nil, ErrDataMissing
}

// String returns a string version of the structure.
func (v *VersionedUniversalProposal) String() string {
	if v.Full != nil {
		return v.Full.String()
	}
	if v.Blinded != nil {
		return v.Blinded.String()
	}
	return ""
}
