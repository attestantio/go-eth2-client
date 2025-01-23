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

package spec

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedAggregateAndProof contains a versioned aggregate and proof.
type VersionedAggregateAndProof struct {
	Version   DataVersion
	Phase0    *phase0.AggregateAndProof
	Altair    *phase0.AggregateAndProof
	Bellatrix *phase0.AggregateAndProof
	Capella   *phase0.AggregateAndProof
	Deneb     *phase0.AggregateAndProof
	Electra   *electra.AggregateAndProof
}

// AggregatorIndex returns the aggregator index of the aggregate.
func (v *VersionedAggregateAndProof) AggregatorIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 aggregate and proof")
		}

		return v.Phase0.AggregatorIndex, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair aggregate and proof")
		}

		return v.Altair.AggregatorIndex, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix aggregate and proof")
		}

		return v.Bellatrix.AggregatorIndex, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella aggregate and proof")
		}

		return v.Capella.AggregatorIndex, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb aggregate and proof")
		}

		return v.Deneb.AggregatorIndex, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra aggregate and proof")
		}

		return v.Electra.AggregatorIndex, nil
	default:
		return 0, errors.New("unknown version for aggregate and proof")
	}
}

// HashTreeRoot returns the hash tree root of the aggregate and proof.
func (v *VersionedAggregateAndProof) HashTreeRoot() ([32]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return [32]byte{}, errors.New("no phase0 aggregate and proof")
		}

		return v.Phase0.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return [32]byte{}, errors.New("no altair aggregate and proof")
		}

		return v.Altair.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, errors.New("no bellatrix aggregate and proof")
		}

		return v.Bellatrix.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella aggregate and proof")
		}

		return v.Capella.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, errors.New("no deneb aggregate and proof")
		}

		return v.Deneb.HashTreeRoot()
	case DataVersionElectra:
		if v.Electra == nil {
			return [32]byte{}, errors.New("no electra aggregate and proof")
		}

		return v.Electra.HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// IsEmpty returns true if there is no aggregate and proof.
func (v *VersionedAggregateAndProof) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
}

// String returns a string version of the structure.
func (v *VersionedAggregateAndProof) String() string {
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

// SelectionProof returns the selection proof of the aggregate.
func (v *VersionedAggregateAndProof) SelectionProof() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 aggregate and proof")
		}

		return v.Phase0.SelectionProof, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no altair aggregate and proof")
		}

		return v.Altair.SelectionProof, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix aggregate and proof")
		}

		return v.Bellatrix.SelectionProof, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella aggregate and proof")
		}

		return v.Capella.SelectionProof, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no deneb aggregate and proof")
		}

		return v.Deneb.SelectionProof, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no electra aggregate and proof")
		}

		return v.Electra.SelectionProof, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}
