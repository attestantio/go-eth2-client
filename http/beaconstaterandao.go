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
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type beaconStateRandaoJSON struct {
	Randao phase0.Root `json:"randao"`
}

// BeaconStateRandao fetches the beacon state RANDAO given a set of options.
func (s *Service) BeaconStateRandao(ctx context.Context, opts *api.BeaconStateRandaoOpts) (*api.Response[*phase0.Root], error) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.State == "" {
		return nil, errors.New("no state specified")
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/randao", opts.State)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), beaconStateRandaoJSON{})
	if err != nil {
		return nil, err
	}

	return &api.Response[*phase0.Root]{
		Data:     &data.Randao,
		Metadata: metadata,
	}, nil
}
