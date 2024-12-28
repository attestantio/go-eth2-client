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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// SubmitProposalSlashing submits a proposal slashing.
func (s *Service) SubmitProposalSlashing(ctx context.Context, slashing *phase0.ProposerSlashing) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	specJSON, err := json.Marshal(slashing)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v1/beacon/pool/proposer_slashings"
	query := ""

	if _, err := s.post(ctx,
		endpoint,
		query,
		&api.CommonOpts{},
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		map[string]string{},
	); err != nil {
		return errors.Join(errors.New("failed to submit proposal slashing"), err)
	}

	return nil
}
