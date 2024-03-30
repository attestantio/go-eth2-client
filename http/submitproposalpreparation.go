// Copyright © 2022, 2024 Attestant Limited.
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
	"errors"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// SubmitProposalPreparations provides the beacon node with information required if a proposal for the given validators
// shows up in the next epoch.
func (s *Service) SubmitProposalPreparations(ctx context.Context, preparations []*apiv1.ProposalPreparation) error {
	if err := s.assertIsActive(ctx); err != nil {
		return err
	}

	var reqBodyReader bytes.Buffer
	if err := json.NewEncoder(&reqBodyReader).Encode(preparations); err != nil {
		return errors.Join(errors.New("failed to encode proposal preparations"), err)
	}

	_, err := s.post(ctx, "/eth/v1/validator/prepare_beacon_proposer", &reqBodyReader)
	if err != nil {
		return errors.Join(errors.New("failed to send proposal preparations"), err)
	}

	return nil
}
