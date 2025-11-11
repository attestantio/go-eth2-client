// Copyright © 2025 Attestant Limited.
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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// Blobs fetches the blobs given options.
func (s *Service) Blobs(ctx context.Context,
	opts *api.BlobsOpts,
) (
	*api.Response[apiv1.Blobs],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/blobs/%s", opts.Block)

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	var response *api.Response[apiv1.Blobs]

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.blobsFromSSZ(httpResponse)
	case ContentTypeJSON:
		response, err = s.blobsFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (*Service) blobsFromSSZ(res *httpResponse) (*api.Response[apiv1.Blobs], error) {
	response := &api.Response[apiv1.Blobs]{}

	if len(res.body) == 0 {
		// This is a valid response when there are no blobs for the request.
		response.Data = make(apiv1.Blobs, 0)
		response.Metadata = make(map[string]any)

		return response, nil
	}

	if err := response.Data.UnmarshalSSZ(res.body); err != nil {
		return nil, errors.Join(errors.New("failed to decode blobs"), err)
	}

	return response, nil
}

func (*Service) blobsFromJSON(res *httpResponse) (*api.Response[apiv1.Blobs], error) {
	response := &api.Response[apiv1.Blobs]{}

	var err error

	response.Data, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), apiv1.Blobs{})
	if err != nil {
		return nil, err
	}

	return response, nil
}
