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

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/pkg/errors"
)

// SubmitBLSToExecutionChange submits a BLS to execution address change operation.
func (s *Service) SubmitBLSToExecutionChange(ctx context.Context, blsToExecutionChange *capella.SignedBLSToExecutionChange) error {
	specJSON, err := json.Marshal(blsToExecutionChange)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	_, err = s.post(ctx, "/eth/v1/beacon/pool/bls_to_execution_changes", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit BLS to execution change")
	}

	return nil
}
