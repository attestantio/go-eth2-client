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

// VersionedSignedValidatorRegistration contains a versioned SignedValidatorRegistrationV1.
type VersionedSignedValidatorRegistration struct {
	Version spec.BuilderVersion
	V1      *apiv1.SignedValidatorRegistration
}

// FeeRecipient returns the fee recipient of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no validator registration")
		}
		return v.V1.Message.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, errors.New("unsupported version")
	}
}

// GasLimit returns the gas limit of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) GasLimit() (uint64, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return 0, errors.New("no validator registration")
		}
		return v.V1.Message.GasLimit, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// Timestamp returns the timestamp of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) Timestamp() (time.Time, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return time.Time{}, errors.New("no validator registration")
		}
		return v.V1.Message.Timestamp, nil
	default:
		return time.Time{}, errors.New("unsupported version")
	}
}

// PubKey returns the public key of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) PubKey() (phase0.BLSPubKey, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.BLSPubKey{}, errors.New("no validator registration")
		}
		return v.V1.Message.Pubkey, nil
	default:
		return phase0.BLSPubKey{}, errors.New("unsupported version")
	}
}

// Root returns the root of the validator registration.
func (v *VersionedSignedValidatorRegistration) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.Root{}, errors.New("no V1 registration")
		}
		return v.V1.Message.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}
