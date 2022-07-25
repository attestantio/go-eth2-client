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

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBlindedBeaconBlock contains a versioned blinded beacon block.
type VersionedBlindedBeaconBlock struct {
	Version   spec.DataVersion
	Bellatrix *apiv1.BlindedBeaconBlock
}

// IsEmpty returns true if there is no block.
func (v *VersionedBlindedBeaconBlock) IsEmpty() bool {
	return v.Bellatrix == nil
}

// Slot returns the slot of the beacon block.
func (v *VersionedBlindedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Slot, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedBlindedBeaconBlock) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Body.Attestations, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedBlindedBeaconBlock) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedBlindedBeaconBlock) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedBlindedBeaconBlock) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedBlindedBeaconBlock) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// String returns a string version of the structure.
func (v *VersionedBlindedBeaconBlock) String() string {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}
		return v.Bellatrix.String()
	default:
		return "unknown version"
	}
}
