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

package spec

import (
	"errors"
	"fmt"

	"github.com/prysmaticlabs/go-bitfield"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedAttestation contains a versioned attestation.
type VersionedAttestation struct {
	Version        DataVersion
	ValidatorIndex *phase0.ValidatorIndex
	Phase0         *phase0.Attestation
	Altair         *phase0.Attestation
	Bellatrix      *phase0.Attestation
	Capella        *phase0.Attestation
	Deneb          *phase0.Attestation
	Electra        *electra.Attestation
}

// IsEmpty returns true if there is no block.
func (v *VersionedAttestation) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
}

// AggregationBits returns the aggregation bits of the attestation.
func (v *VersionedAttestation) AggregationBits() (bitfield.Bitlist, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 attestation")
		}

		return v.Phase0.AggregationBits, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair attestation")
		}

		return v.Altair.AggregationBits, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix attestation")
		}

		return v.Bellatrix.AggregationBits, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella attestation")
		}

		return v.Capella.AggregationBits, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb attestation")
		}

		return v.Deneb.AggregationBits, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra attestation")
		}

		return v.Electra.AggregationBits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Data returns the data of the attestation.
func (v *VersionedAttestation) Data() (*phase0.AttestationData, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 attestation")
		}

		return v.Phase0.Data, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair attestation")
		}

		return v.Altair.Data, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix attestation")
		}

		return v.Bellatrix.Data, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella attestation")
		}

		return v.Capella.Data, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb attestation")
		}

		return v.Deneb.Data, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra attestation")
		}

		return v.Electra.Data, nil
	default:
		return nil, fmt.Errorf("unknown version: %d", v.Version)
	}
}

// CommitteeBits returns the committee bits of the attestation.
func (v *VersionedAttestation) CommitteeBits() (bitfield.Bitvector64, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("attestation does not provide committee bits")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra attestation")
		}

		return v.Electra.CommitteeBits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// CommitteeIndex returns the index if only one bit is set, otherwise error.
func (v *VersionedAttestation) CommitteeIndex() (phase0.CommitteeIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no Phase0 attestation")
		}

		return v.Phase0.Data.Index, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no Altair attestation")
		}

		return v.Altair.Data.Index, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no Bellatrix attestation")
		}

		return v.Bellatrix.Data.Index, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella attestation")
		}

		return v.Capella.Data.Index, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no Deneb attestation")
		}

		return v.Deneb.Data.Index, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra attestation")
		}

		return v.Electra.CommitteeIndex()
	default:
		return 0, errors.New("unknown version")
	}
}

func (v *VersionedAttestation) HashTreeRoot() ([32]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return [32]byte{}, errors.New("no Phase0 attestation")
		}

		return v.Phase0.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return [32]byte{}, errors.New("no Altair attestation")
		}

		return v.Altair.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, errors.New("no Bellatrix attestation")
		}

		return v.Bellatrix.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no Capella attestation")
		}

		return v.Capella.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, errors.New("no Deneb attestation")
		}

		return v.Deneb.HashTreeRoot()
	case DataVersionElectra:
		if v.Electra == nil {
			return [32]byte{}, errors.New("no Electra attestation")
		}

		return v.Electra.HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// Signature returns the signature of the attestation.
func (v *VersionedAttestation) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no Phase0 attestation")
		}

		return v.Phase0.Signature, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no Altair attestation")
		}

		return v.Altair.Signature, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no Bellatrix attestation")
		}

		return v.Bellatrix.Signature, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no Capella attestation")
		}

		return v.Capella.Signature, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no Deneb attestation")
		}

		return v.Deneb.Signature, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no Electra attestation")
		}

		return v.Electra.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedAttestation) String() string {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unknown version"
	}
}
