// Copyright © 2021 - 2025 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/eip7732"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/fulu"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	proofutil "github.com/attestantio/go-eth2-client/util/proof"
	ssz "github.com/ferranbt/fastssz"
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
	Fulu      *fulu.BeaconState
	EIP7732   *eip7732.BeaconState
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconState) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil &&
		v.Electra == nil && v.Fulu == nil && v.EIP7732 == nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.Slot, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.Slot, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.NextWithdrawalValidatorIndex, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.NextWithdrawalValidatorIndex, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.Validators, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return nil, errors.New("no EIP7732 state")
		}

		return v.EIP7732.Validators, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.Balances, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return nil, errors.New("no EIP7732 state")
		}

		return v.EIP7732.Balances, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// DepositRequestsStartIndex returns the deposit requests start index of the state.
func (v *VersionedBeaconState) DepositRequestsStartIndex() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return 0, errors.New("state does not provide deposit requests start index")
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return v.Electra.DepositRequestsStartIndex, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.DepositRequestsStartIndex, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.DepositRequestsStartIndex, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.DepositBalanceToConsume, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.DepositBalanceToConsume, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.ExitBalanceToConsume, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.ExitBalanceToConsume, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.EarliestExitEpoch, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.EarliestExitEpoch, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.ConsolidationBalanceToConsume, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.ConsolidationBalanceToConsume, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return v.Fulu.EarliestConsolidationEpoch, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return 0, errors.New("no EIP7732 state")
		}

		return v.EIP7732.EarliestConsolidationEpoch, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.PendingDeposits, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return nil, errors.New("no EIP7732 state")
		}

		return v.EIP7732.PendingDeposits, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.PendingPartialWithdrawals, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return nil, errors.New("no EIP7732 state")
		}

		return v.EIP7732.PendingPartialWithdrawals, nil
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.PendingConsolidations, nil
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return nil, errors.New("no EIP7732 state")
		}

		return v.EIP7732.PendingConsolidations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ValidatorAtIndex returns the validator at the given index.
// This is a convenience method that handles accessing the validators array.
// Parameters:
//   - index: The index of the validator to retrieve
//
// Returns:
//   - *phase0.Validator: The validator at the given index
//   - error: If the index is invalid or there's an error accessing the validators
func (v *VersionedBeaconState) ValidatorAtIndex(index phase0.ValidatorIndex) (*phase0.Validator, error) {
	validators, err := v.Validators()
	if err != nil {
		return nil, err
	}

	if index >= phase0.ValidatorIndex(len(validators)) {
		return nil, errors.New("validator index out of bounds")
	}

	return validators[index], nil
}

// ValidatorBalance returns the balance of the validator at the given index.
// This is a convenience method that handles accessing the balances array.
// Parameters:
//   - index: The index of the validator whose balance to retrieve
//
// Returns:
//   - phase0.Gwei: The balance in Gwei
//   - error: If the index is invalid or there's an error accessing the balances
func (v *VersionedBeaconState) ValidatorBalance(index phase0.ValidatorIndex) (phase0.Gwei, error) {
	balances, err := v.ValidatorBalances()
	if err != nil {
		return 0, err
	}

	if index >= phase0.ValidatorIndex(len(balances)) {
		return 0, errors.New("validator index out of bounds")
	}

	return balances[index], nil
}

// GetTree returns the GetTree of the specific beacon state version.
func (v *VersionedBeaconState) GetTree() (*ssz.Node, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 state")
		}

		return v.Phase0.GetTree()
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair state")
		}

		return v.Altair.GetTree()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix state")
		}

		return v.Bellatrix.GetTree()
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.GetTree()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb state")
		}

		return v.Deneb.GetTree()
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra state")
		}

		return v.Electra.GetTree()
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no Fulu state")
		}

		return v.Fulu.GetTree()
	default:
		return nil, errors.New("unknown version")
	}
}

// HashTreeRoot returns the HashTreeRoot of the specific beacon state version.
func (v *VersionedBeaconState) HashTreeRoot() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Hash32{}, errors.New("no Phase0 state")
		}

		return v.Phase0.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Hash32{}, errors.New("no Altair state")
		}

		return v.Altair.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Hash32{}, errors.New("no Bellatrix state")
		}

		return v.Bellatrix.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Hash32{}, errors.New("no Capella state")
		}

		return v.Capella.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Hash32{}, errors.New("no Deneb state")
		}

		return v.Deneb.HashTreeRoot()
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.Hash32{}, errors.New("no Electra state")
		}

		return v.Electra.HashTreeRoot()
	case DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Hash32{}, errors.New("no Fulu state")
		}

		return v.Fulu.HashTreeRoot()
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// FieldIndex returns the struct field index for a given field name.
// The index represents the field's position in the struct's memory layout.
// Returns an error if the field doesn't exist or the state is empty.
func (v *VersionedBeaconState) FieldIndex(name string) (int, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no Phase0 state")
		}

		return proofutil.FieldIndex(v.Phase0, name)
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no Altair state")
		}

		return proofutil.FieldIndex(v.Altair, name)
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no Bellatrix state")
		}

		return proofutil.FieldIndex(v.Bellatrix, name)
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return proofutil.FieldIndex(v.Capella, name)
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no Deneb state")
		}

		return proofutil.FieldIndex(v.Deneb, name)
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return proofutil.FieldIndex(v.Electra, name)
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return proofutil.FieldIndex(v.Fulu, name)
	default:
		return 0, errors.New("unknown version")
	}
}

// FieldGeneralizedIndex returns the generalized index for a given field name.
// The generalized index represents the field's absolute position in the Merkle tree.
// This is used for generating and verifying Merkle proofs.
// Returns an error if the field doesn't exist or the state is empty.
func (v *VersionedBeaconState) FieldGeneralizedIndex(name string) (int, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no Phase0 state")
		}

		return proofutil.FieldGeneralizedIndex(v.Phase0, name)
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no Altair state")
		}

		return proofutil.FieldGeneralizedIndex(v.Altair, name)
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no Bellatrix state")
		}

		return proofutil.FieldGeneralizedIndex(v.Bellatrix, name)
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return proofutil.FieldGeneralizedIndex(v.Capella, name)
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no Deneb state")
		}

		return proofutil.FieldGeneralizedIndex(v.Deneb, name)
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no Electra state")
		}

		return proofutil.FieldGeneralizedIndex(v.Electra, name)
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no Fulu state")
		}

		return proofutil.FieldGeneralizedIndex(v.Fulu, name)
	default:
		return 0, errors.New("unknown version")
	}
}

// FieldRoot returns the SSZ hash root of a specific field in the beacon state.
// Parameters:
//   - name: The name of the field to get the root for
//
// Returns:
//   - phase0.Hash32: The SSZ hash root of the field
//   - error: If the field doesn't exist, the state is empty, or the field is not hash tree rootable
func (v *VersionedBeaconState) FieldRoot(name string) (phase0.Hash32, error) {
	fieldTree, err := v.FieldTree(name)
	if err != nil {
		return phase0.Hash32{}, err
	}
	var root phase0.Hash32
	copy(root[:], fieldTree.Hash())

	return root, nil
}

// FieldTree returns the Merkle subtree for a specific field in the beacon state.
func (v *VersionedBeaconState) FieldTree(name string) (*ssz.Node, error) {
	stateTree, err := v.GetTree()
	if err != nil {
		return nil, err
	}

	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return nil, err
	}

	return stateTree.Get(fieldGeneralizedIndex)
}

// ProveField generates a Merkle proof for a specific field against the beacon state root.
// Parameters:
//   - name: The name of the field to generate a proof for
//
// Returns:
//   - []phase0.Hash32: The Merkle proof as a sequence of 32-byte hashes
//   - error: If the field doesn't exist or there's an error generating the proof
func (v *VersionedBeaconState) ProveField(name string) ([]phase0.Hash32, error) {
	stateTree, err := v.GetTree()
	if err != nil {
		return nil, err
	}

	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return nil, err
	}

	proof, err := stateTree.Prove(fieldGeneralizedIndex)
	if err != nil {
		return nil, err
	}

	proofBytes := make([]phase0.Hash32, len(proof.Hashes))
	for i, hash := range proof.Hashes {
		copy(proofBytes[i][:], hash)
	}

	return proofBytes, nil
}

// VerifyFieldProof verifies a Merkle proof for a field against the beacon state root.
// Parameters:
//   - proof: The Merkle proof as a sequence of 32-byte hashes
//   - name: The name of the field the proof is for
//
// Returns:
//   - bool: True if the proof is valid, false otherwise
//   - error: If there's an error during verification
func (v *VersionedBeaconState) VerifyFieldProof(proof []phase0.Hash32, name string) (bool, error) {
	// Get the state root
	stateRoot, err := v.HashTreeRoot()
	if err != nil {
		return false, err
	}

	// Get the field's generalized index
	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return false, err
	}

	// Get the field's root
	fieldRoot, err := v.FieldRoot(name)
	if err != nil {
		return false, err
	}

	// Convert proof to ssz.Proof format
	proofHashes := make([][]byte, len(proof))
	for i, hash := range proof {
		proofHashes[i] = make([]byte, 32)
		copy(proofHashes[i], hash[:])
	}

	// Create and verify the proof
	sszProof := &ssz.Proof{
		Index:  fieldGeneralizedIndex,
		Leaf:   fieldRoot[:],
		Hashes: proofHashes,
	}

	return ssz.VerifyProof(stateRoot[:], sszProof)
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
	case DataVersionFulu:
		if v.Fulu == nil {
			return ""
		}

		return v.Fulu.String()
	case DataVersionEIP7732:
		if v.EIP7732 == nil {
			return ""
		}

		return v.EIP7732.String()
	default:
		return "unknown version"
	}
}
