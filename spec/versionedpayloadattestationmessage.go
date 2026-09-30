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

// VersionedPayloadAttestationMessage contains a versioned payload attestation
// message (EIP-7732): the unaggregated form a single payload timeliness
// committee member broadcasts.  It is a gloas-onwards container; the wrapper
// exists so that the Eth-Consensus-Version header required when submitting one
// is derived from the caller's data rather than hard-coded.
type VersionedPayloadAttestationMessage struct {
	Version DataVersion
	Gloas   *gloas.PayloadAttestationMessage
}

// IsEmpty returns true if no fork-specific message is populated.
func (v *VersionedPayloadAttestationMessage) IsEmpty() bool {
	return v.Gloas == nil
}

// String returns a string version of the structure.
func (v *VersionedPayloadAttestationMessage) String() string {
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

// ValidatorIndex returns the index of the attesting validator.
func (v *VersionedPayloadAttestationMessage) ValidatorIndex() (phase0.ValidatorIndex, error) {
	message, err := v.message()
	if err != nil {
		return 0, err
	}

	return message.ValidatorIndex, nil
}

// Data returns the data being attested to.
func (v *VersionedPayloadAttestationMessage) Data() (*gloas.PayloadAttestationData, error) {
	message, err := v.message()
	if err != nil {
		return nil, err
	}

	if message.Data == nil {
		return nil, errors.New("no gloas payload attestation data")
	}

	return message.Data, nil
}

// Signature returns the signature of the message.
func (v *VersionedPayloadAttestationMessage) Signature() (phase0.BLSSignature, error) {
	message, err := v.message()
	if err != nil {
		return phase0.BLSSignature{}, err
	}

	return message.Signature, nil
}

// message returns the fork-specific message, or an error explaining why there
// is none.  Every accessor funnels through here so that the version and nil
// checks cannot drift between them.
func (v *VersionedPayloadAttestationMessage) message() (*gloas.PayloadAttestationMessage, error) {
	switch v.Version {
	case DataVersionPhase0, DataVersionAltair, DataVersionBellatrix, DataVersionCapella,
		DataVersionDeneb, DataVersionElectra, DataVersionFulu:
		return nil, fmt.Errorf("no payload attestation message in %s", v.Version)
	case DataVersionGloas:
		if v.Gloas == nil {
			return nil, errors.New("no gloas payload attestation message")
		}

		return v.Gloas, nil
	default:
		return nil, errors.New("unknown version")
	}
}
