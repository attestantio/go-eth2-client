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

// VersionedV3Proposal contains a versioned generic beacon block.
type VersionedV3Proposal struct {
	Version                 spec.DataVersion
	ExecutionPayloadBlinded bool
	ExecutionPayloadValue   string
	BlindedBellatrix        *apiv1bellatrix.BlindedBeaconBlock
	BlindedCapella          *apiv1capella.BlindedBeaconBlock
	BlindedDeneb            *apiv1deneb.BlindedBeaconBlock
	Phase0                  *phase0.BeaconBlock
	Altair                  *altair.BeaconBlock
	Bellatrix               *bellatrix.BeaconBlock
	Capella                 *capella.BeaconBlock
	Deneb                   *apiv1deneb.BlockContents
}

// IsEmpty returns true if there is no block.
func (v *VersionedV3Proposal) IsEmpty() bool {
	if v.ExecutionPayloadBlinded {
		return v.BlindedBellatrix == nil && v.BlindedCapella == nil && v.BlindedDeneb == nil
	}
	return v.Phase0 == nil &&
		v.Altair == nil &&
		v.Bellatrix == nil &&
		v.Capella == nil &&
		v.Deneb == nil

}

// Slot returns the slot of the generic beacon block.
func (v *VersionedV3Proposal) Slot() (phase0.Slot, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return 0, ErrDataMissing
			}
			return v.BlindedBellatrix.Slot, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return 0, ErrDataMissing
			}
			return v.BlindedCapella.Slot, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return 0, ErrDataMissing
			}
			return v.BlindedDeneb.Slot, nil
		default:
			return 0, ErrUnsupportedVersion
		}
	} else {
		switch v.Version {
		case spec.DataVersionPhase0:
			if v.Phase0 == nil {
				return 0, ErrDataMissing
			}
			return v.Phase0.Slot, nil
		case spec.DataVersionAltair:
			if v.Altair == nil {
				return 0, ErrDataMissing
			}
			return v.Altair.Slot, nil
		case spec.DataVersionBellatrix:
			if v.Bellatrix == nil {
				return 0, ErrDataMissing
			}
			return v.Bellatrix.Slot, nil
		case spec.DataVersionCapella:
			if v.Capella == nil {
				return 0, ErrDataMissing
			}
			return v.Capella.Slot, nil
		case spec.DataVersionDeneb:
			if v.Deneb == nil {
				return 0, ErrDataMissing
			}
			return v.Deneb.Block.Slot, nil
		default:
			return 0, ErrUnsupportedVersion
		}
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedV3Proposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedBellatrix.ProposerIndex, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedCapella.ProposerIndex, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedDeneb.ProposerIndex, nil
		default:
			return 0, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, ErrDataMissing
		}
		return v.Phase0.ProposerIndex, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return 0, ErrDataMissing
		}
		return v.Altair.ProposerIndex, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, ErrDataMissing
		}
		return v.Bellatrix.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, ErrDataMissing
		}
		return v.Capella.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return 0, ErrDataMissing
		}
		return v.Deneb.Block.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}

}

// RandaoReveal returns the RANDAO reveal of the generic beacon block.
func (v *VersionedV3Proposal) RandaoReveal() (phase0.BLSSignature, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil ||
				v.BlindedBellatrix.Body == nil {
				return phase0.BLSSignature{}, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.RANDAOReveal, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil ||
				v.BlindedCapella.Body == nil {
				return phase0.BLSSignature{}, ErrDataMissing
			}

			return v.BlindedCapella.Body.RANDAOReveal, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil ||
				v.BlindedDeneb.Body == nil {
				return phase0.BLSSignature{}, ErrDataMissing
			}

			return v.BlindedDeneb.Body.RANDAOReveal, nil
		default:
			return phase0.BLSSignature{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}
		return v.Phase0.Body.RANDAOReveal, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}
		return v.Altair.Body.RANDAOReveal, nil

	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}
		return v.Bellatrix.Body.RANDAOReveal, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}
		return v.Capella.Body.RANDAOReveal, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}
		return v.Deneb.Block.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the generic beacon block.
func (v *VersionedV3Proposal) Graffiti() ([32]byte, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil ||
				v.BlindedBellatrix.Body == nil {
				return [32]byte{}, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.Graffiti, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil ||
				v.BlindedCapella.Body == nil {
				return [32]byte{}, ErrDataMissing
			}

			return v.BlindedCapella.Body.Graffiti, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil ||
				v.BlindedDeneb.Body == nil {
				return [32]byte{}, ErrDataMissing
			}

			return v.BlindedDeneb.Body.Graffiti, nil
		default:
			return [32]byte{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return [32]byte{}, ErrDataMissing
		}
		return v.Phase0.Body.Graffiti, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return [32]byte{}, ErrDataMissing
		}
		return v.Altair.Body.Graffiti, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, ErrDataMissing
		}
		return v.Bellatrix.Body.Graffiti, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, ErrDataMissing
		}
		return v.Capella.Body.Graffiti, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, ErrDataMissing
		}
		return v.Deneb.Block.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the generic beacon block.
func (v *VersionedV3Proposal) Attestations() ([]*phase0.Attestation, error) {
	if v.ExecutionPayloadBlinded {

		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil ||
				v.BlindedBellatrix.Body == nil {
				return nil, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.Attestations, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil ||
				v.BlindedCapella.Body == nil {
				return nil, ErrDataMissing
			}

			return v.BlindedCapella.Body.Attestations, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil ||
				v.BlindedDeneb.Body == nil {
				return nil, ErrDataMissing
			}

			return v.BlindedDeneb.Body.Attestations, nil
		default:
			return nil, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil ||
			v.Phase0.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Phase0.Body.Attestations, nil
	case spec.DataVersionAltair:
		if v.Altair == nil ||
			v.Altair.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Altair.Body.Attestations, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Bellatrix.Body.Attestations, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.Attestations, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Block.Body.Attestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Root returns the root of the generic beacon block.
func (v *VersionedV3Proposal) Root() (phase0.Root, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedBellatrix.HashTreeRoot()
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedCapella.HashTreeRoot()
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedDeneb.HashTreeRoot()
		default:
			return phase0.Root{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.HashTreeRoot()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.HashTreeRoot()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the generic beacon block.
func (v *VersionedV3Proposal) BodyRoot() (phase0.Root, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.HashTreeRoot()
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedCapella.Body.HashTreeRoot()
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedDeneb.Body.HashTreeRoot()
		default:
			return phase0.Root{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.Body.HashTreeRoot()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.Body.HashTreeRoot()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the generic beacon block.
func (v *VersionedV3Proposal) ParentRoot() (phase0.Root, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedBellatrix.ParentRoot, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedCapella.ParentRoot, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedDeneb.ParentRoot, nil
		default:
			return phase0.Root{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.ParentRoot, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.ParentRoot, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.ParentRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the generic beacon block.
func (v *VersionedV3Proposal) StateRoot() (phase0.Root, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedBellatrix.StateRoot, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedCapella.StateRoot, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil {
				return phase0.Root{}, ErrDataMissing
			}

			return v.BlindedDeneb.StateRoot, nil
		default:
			return phase0.Root{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.StateRoot, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.StateRoot, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.StateRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// Transactions returns the transactions of the proposal.
func (v *VersionedV3Proposal) Transactions() ([]bellatrix.Transaction, error) {
	if v.ExecutionPayloadBlinded {
		return nil, ErrDataMissing
	}
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// TransactionsRoot returns the transactions root of the generic beacon block.
func (v *VersionedV3Proposal) TransactionsRoot() (phase0.Root, error) {
	if !v.ExecutionPayloadBlinded {
		return phase0.Root{}, ErrDataMissing
	}
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.BlindedBellatrix == nil ||
			v.BlindedBellatrix.Body == nil ||
			v.BlindedBellatrix.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.BlindedBellatrix.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionCapella:
		if v.BlindedCapella == nil ||
			v.BlindedCapella.Body == nil ||
			v.BlindedCapella.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.BlindedCapella.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionDeneb:
		if v.BlindedDeneb == nil ||
			v.BlindedDeneb.Body == nil ||
			v.BlindedDeneb.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.BlindedDeneb.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the generic beacon block.
func (v *VersionedV3Proposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil ||
				v.BlindedBellatrix.Body == nil ||
				v.BlindedBellatrix.Body.ExecutionPayloadHeader == nil {
				return bellatrix.ExecutionAddress{}, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.ExecutionPayloadHeader.FeeRecipient, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil ||
				v.BlindedCapella.Body == nil ||
				v.BlindedCapella.Body.ExecutionPayloadHeader == nil {
				return bellatrix.ExecutionAddress{}, ErrDataMissing
			}

			return v.BlindedCapella.Body.ExecutionPayloadHeader.FeeRecipient, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil ||
				v.BlindedDeneb.Body == nil ||
				v.BlindedDeneb.Body.ExecutionPayloadHeader == nil {
				return bellatrix.ExecutionAddress{}, ErrDataMissing
			}

			return v.BlindedDeneb.Body.ExecutionPayloadHeader.FeeRecipient, nil
		default:
			return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the generic beacon block.
func (v *VersionedV3Proposal) Timestamp() (uint64, error) {
	if v.ExecutionPayloadBlinded {
		switch v.Version {
		case spec.DataVersionBellatrix:
			if v.BlindedBellatrix == nil ||
				v.BlindedBellatrix.Body == nil ||
				v.BlindedBellatrix.Body.ExecutionPayloadHeader == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedBellatrix.Body.ExecutionPayloadHeader.Timestamp, nil
		case spec.DataVersionCapella:
			if v.BlindedCapella == nil ||
				v.BlindedCapella.Body == nil ||
				v.BlindedCapella.Body.ExecutionPayloadHeader == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedCapella.Body.ExecutionPayloadHeader.Timestamp, nil
		case spec.DataVersionDeneb:
			if v.BlindedDeneb == nil ||
				v.BlindedDeneb.Body == nil ||
				v.BlindedDeneb.Body.ExecutionPayloadHeader == nil {
				return 0, ErrDataMissing
			}

			return v.BlindedDeneb.Body.ExecutionPayloadHeader.Timestamp, nil
		default:
			return 0, ErrUnsupportedVersion
		}
	}
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Blobs returns the blobs of the proposal.
func (v *VersionedV3Proposal) Blobs() ([]deneb.Blob, error) {
	if v.ExecutionPayloadBlinded {
		return nil, ErrDataMissing
	}
	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Blobs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// KZGProofs returns the KZG proofs of the proposal.
func (v *VersionedV3Proposal) KZGProofs() ([]deneb.KZGProof, error) {
	if v.ExecutionPayloadBlinded {
		return nil, ErrDataMissing
	}
	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.KZGProofs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedV3Proposal) String() string {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.BlindedBellatrix == nil {
			return ""
		}

		return v.BlindedBellatrix.String()
	case spec.DataVersionCapella:
		if v.BlindedCapella == nil {
			return ""
		}

		return v.BlindedCapella.String()
	case spec.DataVersionDeneb:
		if v.BlindedDeneb == nil {
			return ""
		}

		return v.BlindedDeneb.String()
	default:
		return "unknown version"
	}
}
