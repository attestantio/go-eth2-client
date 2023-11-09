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

// VersionedSignedBeaconBlock contains a versioned signed beacon block.
type VersionedSignedBeaconBlock struct {
	Version   DataVersion
	Phase0    *phase0.SignedBeaconBlock
	Altair    *altair.SignedBeaconBlock
	Bellatrix *bellatrix.SignedBeaconBlock
	Capella   *capella.SignedBeaconBlock
	Deneb     *deneb.SignedBeaconBlock
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedSignedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil {
			return 0, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Slot, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil {
			return 0, errors.New("no altair block")
		}

		return v.Altair.Message.Slot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Slot, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.Slot, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return 0, errors.New("no denb block")
		}

		return v.Deneb.Message.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedSignedBeaconBlock) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 block")
		}

		return v.Phase0.Message.ProposerIndex, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair block")
		}

		return v.Altair.Message.ProposerIndex, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.ProposerIndex, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.ProposerIndex, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return 0, errors.New("no deneb block")
		}

		return v.Deneb.Message.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionBlockHash returns the block hash of the beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionBlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.ExecutionPayload.BlockHash, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockHash, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil || v.Deneb.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no denb block")
		}

		return v.Deneb.Message.Body.ExecutionPayload.BlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// ExecutionBlockNumber returns the block number of the beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionBlockNumber() (uint64, error) {
	switch v.Version {
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayload == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.ExecutionPayload.BlockNumber, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockNumber, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil || v.Deneb.Message.Body.ExecutionPayload == nil {
			return 0, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.ExecutionPayload.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionTransactions returs the execution payload transactions for the block.
func (v *VersionedSignedBeaconBlock) ExecutionTransactions() ([]bellatrix.Transaction, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("phase0 block does not have execution transactions")
	case DataVersionAltair:
		return nil, errors.New("altair block does not have execution transactions")
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.ExecutionPayload.Transactions, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.Transactions, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil || v.Deneb.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Graffiti returns the graffiti for the block.
func (v *VersionedSignedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return [32]byte{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.Graffiti, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return [32]byte{}, errors.New("no altair block")
		}

		return v.Altair.Message.Body.Graffiti, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return [32]byte{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.Graffiti, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return [32]byte{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.Graffiti, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil || v.Deneb.Message.Body.ExecutionPayload == nil {
			return [32]byte{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.Graffiti, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedSignedBeaconBlock) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.Attestations, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.Attestations, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.Attestations, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.Attestations, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.Attestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedSignedBeaconBlock) Root() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.Message.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedSignedBeaconBlock) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.Message.Body.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedSignedBeaconBlock) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.ParentRoot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.Message.ParentRoot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.ParentRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.ParentRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedSignedBeaconBlock) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.StateRoot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.Message.StateRoot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.StateRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.StateRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// RandaoReveal returns the randao reveal of the beacon block.
func (v *VersionedSignedBeaconBlock) RandaoReveal() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.RANDAOReveal, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return phase0.BLSSignature{}, errors.New("no altair block")
		}

		return v.Altair.Message.Body.RANDAOReveal, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.RANDAOReveal, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return phase0.BLSSignature{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.RANDAOReveal, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return phase0.BLSSignature{}, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// ETH1Data returns the eth1 data of the beacon block.
func (v *VersionedSignedBeaconBlock) ETH1Data() (*phase0.ETH1Data, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.ETH1Data, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.ETH1Data, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.ETH1Data, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ETH1Data, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.ETH1Data, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Deposits returns the deposits of the beacon block.
func (v *VersionedSignedBeaconBlock) Deposits() ([]*phase0.Deposit, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.Deposits, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.Deposits, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.Deposits, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.Deposits, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.Deposits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// VoluntaryExits returns the voluntary exits of the beacon block.
func (v *VersionedSignedBeaconBlock) VoluntaryExits() ([]*phase0.SignedVoluntaryExit, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Message == nil || v.Phase0.Message.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.VoluntaryExits, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Message == nil || v.Altair.Message.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.VoluntaryExits, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.VoluntaryExits, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.VoluntaryExits, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.VoluntaryExits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedSignedBeaconBlock) AttesterSlashings() ([]*phase0.AttesterSlashing, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.AttesterSlashings, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.AttesterSlashings, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.AttesterSlashings, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.AttesterSlashings, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.AttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedSignedBeaconBlock) ProposerSlashings() ([]*phase0.ProposerSlashing, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Message.Body.ProposerSlashings, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.ProposerSlashings, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.ProposerSlashings, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ProposerSlashings, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// SyncAggregate returns the sync aggregate of the beacon block.
func (v *VersionedSignedBeaconBlock) SyncAggregate() (*altair.SyncAggregate, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("phase0 block does not have sync aggregate")
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Message.Body.SyncAggregate, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Message.Body.SyncAggregate, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.SyncAggregate, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.SyncAggregate, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BLSToExecutionChanges returns the bls to execution changes of the beacon block.
func (v *VersionedSignedBeaconBlock) BLSToExecutionChanges() ([]*capella.SignedBLSToExecutionChange, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("phase0 block does not have bls to execution changes")
	case DataVersionAltair:
		return nil, errors.New("altair block does not have bls to execution changes")
	case DataVersionBellatrix:
		return nil, errors.New("bellatrix block does not have bls to execution changes")
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.BLSToExecutionChanges, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.BLSToExecutionChanges, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals of the beacon block.
func (v *VersionedSignedBeaconBlock) Withdrawals() ([]*capella.Withdrawal, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("phase0 block does not have execution withdrawals")
	case DataVersionAltair:
		return nil, errors.New("altair block does not have execution withdrawals")
	case DataVersionBellatrix:
		return nil, errors.New("bellatrix block does not have execution withdrawals")
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.Withdrawals, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil || v.Deneb.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.ExecutionPayload.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BlobKZGCommitments returns the blob KZG commitments of the beacon block.
func (v *VersionedSignedBeaconBlock) BlobKZGCommitments() ([]deneb.KZGCommitment, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("phase0 block does not have kzg commitments")
	case DataVersionAltair:
		return nil, errors.New("altair block does not have kzg commitments")
	case DataVersionBellatrix:
		return nil, errors.New("bellatrix block does not have kzg commitments")
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Message.Body.BlobKZGCommitments, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedBeaconBlock) String() string {
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
