// Copyright © 2021 - 2023 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBeaconState contains a versioned beacon state.
type VersionedBeaconState struct {
	Version   DataVersion
	Phase0    *phase0.BeaconState
	Altair    *altair.BeaconState
	Bellatrix *bellatrix.BeaconState
	Capella   *capella.BeaconState
	Deneb     *deneb.BeaconState
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconState) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil
}

// Slot returns the slot of the state.
func (v *VersionedBeaconState) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no Phase0 state")
		}

		return v.Phase0.Slot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no Altair state")
		}

		return v.Altair.Slot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no Bellatrix state")
		}

		return v.Bellatrix.Slot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return v.Capella.Slot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no Deneb state")
		}

		return v.Deneb.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// NextWithdrawalValidatorIndex returns the next withdrawal validator index of the state.
func (v *VersionedBeaconState) NextWithdrawalValidatorIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix:
		return 0, errors.New("state does not provide next withdrawal validator index")
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return v.Capella.NextWithdrawalValidatorIndex, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no Deneb state")
		}

		return v.Deneb.NextWithdrawalValidatorIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Validators returns the validators of the state.
func (v *VersionedBeaconState) Validators() ([]*phase0.Validator, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 state")
		}

		return v.Phase0.Validators, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair state")
		}

		return v.Altair.Validators, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix state")
		}

		return v.Bellatrix.Validators, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.Validators, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb state")
		}

		return v.Deneb.Validators, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ValidatorBalances returns the validator balances of the state.
func (v *VersionedBeaconState) ValidatorBalances() ([]phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 state")
		}

		return v.Phase0.Balances, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair state")
		}

		return v.Altair.Balances, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix state")
		}

		return v.Bellatrix.Balances, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.Balances, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb state")
		}

		return v.Deneb.Balances, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedBeaconState) String() string {
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
	default:
		return "unknown version"
	}
}
