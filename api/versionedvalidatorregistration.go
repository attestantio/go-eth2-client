// Copyright © 2022 Attestant Limited.
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
	"time"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedValidatorRegistration contains a versioned ValidatorRegistrationV1.
type VersionedValidatorRegistration struct {
	Version spec.BuilderVersion
	V1      *apiv1.ValidatorRegistration
}

// IsEmpty returns true if there is no block.
func (v *VersionedValidatorRegistration) IsEmpty() bool {
	return v.V1 == nil
}

// FeeRecipient returns the fee recipient of the validator registration.
func (v *VersionedValidatorRegistration) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no validator registration")
		}

		return v.V1.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, errors.New("unsupported version")
	}
}

// GasLimit returns the gas limit of the validator registration.
func (v *VersionedValidatorRegistration) GasLimit() (uint64, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return 0, errors.New("no validator registration")
		}

		return v.V1.GasLimit, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// Timestamp returns the timestamp of the validator registration.
func (v *VersionedValidatorRegistration) Timestamp() (time.Time, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return time.Time{}, errors.New("no validator registration")
		}

		return v.V1.Timestamp, nil
	default:
		return time.Time{}, errors.New("unsupported version")
	}
}

// PubKey returns the public key of the validator registration.
func (v *VersionedValidatorRegistration) PubKey() (phase0.BLSPubKey, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.BLSPubKey{}, errors.New("no validator registration")
		}

		return v.V1.Pubkey, nil
	default:
		return phase0.BLSPubKey{}, errors.New("unsupported version")
	}
}

// Root returns the root of the validator registration.
func (v *VersionedValidatorRegistration) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.Root{}, errors.New("no V1 registration")
		}

		return v.V1.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}
