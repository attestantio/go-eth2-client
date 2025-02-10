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

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
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
	Electra   *electra.BeaconState
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconState) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
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
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.Slot, nil
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
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.NextWithdrawalValidatorIndex, nil
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
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.Validators, nil
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
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.Balances, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// DepositReceiptsStartIndex returns the deposit requests start index of the state.
func (v *VersionedBeaconState) DepositRequestsStartIndex() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide deposit requests start index")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.DepositRequestsStartIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// DepositBalanceToConsume returns the deposit balance to consume of the state.
func (v *VersionedBeaconState) DepositBalanceToConsume() (phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide deposit balance to consume")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.DepositBalanceToConsume, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExitBalanceToConsume returns the deposit balance to consume of the state.
func (v *VersionedBeaconState) ExitBalanceToConsume() (phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide exit balance to consume")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.ExitBalanceToConsume, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// EarliestExitEpoch returns the earliest exit epoch of the state.
func (v *VersionedBeaconState) EarliestExitEpoch() (phase0.Epoch, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide earliest exit epoch")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.EarliestExitEpoch, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ConsolidationBalanceToConsume returns the consolidation balance to consume of the state.
func (v *VersionedBeaconState) ConsolidationBalanceToConsume() (phase0.Gwei, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide consolidation balance to consume")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.ConsolidationBalanceToConsume, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// EarliestConsolidationEpoch returns the earliest consolidation epoch of the state.
func (v *VersionedBeaconState) EarliestConsolidationEpoch() (phase0.Epoch, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide earliest consolidation epoch")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.EarliestConsolidationEpoch, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// PendingDeposits returns the pending deposits of the state.
func (v *VersionedBeaconState) PendingDeposits() ([]*electra.PendingDeposit, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("state does not provide pending deposits")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.PendingDeposits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// PendingPartialWithdrawals returns the pending partial withdrawals of the state.
func (v *VersionedBeaconState) PendingPartialWithdrawals() ([]*electra.PendingPartialWithdrawal, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("state does not provide pending partial withdrawals")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.PendingPartialWithdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// PendingConsolidations returns the pending consolidations of the state.
func (v *VersionedBeaconState) PendingConsolidations() ([]*electra.PendingConsolidation, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("state does not provide pending consolidations")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.PendingConsolidations, nil
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
	case DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unknown version"
	}
}
