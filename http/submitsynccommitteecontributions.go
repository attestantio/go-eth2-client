// Copyright © 2021, 2024 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/altair"
)

// SubmitSyncCommitteeContributions submits sync committee contributions.
func (s *Service) SubmitSyncCommitteeContributions(ctx context.Context,
	contributionAndProofs []*altair.SignedContributionAndProof,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	specJSON, err := json.Marshal(contributionAndProofs)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v1/validator/contribution_and_proofs"
	query := ""

	if _, err := s.post(ctx,
		endpoint,
		query,
		&api.CommonOpts{},
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		map[string]string{},
	); err != nil {
		return errors.Join(errors.New("failed to submit contribution and proofs"), err)
	}

	return nil
}
