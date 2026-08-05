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

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedPayloadAttestation contains a versioned payload attestation
// (EIP-7732): the aggregate form held in the operations pool, in which a
// single signature covers the payload timeliness committee members named by
// the aggregation bits.  It is a gloas-onwards container; the wrapper exists
// so callers can branch on Version uniformly with the other versioned types.
type VersionedPayloadAttestation struct {
	Version DataVersion
	Gloas   *gloas.PayloadAttestation
}

// IsEmpty returns true if no fork-specific payload attestation is populated.
func (v *VersionedPayloadAttestation) IsEmpty() bool {
	return v.Gloas == nil
}

// String returns a string version of the structure.
func (v *VersionedPayloadAttestation) String() string {
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

// Data returns the data being attested to.
func (v *VersionedPayloadAttestation) Data() (*gloas.PayloadAttestationData, error) {
	attestation, err := v.attestation()
	if err != nil {
		return nil, err
	}

	if attestation.Data == nil {
		return nil, errors.New("no gloas payload attestation data")
	}

	return attestation.Data, nil
}

// AggregationBits returns the bits of the payload timeliness committee covered
// by the attestation's signature.
//
// The returned value's own BitAt, SetBitAt, Len and Shift methods are mainnet-only and
// must not be used; read and write participation through gloas.PayloadAttestation's
// AttestedAt, SetAttestedAt and PTCSize instead, which hold at any preset width.  Note
// that width is the producing node's claim rather than a verified fact; see PTCSize.
func (v *VersionedPayloadAttestation) AggregationBits() (bitfield.Bitvector512, error) {
	attestation, err := v.attestation()
	if err != nil {
		return nil, err
	}

	return attestation.AggregationBits, nil
}

// Signature returns the aggregate signature of the payload attestation.
func (v *VersionedPayloadAttestation) Signature() (phase0.BLSSignature, error) {
	attestation, err := v.attestation()
	if err != nil {
		return phase0.BLSSignature{}, err
	}

	return attestation.Signature, nil
}

// attestation returns the fork-specific payload attestation, or an error
// explaining why there is none.  Every accessor funnels through here so that
// the version and nil checks cannot drift between them.
func (v *VersionedPayloadAttestation) attestation() (*gloas.PayloadAttestation, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella,
		DataVersionDeneb, DataVersionElectra, DataVersionFulu:
		return nil, fmt.Errorf("no payload attestation in %s", v.Version)
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas payload attestation")
		}

		return v.Gloas, nil
	default:
		return nil, errors.New("unknown version")
	}
}
