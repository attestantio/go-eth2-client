// Copyright © 2023, 2024 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/deneb"
	dynssz "github.com/pk910/dynamic-ssz"
)

// BlobSidecars fetches the blobs sidecars given options.
func (s *Service) BlobSidecars(ctx context.Context,
	opts *api.BlobSidecarsOpts,
) (
	*api.Response[[]*deneb.BlobSidecar],
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

	endpoint := fmt.Sprintf("/eth/v1/beacon/blob_sidecars/%s", opts.Block)

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	var response *api.Response[[]*deneb.BlobSidecar]

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.blobSidecarsFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		response, err = s.blobSidecarsFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) blobSidecarsFromSSZ(ctx context.Context, res *httpResponse) (*api.Response[[]*deneb.BlobSidecar], error) {
	response := &api.Response[[]*deneb.BlobSidecar]{}

	if len(res.body) == 0 {
		// This is a valid response when there are no blobs for the request.
		response.Data = make([]*deneb.BlobSidecar, 0)
		response.Metadata = make(map[string]any)

		return response, nil
	}

	data := &api.BlobSidecars{}

	var err error
	if s.customSpecSupport {
		var specs *api.Response[map[string]any]
		if specs, err = s.Spec(ctx, &api.SpecOpts{}); err != nil {
			return nil, errors.Join(errors.New("failed to request specs"), err)
		}
		dynSsz := dynssz.NewDynSsz(specs.Data)
		err = dynSsz.UnmarshalSSZ(data, res.body)
	} else {
		err = data.UnmarshalSSZ(res.body)
	}

	if err != nil {
		return nil, errors.Join(errors.New("failed to decode blob sidecars"), err)
	}

	response.Data = data.Sidecars

	return response, nil
}

func (*Service) blobSidecarsFromJSON(res *httpResponse) (*api.Response[[]*deneb.BlobSidecar], error) {
	response := &api.Response[[]*deneb.BlobSidecar]{}

	var err error

	response.Data, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), []*deneb.BlobSidecar{})
	if err != nil {
		return nil, err
	}

	return response, nil
}
