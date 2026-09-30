// Copyright © 2026 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

// VersionedExecutionRequests contains versioned execution layer triggered
// requests. The container is introduced in Electra; builder requests are added
// in Gloas (EIP-8282).
type VersionedExecutionRequests struct {
	Version DataVersion
	Electra *electra.ExecutionRequests
	Fulu    *electra.ExecutionRequests
	Gloas   *gloas.ExecutionRequests
}

// IsEmpty returns true if no fork-specific execution requests are populated.
func (v *VersionedExecutionRequests) IsEmpty() bool {
	return v.Electra == nil && v.Fulu == nil && v.Gloas == nil
}

// Deposits returns the deposit requests.
func (v *VersionedExecutionRequests) Deposits() ([]*electra.DepositRequest, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("execution requests are not present before electra")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution requests")
		}

		return v.Electra.Deposits, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution requests")
		}

		return v.Fulu.Deposits, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution requests")
		}

		return v.Gloas.Deposits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawal requests.
func (v *VersionedExecutionRequests) Withdrawals() ([]*electra.WithdrawalRequest, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("execution requests are not present before electra")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution requests")
		}

		return v.Electra.Withdrawals, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution requests")
		}

		return v.Fulu.Withdrawals, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution requests")
		}

		return v.Gloas.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Consolidations returns the consolidation requests.
func (v *VersionedExecutionRequests) Consolidations() ([]*electra.ConsolidationRequest, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return nil, errors.New("execution requests are not present before electra")
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no electra execution requests")
		}

		return v.Electra.Consolidations, nil
	case DataVersionFulu:
		if v.Fulu == nil {
			return nil, errors.New("no fulu execution requests")
		}

		return v.Fulu.Consolidations, nil
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution requests")
		}

		return v.Gloas.Consolidations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BuilderDeposits returns the builder deposit requests (EIP-8282, gloas onwards).
func (v *VersionedExecutionRequests) BuilderDeposits() ([]*gloas.BuilderDepositRequest, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb,
		DataVersionElectra, DataVersionFulu:
		return nil, errors.New("builder deposit requests are not present before gloas")
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution requests")
		}

		return v.Gloas.BuilderDeposits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BuilderExits returns the builder exit requests (EIP-8282, gloas onwards).
func (v *VersionedExecutionRequests) BuilderExits() ([]*gloas.BuilderExitRequest, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb,
		DataVersionElectra, DataVersionFulu:
		return nil, errors.New("builder exit requests are not present before gloas")
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas execution requests")
		}

		return v.Gloas.BuilderExits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedExecutionRequests) String() string {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella, DataVersionDeneb:
		return ""
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
