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

package http

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
)

// SubmitValidatorRegistrations submits a validator registration.
func (s *Service) SubmitValidatorRegistrations(ctx context.Context, registrations []*api.VersionedSignedValidatorRegistration) error {
	if len(registrations) == 0 {
		return errors.New("no registrations supplied")
	}

	// Unwrap versioned registrations.
	var version *spec.BuilderVersion
	var unversionedRegistrations []interface{}

	for i := range registrations {
		if registrations[i] == nil {
			return errors.New("nil registration supplied")
		}

		// Ensure consistent versioning.
		if version == nil {
			version = &registrations[i].Version
		} else if *version != registrations[i].Version {
			return errors.New("registrations must all be of the same version")
		}

		// Append to unversionedRegistrations.
		switch registrations[i].Version {
		case spec.BuilderVersionV1:
			unversionedRegistrations = append(unversionedRegistrations, registrations[i].V1)
		default:
			return errors.New("unknown validator registration version")
		}
	}

	specJSON, err := json.Marshal(unversionedRegistrations)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}
	_, err = s.post(ctx, "/eth/v1/validator/register_validator", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit validator registration")
	}

	return nil
}
