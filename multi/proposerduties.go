// Copyright © 2021 Attestant Limited.
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

package multi

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"strings"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// ProposerDutiesV2 obtains proposer duties for the given epoch using the v2 API.
// If opts.Indices is empty all duties are returned, otherwise only matching duties are returned.
//
// v2 is an optional endpoint, so a client that cannot serve it is skipped rather
// than deactivated: support is a static property of the node, not a sign that it
// is unhealthy, and evicting it here would take it out of the rotation for every
// other call until the next recheck.
func (s *Service) ProposerDutiesV2(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*apiv1.ProposerDuty],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		provider, supported := client.(consensusclient.ProposerDutiesV2Provider)
		if !supported {
			return nil, fmt.Errorf("%s@%s does not implement proposer duties v2: %w",
				client.Name(), client.Address(), ErrCallNotSupported)
		}

		duties, err := provider.ProposerDutiesV2(ctx, opts)
		if err != nil {
			// *http.Service implements the provider whether or not the node
			// serves the route, so a missing route arrives as 404 or 405 on
			// that path.  A 404 from a later call, such as the spec fetch,
			// is a different failure and must not be rewritten.
			var apiErr *api.Error
			if errors.As(err, &apiErr) && missingProposerDutiesV2(apiErr) {
				return nil, fmt.Errorf("%s@%s does not serve proposer duties v2: %w",
					client.Name(), client.Address(), errors.Join(ErrCallNotSupported, err))
			}

			return nil, err
		}

		return duties, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[[]*apiv1.ProposerDuty])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}

// missingProposerDutiesV2 reports a 404 or 405 from the v2 duties route itself.
// Other 404s, including one from the spec endpoint after a duties response,
// are not a missing route.
func missingProposerDutiesV2(apiErr *api.Error) bool {
	if apiErr.StatusCode != nethttp.StatusNotFound && apiErr.StatusCode != nethttp.StatusMethodNotAllowed {
		return false
	}

	return strings.HasPrefix(apiErr.Endpoint, "/eth/v2/validator/duties/proposer")
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validatorIndices is empty all duties are returned, otherwise only matching duties are returned.
func (s *Service) ProposerDuties(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*apiv1.ProposerDuty],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		block, err := client.(consensusclient.ProposerDutiesProvider).ProposerDuties(ctx, opts)
		if err != nil {
			return nil, err
		}

		return block, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[[]*apiv1.ProposerDuty])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
