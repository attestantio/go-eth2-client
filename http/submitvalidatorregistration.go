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

// SubmitValidatorRegistration submits a validator registration.
func (s *Service) SubmitValidatorRegistration(ctx context.Context, registration *api.VersionedSignedValidatorRegistration) error {
	var specJSON []byte
	var err error

	if registration == nil {
		return errors.New("no validator registration supplied")
	}

	switch registration.Version {
	case spec.DataVersionPhase0:
		err = errors.New("builder spec validator registration not supported")
	case spec.DataVersionAltair:
		err = errors.New("builder spec validator registration not supported")
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(registration.Bellatrix)
	default:
		err = errors.New("unknown validator registration version")
	}
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	_, err = s.post(ctx, "/eth/v1/validator/register_validator", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit validator registration")
	}

	return nil
}
