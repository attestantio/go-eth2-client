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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedPayloadAttestationData contains a versioned payload attestation
// data (EIP-7732).  This is the datum a payload timeliness committee member
// signs, and is a gloas-onwards container; the wrapper exists so callers can
// branch on Version uniformly with the other versioned types.
type VersionedPayloadAttestationData struct {
	Version DataVersion
	Gloas   *gloas.PayloadAttestationData
}

// IsEmpty returns true if no fork-specific payload attestation data is populated.
func (v *VersionedPayloadAttestationData) IsEmpty() bool {
	return v.Gloas == nil
}

// String returns a string version of the structure.
func (v *VersionedPayloadAttestationData) String() string {
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

// Slot returns the slot of the payload attestation data.
func (v *VersionedPayloadAttestationData) Slot() (phase0.Slot, error) {
	data, err := v.data()
	if err != nil {
		return 0, err
	}

	return data.Slot, nil
}

// BeaconBlockRoot returns the root of the beacon block being attested to.
func (v *VersionedPayloadAttestationData) BeaconBlockRoot() (phase0.Root, error) {
	data, err := v.data()
	if err != nil {
		return phase0.Root{}, err
	}

	return data.BeaconBlockRoot, nil
}

// PayloadPresent returns true if the execution payload was revealed on time.
func (v *VersionedPayloadAttestationData) PayloadPresent() (bool, error) {
	data, err := v.data()
	if err != nil {
		return false, err
	}

	return data.PayloadPresent, nil
}

// BlobDataAvailable returns true if the blob data for the payload was available.
func (v *VersionedPayloadAttestationData) BlobDataAvailable() (bool, error) {
	data, err := v.data()
	if err != nil {
		return false, err
	}

	return data.BlobDataAvailable, nil
}

// data returns the fork-specific payload attestation data, or an error
// explaining why there is none.  Every accessor funnels through here so that
// the version and nil checks cannot drift between them.
func (v *VersionedPayloadAttestationData) data() (*gloas.PayloadAttestationData, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella,
		DataVersionDeneb, DataVersionElectra, DataVersionFulu:
		return nil, fmt.Errorf("no payload attestation data in %s", v.Version)
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas payload attestation data")
		}

		return v.Gloas, nil
	default:
		return nil, errors.New("unknown version")
	}
}
