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

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedExecutionPayloadEnvelope contains a versioned signed
// execution payload envelope (EIP-7732). The envelope is a gloas-onwards
// container; this wrapper exists so callers can treat it uniformly with
// other versioned types.
type VersionedSignedExecutionPayloadEnvelope struct {
	Version DataVersion
	Gloas   *gloas.SignedExecutionPayloadEnvelope
}

// IsEmpty returns true if no fork-specific envelope is populated.
func (v *VersionedSignedExecutionPayloadEnvelope) IsEmpty() bool {
	return v.Gloas == nil
}

// String returns a string version of the structure.
func (v *VersionedSignedExecutionPayloadEnvelope) String() string {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella,
		DataVersionDeneb, DataVersionElectra, DataVersionFulu:
		return ""
	case DataVersionGloas:
		if v.Gloas == nil {
			return ""
		}

		return v.Gloas.String()
	default:
		return "unknown version"
	}
}

// Message returns the inner ExecutionPayloadEnvelope.
func (v *VersionedSignedExecutionPayloadEnvelope) Message() (*gloas.ExecutionPayloadEnvelope, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return nil, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return nil, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return nil, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return nil, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return nil, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return nil, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Message, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Payload returns the execution payload contained in the signed envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) Payload() (*gloas.ExecutionPayload, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return nil, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return nil, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return nil, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return nil, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return nil, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return nil, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Message.Payload, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ExecutionRequests returns the execution requests contained in the signed envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) ExecutionRequests() (*VersionedExecutionRequests, error) {
	switch v.Version {
	case DataVersionPhase0:
		return nil, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return nil, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return nil, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return nil, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return nil, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return nil, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return nil, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return nil, errors.New("no gloas signed execution payload envelope")
		}

		return &VersionedExecutionRequests{
			Version: DataVersionGloas,
			Gloas:   v.Gloas.Message.ExecutionRequests,
		}, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BuilderIndex returns the builder index of the signed envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) BuilderIndex() (gloas.BuilderIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		return 0, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return 0, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return 0, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return 0, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return 0, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return 0, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return 0, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return 0, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Message.BuilderIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// BeaconBlockRoot returns the beacon block root of the signed envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) BeaconBlockRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Root{}, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return phase0.Root{}, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return phase0.Root{}, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return phase0.Root{}, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return phase0.Root{}, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return phase0.Root{}, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return phase0.Root{}, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return phase0.Root{}, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Message.BeaconBlockRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// ParentBeaconBlockRoot returns the parent beacon block root of the signed envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) ParentBeaconBlockRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.Root{}, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return phase0.Root{}, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return phase0.Root{}, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return phase0.Root{}, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return phase0.Root{}, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return phase0.Root{}, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return phase0.Root{}, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Message == nil {
			return phase0.Root{}, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Message.ParentBeaconBlockRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// Signature returns the signature on the envelope.
func (v *VersionedSignedExecutionPayloadEnvelope) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in phase0")
	case DataVersionAltair:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in altair")
	case DataVersionBellatrix:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in bellatrix")
	case DataVersionCapella:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in capella")
	case DataVersionDeneb:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in deneb")
	case DataVersionElectra:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in electra")
	case DataVersionFulu:
		return phase0.BLSSignature{}, errors.New("no execution payload envelope in fulu")
	case DataVersionGloas:
		if v.Gloas == nil {
			return phase0.BLSSignature{}, errors.New("no gloas signed execution payload envelope")
		}

		return v.Gloas.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}
