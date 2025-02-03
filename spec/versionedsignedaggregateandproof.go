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

// VersionedSignedAggregateAndProof contains a versioned signed aggregate and proof.
type VersionedSignedAggregateAndProof struct {
	Version   DataVersion
	Phase0    *phase0.SignedAggregateAndProof
	Altair    *phase0.SignedAggregateAndProof
	Bellatrix *phase0.SignedAggregateAndProof
	Capella   *phase0.SignedAggregateAndProof
	Deneb     *phase0.SignedAggregateAndProof
	Electra   *electra.SignedAggregateAndProof
}

// AggregatorIndex returns the aggregator index of the aggregate.
func (v *VersionedSignedAggregateAndProof) AggregatorIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 signed aggregate and proof")
		}

		return v.Phase0.Message.AggregatorIndex, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair signed aggregate and proof")
		}

		return v.Altair.Message.AggregatorIndex, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix signed aggregate and proof")
		}

		return v.Bellatrix.Message.AggregatorIndex, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.AggregatorIndex, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb signed aggregate and proof")
		}

		return v.Deneb.Message.AggregatorIndex, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra signed aggregate and proof")
		}

		return v.Electra.Message.AggregatorIndex, nil
	default:
		return 0, errors.New("unknown version for signed aggregate and proof")
	}
}

// IsEmpty returns true if there is no aggregate and proof.
func (v *VersionedSignedAggregateAndProof) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
}

// SelectionProof returns the selection proof of the signed aggregate.
func (v *VersionedSignedAggregateAndProof) SelectionProof() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 signed aggregate and proof")
		}

		return v.Phase0.Message.SelectionProof, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no altair signed aggregate and proof")
		}

		return v.Altair.Message.SelectionProof, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix signed aggregate and proof")
		}

		return v.Bellatrix.Message.SelectionProof, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.SelectionProof, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no deneb signed aggregate and proof")
		}

		return v.Deneb.Message.SelectionProof, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no electra signed aggregate and proof")
		}

		return v.Electra.Message.SelectionProof, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// Signature returns the signature of the signed aggregate and proof.
func (v *VersionedSignedAggregateAndProof) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 signed aggregate and proof")
		}

		return v.Phase0.Signature, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no altair signed aggregate and proof")
		}

		return v.Altair.Signature, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix signed aggregate and proof")
		}

		return v.Bellatrix.Signature, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Signature, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no deneb signed aggregate and proof")
		}

		return v.Deneb.Signature, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no electra signed aggregate and proof")
		}

		return v.Electra.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// Slot returns the slot of the signed aggregate and proof.
func (v *VersionedSignedAggregateAndProof) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 signed aggregate and proof")
		}

		return v.Phase0.Message.Aggregate.Data.Slot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair signed aggregate and proof")
		}

		return v.Altair.Message.Aggregate.Data.Slot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix signed aggregate and proof")
		}

		return v.Bellatrix.Message.Aggregate.Data.Slot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.Aggregate.Data.Slot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb signed aggregate and proof")
		}

		return v.Deneb.Message.Aggregate.Data.Slot, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra signed aggregate and proof")
		}

		return v.Electra.Message.Aggregate.Data.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedAggregateAndProof) String() string {
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
