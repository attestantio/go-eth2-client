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
	"time"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedValidatorRegistration contains a versioned SignedValidatorRegistrationV1.
type VersionedSignedValidatorRegistration struct {
	Version spec.BuilderVersion                `json:"version"`
	V1      *apiv1.SignedValidatorRegistration `json:"v1"`
}

// FeeRecipient returns the fee recipient of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.V1.Message.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// GasLimit returns the gas limit of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) GasLimit() (uint64, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return 0, ErrDataMissing
		}

		return v.V1.Message.GasLimit, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) Timestamp() (time.Time, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return time.Time{}, ErrDataMissing
		}

		return v.V1.Message.Timestamp, nil
	default:
		return time.Time{}, ErrUnsupportedVersion
	}
}

// PubKey returns the public key of the signed validator registration.
func (v *VersionedSignedValidatorRegistration) PubKey() (phase0.BLSPubKey, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.BLSPubKey{}, ErrDataMissing
		}

		return v.V1.Message.Pubkey, nil
	default:
		return phase0.BLSPubKey{}, ErrUnsupportedVersion
	}
}

// Root returns the root of the validator registration.
func (v *VersionedSignedValidatorRegistration) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.V1.Message.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}
