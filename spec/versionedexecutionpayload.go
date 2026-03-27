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

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
)

// VersionedExecutionPayload contains a versioned execution payload.
type VersionedExecutionPayload struct {
	Version   DataVersion
	Bellatrix *bellatrix.ExecutionPayload
	Capella   *capella.ExecutionPayload
	Deneb     *deneb.ExecutionPayload
	Electra   *deneb.ExecutionPayload
	Fulu      *deneb.ExecutionPayload
	Gloas     *gloas.ExecutionPayload
}

// IsEmpty returns true if there is no block.
func (v *VersionedExecutionPayload) IsEmpty() bool {
	return v.Version < DataVersionBellatrix || (v.Bellatrix == nil &&
		v.Capella == nil && v.Deneb == nil && v.Electra == nil && v.Fulu == nil && v.Gloas == nil)
}

// ParentHash returns the parent hash of the execution payload.
func (v *VersionedExecutionPayload) ParentHash() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Hash32{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return phase0.Hash32{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Hash32{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.ParentHash, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Hash32{}, errors.New("no capella execution payload")
		}

		return v.Capella.ParentHash, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Hash32{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.ParentHash, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.Hash32{}, errors.New("no electra execution payload")
		}

		return v.Electra.ParentHash, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Hash32{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.ParentHash, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Hash32{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.ParentHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// FeeRecipient returns the fee recipient of the execution payload.
func (v *VersionedExecutionPayload) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case DataVersionPhase0:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return bellatrix.ExecutionAddress{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.FeeRecipient, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no capella execution payload")
		}

		return v.Capella.FeeRecipient, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.FeeRecipient, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no electra execution payload")
		}

		return v.Electra.FeeRecipient, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.FeeRecipient, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the execution payload.
func (v *VersionedExecutionPayload) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Root{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return phase0.Root{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.StateRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella execution payload")
		}

		return v.Capella.StateRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.StateRoot, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, errors.New("no electra execution payload")
		}

		return v.Electra.StateRoot, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.StateRoot, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// ReceiptsRoot returns the receipts root of the execution payload.
func (v *VersionedExecutionPayload) ReceiptsRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Root{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return phase0.Root{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.ReceiptsRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella execution payload")
		}

		return v.Capella.ReceiptsRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.ReceiptsRoot, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, errors.New("no electra execution payload")
		}

		return v.Electra.ReceiptsRoot, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.ReceiptsRoot, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.ReceiptsRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// LogsBloom returns the logs bloom of the execution payload.
func (v *VersionedExecutionPayload) LogsBloom() ([256]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		return [256]byte{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return [256]byte{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [256]byte{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.LogsBloom, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return [256]byte{}, errors.New("no capella execution payload")
		}

		return v.Capella.LogsBloom, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return [256]byte{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.LogsBloom, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return [256]byte{}, errors.New("no electra execution payload")
		}

		return v.Electra.LogsBloom, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return [256]byte{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.LogsBloom, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return [256]byte{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.LogsBloom, nil
	default:
		return [256]byte{}, errors.New("unknown version")
	}
}

// PrevRandao returns the prev randao of the execution payload.
func (v *VersionedExecutionPayload) PrevRandao() ([32]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		return [32]byte{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return [32]byte{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.PrevRandao, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella execution payload")
		}

		return v.Capella.PrevRandao, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.PrevRandao, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return [32]byte{}, errors.New("no electra execution payload")
		}

		return v.Electra.PrevRandao, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return [32]byte{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.PrevRandao, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return [32]byte{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.PrevRandao, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// BlockNumber returns the block number of the execution payload.
func (v *VersionedExecutionPayload) BlockNumber() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.BlockNumber, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.BlockNumber, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.BlockNumber, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.BlockNumber, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.BlockNumber, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// GasLimit returns the gas limit of the execution payload.
func (v *VersionedExecutionPayload) GasLimit() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.GasLimit, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.GasLimit, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.GasLimit, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.GasLimit, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.GasLimit, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.GasLimit, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// GasUsed returns the gas used of the execution payload.
func (v *VersionedExecutionPayload) GasUsed() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.GasUsed, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.GasUsed, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.GasUsed, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.GasUsed, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.GasUsed, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.GasUsed, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Timestamp returns the timestamp of the execution payload.
func (v *VersionedExecutionPayload) Timestamp() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.Timestamp, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.Timestamp, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.Timestamp, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.Timestamp, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.Timestamp, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.Timestamp, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExtraData returns the extra data of the execution payload.
func (v *VersionedExecutionPayload) ExtraData() ([]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.ExtraData, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.ExtraData, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb execution payload")
		}

		return v.Deneb.ExtraData, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution payload")
		}

		return v.Electra.ExtraData, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution payload")
		}

		return v.Fulu.ExtraData, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution payload")
		}

		return v.Gloas.ExtraData, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BaseFeePerGas returns the base fee per gas of the execution payload.
func (v *VersionedExecutionPayload) BaseFeePerGas() (*uint256.Int, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix execution payload")
		}

		return uint256.NewInt(0).SetBytes(v.Bellatrix.BaseFeePerGas[:]), nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return uint256.NewInt(0).SetBytes(v.Capella.BaseFeePerGas[:]), nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb execution payload")
		}

		return v.Deneb.BaseFeePerGas, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution payload")
		}

		return v.Electra.BaseFeePerGas, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution payload")
		}

		return v.Fulu.BaseFeePerGas, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution payload")
		}

		return v.Gloas.BaseFeePerGas, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BlockHash returns the block hash of the execution payload.
func (v *VersionedExecutionPayload) BlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Hash32{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return phase0.Hash32{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Hash32{}, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.BlockHash, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Hash32{}, errors.New("no capella execution payload")
		}

		return v.Capella.BlockHash, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Hash32{}, errors.New("no deneb execution payload")
		}

		return v.Deneb.BlockHash, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.Hash32{}, errors.New("no electra execution payload")
		}

		return v.Electra.BlockHash, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Hash32{}, errors.New("no fulu execution payload")
		}

		return v.Fulu.BlockHash, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Hash32{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.BlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// Transactions returns the transactions of the execution payload.
func (v *VersionedExecutionPayload) Transactions() ([]bellatrix.Transaction, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix execution payload")
		}

		return v.Bellatrix.Transactions, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.Transactions, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb execution payload")
		}

		return v.Deneb.Transactions, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution payload")
		}

		return v.Electra.Transactions, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution payload")
		}

		return v.Fulu.Transactions, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution payload")
		}

		return v.Gloas.Transactions, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals of the execution payload.
func (v *VersionedExecutionPayload) Withdrawals() ([]*capella.Withdrawal, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		return nil, errors.New("no withdrawals in bellatrix")
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.Withdrawals, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb execution payload")
		}

		return v.Deneb.Withdrawals, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution payload")
		}

		return v.Electra.Withdrawals, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution payload")
		}

		return v.Fulu.Withdrawals, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution payload")
		}

		return v.Gloas.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BlobGasUsed returns the blob gas used of the execution payload.
func (v *VersionedExecutionPayload) BlobGasUsed() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no blob gas used in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no blob gas used in capella")
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.BlobGasUsed, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.BlobGasUsed, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.BlobGasUsed, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.BlobGasUsed, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExcessBlobGas returns the excess blob gas of the execution payload.
func (v *VersionedExecutionPayload) ExcessBlobGas() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no excess blob gas in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no excess blob gas in capella")
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb execution payload")
		}

		return v.Deneb.ExcessBlobGas, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return 0, errors.New("no electra execution payload")
		}

		return v.Electra.ExcessBlobGas, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return 0, errors.New("no fulu execution payload")
		}

		return v.Fulu.ExcessBlobGas, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.ExcessBlobGas, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BlockAccessList returns the block access list of the execution payload.
func (v *VersionedExecutionPayload) BlockAccessList() (gloas.BlockAccessList, error) {
	switch v.Version {
	case DataVersionPhase0:
		return gloas.BlockAccessList{}, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return gloas.BlockAccessList{}, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		return gloas.BlockAccessList{}, errors.New("no block access list in bellatrix")
	case DataVersionCapella:
		return gloas.BlockAccessList{}, errors.New("no block access list in capella")
	case DataVersionDeneb:
		return gloas.BlockAccessList{}, errors.New("no block access list in deneb")
	case DataVersionElectra:
		return gloas.BlockAccessList{}, errors.New("no block access list in electra")
	case DataVersionFulu:
		return gloas.BlockAccessList{}, errors.New("no block access list in fulu")
	case DataVersionGloas:
		if v.Gloas == nil {
			return gloas.BlockAccessList{}, errors.New("no gloas execution payload")
		}

		return v.Gloas.BlockAccessList, nil
	default:
		return gloas.BlockAccessList{}, errors.New("unknown version")
	}
}

// SlotNumber returns the slot number of the execution payload.
func (v *VersionedExecutionPayload) SlotNumber() (uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no slot number in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no slot number in capella")
	case DataVersionDeneb:
		return 0, errors.New("no slot number in deneb")
	case DataVersionElectra:
		return 0, errors.New("no slot number in electra")
	case DataVersionFulu:
		return 0, errors.New("no slot number in fulu")
	case DataVersionGloas:
		if v.Gloas == nil {
			return 0, errors.New("no gloas execution payload")
		}

		return v.Gloas.SlotNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedExecutionPayload) String() string {
	switch v.Version {
	case DataVersionPhase0:
		return "phase0"
	case DataVersionAltair:
		return "altair"
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
	case DataVersionGloas:
		if v.Gloas == nil {
			return ""
		}

		return v.Gloas.String()
	default:
		return "unknown version"
	}
}
