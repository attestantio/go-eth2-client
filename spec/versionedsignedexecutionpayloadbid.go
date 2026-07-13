// Copyright © 2021 - 2024 Attestant Limited.
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

package spec

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/heze"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedExecutionPayloadBid contains a versioned signed execution payload bid.
type VersionedSignedExecutionPayloadBid struct {
	Version DataVersion
	Gloas   *gloas.SignedExecutionPayloadBid
	Heze    *heze.SignedExecutionPayloadBid
}

// String returns a string version of the structure.
func (v *VersionedSignedExecutionPayloadBid) String() string {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella,
		DataVersionDeneb, DataVersionElectra, DataVersionFulu:
		return ""
	case DataVersionGloas:
		if v.Gloas == nil {
			return ""
		}

		return v.Gloas.String()
	case DataVersionHeze:
		if v.Heze == nil {
			return ""
		}

		return v.Heze.String()
	default:
		return "unknown version"
	}
}

// ParentBlockHash returns the parent block hash of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) ParentBlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Hash32{}, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return phase0.Hash32{}, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return phase0.Hash32{}, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return phase0.Hash32{}, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return phase0.Hash32{}, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return phase0.Hash32{}, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return phase0.Hash32{}, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return phase0.Hash32{}, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.ParentBlockHash, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return phase0.Hash32{}, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.ParentBlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// ParentBlockRoot returns the parent block root of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) ParentBlockRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Root{}, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return phase0.Root{}, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return phase0.Root{}, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return phase0.Root{}, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return phase0.Root{}, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return phase0.Root{}, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return phase0.Root{}, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return phase0.Root{}, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.ParentBlockRoot, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return phase0.Root{}, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.ParentBlockRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// BlockHash returns the block hash of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) BlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Hash32{}, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return phase0.Hash32{}, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return phase0.Hash32{}, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return phase0.Hash32{}, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return phase0.Hash32{}, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return phase0.Hash32{}, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return phase0.Hash32{}, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return phase0.Hash32{}, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.BlockHash, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return phase0.Hash32{}, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.BlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// FeeRecipient returns the fee recipient of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case DataVersionPhase0:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.FeeRecipient, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, errors.New("unknown version")
	}
}

// GasLimit returns the gas limit of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) GasLimit() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.GasLimit, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return 0, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.GasLimit, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BuilderIndex returns the builder index of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) BuilderIndex() (gloas.BuilderIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.BuilderIndex, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return 0, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.BuilderIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Slot returns the slot of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.Slot, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return 0, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Value returns the value of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) Value() (phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.Value, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return 0, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.Value, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionPayment returns the execution payment of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) ExecutionPayment() (phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.ExecutionPayment, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return 0, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.ExecutionPayment, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BlobKZGCommitments returns the blob kzg commitments of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) BlobKZGCommitments() ([]deneb.KZGCommitment, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return nil, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return nil, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return nil, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return nil, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return nil, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return nil, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Message.BlobKZGCommitments, nil
	case DataVersionHeze:
		if v.Heze == nil || v.Heze.Message == nil {
			return nil, errors.New("no heze execution payload bid")
		}

		return v.Heze.Message.BlobKZGCommitments, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Signature returns the signature of the execution payload bid.
func (v *VersionedSignedExecutionPayloadBid) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in phase0")
	case DataVersionAltair:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in altair")
	case DataVersionBellatrix:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in bellatrix")
	case DataVersionCapella:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in capella")
	case DataVersionDeneb:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in deneb")
	case DataVersionElectra:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in electra")
	case DataVersionFulu:
		return phase0.BLSSignature{}, errors.New("no execution payload bid in fulu")
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.BLSSignature{}, errors.New("no gloas execution payload bid")
		}

		return v.Gloas.Signature, nil
	case DataVersionHeze:
		if v.Heze == nil {
			return phase0.BLSSignature{}, errors.New("no heze execution payload bid")
		}

		return v.Heze.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}
