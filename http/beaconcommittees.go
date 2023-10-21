// Copyright © 2020 - 2023 Attestant Limited.
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
	"fmt"

	api "github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// BeaconCommittees fetches all beacon committees for the epoch at the given state.
func (s *Service) BeaconCommittees(ctx context.Context,
	opts *api.BeaconCommitteesOpts,
) (
	*api.Response[[]*apiv1.BeaconCommittee],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.State == "" {
		return nil, errors.New("no state specified")
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/committees", opts.State)
	if opts.Epoch != nil {
		url = fmt.Sprintf("%s?epoch=%d", url, *opts.Epoch)
	}

	httpResponse, err := s.get2(ctx, url)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.BeaconCommittee{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.BeaconCommittee]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
