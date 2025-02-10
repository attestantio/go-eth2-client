// Copyright © 2023, 2024 Attestant Limited.
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
	"errors"
	"math/big"

	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"

	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedProposal contains a versioned signed beacon node proposal.
type VersionedSignedProposal struct {
	Version          spec.DataVersion
	Blinded          bool
	ConsensusValue   *big.Int
	ExecutionValue   *big.Int
	Phase0           *phase0.SignedBeaconBlock
	Altair           *altair.SignedBeaconBlock
	Bellatrix        *bellatrix.SignedBeaconBlock
	BellatrixBlinded *apiv1bellatrix.SignedBlindedBeaconBlock
	Capella          *capella.SignedBeaconBlock
	CapellaBlinded   *apiv1capella.SignedBlindedBeaconBlock
	Deneb            *apiv1deneb.SignedBlockContents
	DenebBlinded     *apiv1deneb.SignedBlindedBeaconBlock
	Electra          *apiv1electra.SignedBlockContents
	ElectraBlinded   *apiv1electra.SignedBlindedBeaconBlock
}

// AssertPresent throws an error if the expected proposal
// given the version and blinded fields is not present.
func (v *VersionedSignedProposal) AssertPresent() error {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return errors.New("phase0 proposal not present")
		}
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return errors.New("altair proposal not present")
		}
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil && !v.Blinded {
			return errors.New("bellatrix proposal not present")
		}
		if v.BellatrixBlinded == nil && v.Blinded {
			return errors.New("blinded bellatrix proposal not present")
		}
	case spec.DataVersionCapella:
		if v.Capella == nil && !v.Blinded {
			return errors.New("capella proposal not present")
		}
		if v.CapellaBlinded == nil && v.Blinded {
			return errors.New("blinded capella proposal not present")
		}
	case spec.DataVersionDeneb:
		if v.Deneb == nil && !v.Blinded {
			return errors.New("deneb proposal not present")
		}
		if v.DenebBlinded == nil && v.Blinded {
			return errors.New("blinded deneb proposal not present")
		}
	case spec.DataVersionElectra:
		if v.Electra == nil && !v.Blinded {
			return errors.New("electra proposal not present")
		}
		if v.ElectraBlinded == nil && v.Blinded {
			return errors.New("blinded electra proposal not present")
		}
	default:
		return errors.New("unsupported version")
	}

	return nil
}

// Slot returns the slot of the signed proposal.
func (v *VersionedSignedProposal) Slot() (phase0.Slot, error) {
	if err := v.assertMessagePresent(); err != nil {
		return 0, err
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Message.Slot, nil
	case spec.DataVersionAltair:
		return v.Altair.Message.Slot, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Message.Slot, nil
		}

		return v.Bellatrix.Message.Slot, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.Slot, nil
		}

		return v.Capella.Message.Slot, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Message.Slot, nil
		}

		return v.Deneb.SignedBlock.Message.Slot, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Message.Slot, nil
		}

		return v.Electra.SignedBlock.Message.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the signed proposal.
func (v *VersionedSignedProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	if err := v.assertMessagePresent(); err != nil {
		return 0, err
	}

	switch v.Version {
	case spec.DataVersionPhase0:
		return v.Phase0.Message.ProposerIndex, nil
	case spec.DataVersionAltair:
		return v.Altair.Message.ProposerIndex, nil
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Message.ProposerIndex, nil
		}

		return v.Bellatrix.Message.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.ProposerIndex, nil
		}

		return v.Capella.Message.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Message.ProposerIndex, nil
		}

		return v.Deneb.SignedBlock.Message.ProposerIndex, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Message.ProposerIndex, nil
		}

		return v.Electra.SignedBlock.Message.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ExecutionBlockHash returns the hash of the execution payload.
func (v *VersionedSignedProposal) ExecutionBlockHash() (phase0.Hash32, error) {
	if err := v.assertExecutionPayloadPresent(); err != nil {
		return phase0.Hash32{}, err
	}

	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			return v.BellatrixBlinded.Message.Body.ExecutionPayloadHeader.BlockHash, nil
		}

		return v.Bellatrix.Message.Body.ExecutionPayload.BlockHash, nil
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.Body.ExecutionPayloadHeader.BlockHash, nil
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockHash, nil
	case spec.DataVersionDeneb:
		if v.Blinded {
			return v.DenebBlinded.Message.Body.ExecutionPayloadHeader.BlockHash, nil
		}

		return v.Deneb.SignedBlock.Message.Body.ExecutionPayload.BlockHash, nil
	case spec.DataVersionElectra:
		if v.Blinded {
			return v.ElectraBlinded.Message.Body.ExecutionPayloadHeader.BlockHash, nil
		}

		return v.Electra.SignedBlock.Message.Body.ExecutionPayload.BlockHash, nil
	default:
		return phase0.Hash32{}, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedProposal) String() string {
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
		if v.Blinded {
			if v.BellatrixBlinded == nil {
				return ""
			}

			return v.BellatrixBlinded.String()
		}

		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil {
				return ""
			}

			return v.CapellaBlinded.String()
		}

		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case spec.DataVersionDeneb:
		if v.Blinded {
			if v.DenebBlinded == nil {
				return ""
			}

			return v.DenebBlinded.String()
		}

		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case spec.DataVersionElectra:
		if v.Blinded {
			if v.ElectraBlinded == nil {
				return ""
			}

			return v.ElectraBlinded.String()
		}

		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unsupported version"
	}
}

// assertMessagePresent throws an error if the expected message
// given the version and blinded fields is not present.
func (v *VersionedSignedProposal) assertMessagePresent() error {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			if v.BellatrixBlinded == nil ||
				v.BellatrixBlinded.Message == nil {
				return ErrDataMissing
			}
		} else {
			if v.Bellatrix == nil ||
				v.Bellatrix.Message == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil ||
				v.CapellaBlinded.Message == nil {
				return ErrDataMissing
			}
		} else {
			if v.Capella == nil ||
				v.Capella.Message == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionDeneb:
		if v.Blinded {
			if v.DenebBlinded == nil ||
				v.DenebBlinded.Message == nil {
				return ErrDataMissing
			}
		} else {
			if v.Deneb == nil ||
				v.Deneb.SignedBlock == nil ||
				v.Deneb.SignedBlock.Message == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionElectra:
		if v.Blinded {
			if v.ElectraBlinded == nil ||
				v.ElectraBlinded.Message == nil {
				return ErrDataMissing
			}
		} else {
			if v.Electra == nil ||
				v.Electra.SignedBlock == nil ||
				v.Electra.SignedBlock.Message == nil {
				return ErrDataMissing
			}
		}
	default:
		return ErrUnsupportedVersion
	}

	return nil
}

// assertExecutionPayloadPresent throws an error if the expected execution payload or payload header
// given the version and blinded fields is not present.
//
//nolint:gocyclo
func (v *VersionedSignedProposal) assertExecutionPayloadPresent() error {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Blinded {
			if v.BellatrixBlinded == nil ||
				v.BellatrixBlinded.Message == nil ||
				v.BellatrixBlinded.Message.Body == nil ||
				v.BellatrixBlinded.Message.Body.ExecutionPayloadHeader == nil {
				return ErrDataMissing
			}
		} else {
			if v.Bellatrix == nil ||
				v.Bellatrix.Message == nil ||
				v.Bellatrix.Message.Body == nil ||
				v.Bellatrix.Message.Body.ExecutionPayload == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil ||
				v.CapellaBlinded.Message == nil ||
				v.CapellaBlinded.Message.Body == nil ||
				v.CapellaBlinded.Message.Body.ExecutionPayloadHeader == nil {
				return ErrDataMissing
			}
		} else {
			if v.Capella == nil ||
				v.Capella.Message == nil ||
				v.Capella.Message.Body == nil ||
				v.Capella.Message.Body.ExecutionPayload == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionDeneb:
		if v.Blinded {
			if v.DenebBlinded == nil ||
				v.DenebBlinded.Message == nil ||
				v.DenebBlinded.Message.Body == nil ||
				v.DenebBlinded.Message.Body.ExecutionPayloadHeader == nil {
				return ErrDataMissing
			}
		} else {
			if v.Deneb == nil ||
				v.Deneb.SignedBlock == nil ||
				v.Deneb.SignedBlock.Message == nil ||
				v.Deneb.SignedBlock.Message.Body == nil ||
				v.Deneb.SignedBlock.Message.Body.ExecutionPayload == nil {
				return ErrDataMissing
			}
		}
	case spec.DataVersionElectra:
		if v.Blinded {
			if v.ElectraBlinded == nil ||
				v.ElectraBlinded.Message == nil ||
				v.ElectraBlinded.Message.Body == nil ||
				v.ElectraBlinded.Message.Body.ExecutionPayloadHeader == nil {
				return ErrDataMissing
			}
		} else {
			if v.Electra == nil ||
				v.Electra.SignedBlock == nil ||
				v.Electra.SignedBlock.Message == nil ||
				v.Electra.SignedBlock.Message.Body == nil ||
				v.Electra.SignedBlock.Message.Body.ExecutionPayload == nil {
				return ErrDataMissing
			}
		}
	default:
		return ErrUnsupportedVersion
	}

	return nil
}
