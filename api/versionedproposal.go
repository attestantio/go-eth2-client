// Copyright © 2022, 2023 Attestant Limited.
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
	"math/big"

	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"

	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedProposal contains a versioned proposal.
type VersionedProposal struct {
	Version          spec.DataVersion
	Blinded          bool
	ConsensusValue   *big.Int
	ExecutionValue   *big.Int
	Phase0           *phase0.BeaconBlock
	Altair           *altair.BeaconBlock
	Bellatrix        *bellatrix.BeaconBlock
	BellatrixBlinded *apiv1bellatrix.BlindedBeaconBlock
	Capella          *capella.BeaconBlock
	CapellaBlinded   *apiv1capella.BlindedBeaconBlock
	Deneb            *apiv1deneb.BlockContents
	DenebBlinded     *apiv1deneb.BlindedBeaconBlock
	Electra          *apiv1electra.BlockContents
	ElectraBlinded   *apiv1electra.BlindedBeaconBlock
}

// IsEmpty returns true if there is no proposal.
func (v *VersionedProposal) IsEmpty() bool {
	return v.Phase0 == nil &&
		v.Altair == nil &&
		v.Bellatrix == nil &&
		v.BellatrixBlinded == nil &&
		v.Capella == nil &&
		v.CapellaBlinded == nil &&
		v.Deneb == nil &&
		v.DenebBlinded == nil &&
		v.Electra == nil &&
		v.ElectraBlinded == nil
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedProposal) BodyRoot() (phase0.Root, error) {
	if !v.bodyPresent() {
		return phase0.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Body.HashTreeRoot()
	case spec.DataVersionAltair:
		return v.Altair.Body.HashTreeRoot()
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.HashTreeRoot()
		}

		return v.Bellatrix.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.HashTreeRoot()
		}

		return v.Capella.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.HashTreeRoot()
		}

		return v.Deneb.Block.Body.HashTreeRoot()
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.HashTreeRoot()
		}

		return v.Electra.Block.Body.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedProposal) ParentRoot() (phase0.Root, error) {
	if !v.proposalPresent() {
		return phase0.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.ParentRoot, nil
	case spec.DataVersionAltair:
		return v.Altair.ParentRoot, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.ParentRoot, nil
		}

		return v.Bellatrix.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.ParentRoot, nil
		}

		return v.Capella.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.ParentRoot, nil
		}

		return v.Deneb.Block.ParentRoot, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.ParentRoot, nil
		}

		return v.Electra.Block.ParentRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	if !v.proposalPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.ProposerIndex, nil
	case spec.DataVersionAltair:
		return v.Altair.ProposerIndex, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.ProposerIndex, nil
		}

		return v.Bellatrix.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.ProposerIndex, nil
		}

		return v.Capella.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.ProposerIndex, nil
		}

		return v.Deneb.Block.ProposerIndex, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.ProposerIndex, nil
		}

		return v.Electra.Block.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Root returns the root of the proposal.
func (v *VersionedProposal) Root() (phase0.Root, error) {
	if !v.proposalPresent() {
		return phase0.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.HashTreeRoot()
	case spec.DataVersionAltair:
		return v.Altair.HashTreeRoot()
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.HashTreeRoot()
		}

		return v.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.HashTreeRoot()
		}

		return v.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.HashTreeRoot()
		}

		return v.Deneb.Block.HashTreeRoot()
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.HashTreeRoot()
		}

		return v.Electra.Block.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// Slot returns the slot of the proposal.
func (v *VersionedProposal) Slot() (phase0.Slot, error) {
	if !v.proposalPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Slot, nil
	case spec.DataVersionAltair:
		return v.Altair.Slot, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Slot, nil
		}

		return v.Bellatrix.Slot, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Slot, nil
		}

		return v.Capella.Slot, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Slot, nil
		}

		return v.Deneb.Block.Slot, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Slot, nil
		}

		return v.Electra.Block.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the proposal.
func (v *VersionedProposal) StateRoot() (phase0.Root, error) {
	if !v.proposalPresent() {
		return phase0.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.StateRoot, nil
	case spec.DataVersionAltair:
		return v.Altair.StateRoot, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.StateRoot, nil
		}

		return v.Bellatrix.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.StateRoot, nil
		}

		return v.Capella.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.StateRoot, nil
		}

		return v.Deneb.Block.StateRoot, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.StateRoot, nil
		}

		return v.Electra.Block.StateRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the proposal.
func (v *VersionedProposal) Attestations() ([]spec.VersionedAttestation, error) {
	if !v.bodyPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		versionedAttestations := make([]spec.VersionedAttestation, len(v.Phase0.Body.Attestations))
		for i, attestation := range v.Phase0.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionPhase0,
				Phase0:  attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionAltair:
		versionedAttestations := make([]spec.VersionedAttestation, len(v.Altair.Body.Attestations))
		for i, attestation := range v.Altair.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionAltair,
				Altair:  attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			versionedAttestations := make([]spec.VersionedAttestation, len(v.BellatrixBlinded.Body.Attestations))
			for i, attestation := range v.BellatrixBlinded.Body.Attestations {
				versionedAttestations[i] = spec.VersionedAttestation{
					Version:   spec.DataVersionBellatrix,
					Bellatrix: attestation,
				}
			}

			return versionedAttestations, nil
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Bellatrix.Body.Attestations))
		for i, attestation := range v.Bellatrix.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version:   spec.DataVersionBellatrix,
				Bellatrix: attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			versionedAttestations := make([]spec.VersionedAttestation, len(v.CapellaBlinded.Body.Attestations))
			for i, attestation := range v.CapellaBlinded.Body.Attestations {
				versionedAttestations[i] = spec.VersionedAttestation{
					Version: spec.DataVersionCapella,
					Capella: attestation,
				}
			}

			return versionedAttestations, nil
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Capella.Body.Attestations))
		for i, attestation := range v.Capella.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionCapella,
				Capella: attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			versionedAttestations := make([]spec.VersionedAttestation, len(v.DenebBlinded.Body.Attestations))
			for i, attestation := range v.DenebBlinded.Body.Attestations {
				versionedAttestations[i] = spec.VersionedAttestation{
					Version: spec.DataVersionDeneb,
					Deneb:   attestation,
				}
			}

			return versionedAttestations, nil
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Deneb.Block.Body.Attestations))
		for i, attestation := range v.Deneb.Block.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionDeneb,
				Deneb:   attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			versionedAttestations := make([]spec.VersionedAttestation, len(v.ElectraBlinded.Body.Attestations))
			for i, attestation := range v.ElectraBlinded.Body.Attestations {
				versionedAttestations[i] = spec.VersionedAttestation{
					Version: spec.DataVersionElectra,
					Electra: attestation,
				}
			}

			return versionedAttestations, nil
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Electra.Block.Body.Attestations))
		for i, attestation := range v.Electra.Block.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionElectra,
				Electra: attestation,
			}
		}

		return versionedAttestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the proposal.
func (v *VersionedProposal) Graffiti() ([32]byte, error) {
	if !v.bodyPresent() {
		return [32]byte{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Body.Graffiti, nil
	case spec.DataVersionAltair:
		return v.Altair.Body.Graffiti, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.Graffiti, nil
		}

		return v.Bellatrix.Body.Graffiti, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.Graffiti, nil
		}

		return v.Capella.Body.Graffiti, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.Graffiti, nil
		}

		return v.Deneb.Block.Body.Graffiti, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.Graffiti, nil
		}

		return v.Electra.Block.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedProposal) RandaoReveal() (phase0.BLSSignature, error) {
	if !v.bodyPresent() {
		return phase0.BLSSignature{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Body.RANDAOReveal, nil
	case spec.DataVersionAltair:
		return v.Altair.Body.RANDAOReveal, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.RANDAOReveal, nil
		}

		return v.Bellatrix.Body.RANDAOReveal, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.RANDAOReveal, nil
		}

		return v.Capella.Body.RANDAOReveal, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.RANDAOReveal, nil
		}

		return v.Deneb.Block.Body.RANDAOReveal, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.RANDAOReveal, nil
		}

		return v.Electra.Block.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Transactions returns the transactions of the proposal.
func (v *VersionedProposal) Transactions() ([]bellatrix.Transaction, error) {
	if v.Version >= spec.DataVersionBellatrix && !v.payloadPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Electra.Block.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	if v.Version >= spec.DataVersionBellatrix && !v.payloadPresent() {
		return bellatrix.ExecutionAddress{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.ExecutionPayloadHeader.FeeRecipient, nil
		}

		return v.Bellatrix.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.FeeRecipient, nil
		}

		return v.Capella.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.ExecutionPayloadHeader.FeeRecipient, nil
		}

		return v.Deneb.Block.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.ExecutionPayloadHeader.FeeRecipient, nil
		}

		return v.Electra.Block.Body.ExecutionPayload.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedProposal) Timestamp() (uint64, error) {
	if v.Version >= spec.DataVersionBellatrix && !v.payloadPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.ExecutionPayloadHeader.Timestamp, nil
		}

		return v.Bellatrix.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.Timestamp, nil
		}

		return v.Capella.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.ExecutionPayloadHeader.Timestamp, nil
		}

		return v.Deneb.Block.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.ExecutionPayloadHeader.Timestamp, nil
		}

		return v.Electra.Block.Body.ExecutionPayload.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// GasLimit returns the gas limit of the proposal.
func (v *VersionedProposal) GasLimit() (uint64, error) {
	if v.Version >= spec.DataVersionBellatrix && !v.payloadPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Body.ExecutionPayloadHeader.GasLimit, nil
		}

		return v.Bellatrix.Body.ExecutionPayload.GasLimit, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.GasLimit, nil
		}

		return v.Capella.Body.ExecutionPayload.GasLimit, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Body.ExecutionPayloadHeader.GasLimit, nil
		}

		return v.Deneb.Block.Body.ExecutionPayload.GasLimit, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Body.ExecutionPayloadHeader.GasLimit, nil
		}

		return v.Electra.Block.Body.ExecutionPayload.GasLimit, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Blobs returns the blobs of the proposal.
func (v *VersionedProposal) Blobs() ([]deneb.Blob, error) {
	if v.Version >= spec.DataVersionDeneb && !v.payloadPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Deneb.Blobs, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Electra.Blobs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// KZGProofs returns the KZG proofs of the proposal.
func (v *VersionedProposal) KZGProofs() ([]deneb.KZGProof, error) {
	if v.Version >= spec.DataVersionDeneb && !v.payloadPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Deneb.KZGProofs, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Electra.KZGProofs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Value returns the value of the proposal, in Wei.
func (v *VersionedProposal) Value() *big.Int {
	value := big.NewInt(0)
	if v.ConsensusValue != nil {
		value = value.Add(value, v.ConsensusValue)
	}
	if v.ExecutionValue != nil {
		value = value.Add(value, v.ExecutionValue)
	}

	return value
}

// String returns a string version of the structure.
func (v *VersionedProposal) String() string {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unknown version"
	}
}

func (v *VersionedProposal) proposalPresent() bool {
	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0 != nil
	case spec.DataVersionAltair:
		return v.Altair != nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded != nil
		}

		return v.Bellatrix != nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil
		}

		return v.Capella != nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded != nil
		}

		return v.Deneb.Block != nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded != nil
		}

		return v.Electra.Block != nil
	}

	return false
}

func (v *VersionedProposal) bodyPresent() bool {
	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0 != nil && v.Phase0.Body != nil
	case spec.DataVersionAltair:
		return v.Altair != nil && v.Altair.Body != nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded != nil && v.BellatrixBlinded.Body != nil
		}

		return v.Bellatrix != nil && v.Bellatrix.Body != nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil && v.CapellaBlinded.Body != nil
		}

		return v.Capella != nil && v.Capella.Body != nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded != nil && v.DenebBlinded.Body != nil
		}

		return v.Deneb != nil && v.Deneb.Block != nil && v.Deneb.Block.Body != nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded != nil && v.ElectraBlinded.Body != nil
		}

		return v.Electra != nil && v.Electra.Block != nil && v.Electra.Block.Body != nil
	}

	return false
}

func (v *VersionedProposal) payloadPresent() bool {
	switch v.Version {
	case spec.DataVersionPhase0:
		return false
	case spec.DataVersionAltair:
		return false
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded != nil &&
				v.BellatrixBlinded.Body != nil &&
				v.BellatrixBlinded.Body.ExecutionPayloadHeader != nil
		}

		return v.Bellatrix != nil && v.Bellatrix.Body != nil && v.Bellatrix.Body.ExecutionPayload != nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil && v.CapellaBlinded.Body != nil && v.CapellaBlinded.Body.ExecutionPayloadHeader != nil
		}

		return v.Capella != nil && v.Capella.Body != nil && v.Capella.Body.ExecutionPayload != nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded != nil && v.DenebBlinded.Body != nil && v.DenebBlinded.Body.ExecutionPayloadHeader != nil
		}

		return v.Deneb != nil && v.Deneb.Block != nil && v.Deneb.Block.Body != nil && v.Deneb.Block.Body.ExecutionPayload != nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded != nil && v.ElectraBlinded.Body != nil && v.ElectraBlinded.Body.ExecutionPayloadHeader != nil
		}

		return v.Electra != nil &&
			v.Electra.Block != nil &&
			v.Electra.Block.Body != nil &&
			v.Electra.Block.Body.ExecutionPayload != nil
	}

	return false
}
