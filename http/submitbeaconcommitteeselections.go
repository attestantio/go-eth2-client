// Copyright © 2020, 2024 Attestant Limited.
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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// SubmitBeaconCommitteeSelections submits beacon committee selections.
func (s *Service) SubmitBeaconCommitteeSelections(ctx context.Context,
	selections []*apiv1.BeaconCommitteeSelection,
) (
	*api.Response[[]*apiv1.BeaconCommitteeSelection],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	specJSON, err := json.Marshal(selections)
	if err != nil {
		return nil, errors.Join(errors.New("failed to encode beacon committee selections"), err)
	}

	endpoint := "/eth/v1/validator/beacon_committee_selections"
	query := ""

	httpResponse, err := s.post(ctx,
		endpoint,
		query,
		&api.CommonOpts{},
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		map[string]string{},
	)

	if err != nil {
		return nil, errors.Join(errors.New("failed to request beacon committee selections"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.BeaconCommitteeSelection{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.BeaconCommitteeSelection]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
