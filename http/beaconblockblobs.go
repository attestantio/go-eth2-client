// Copyright © 2023 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/pkg/errors"
)

// BeaconBlockBlobs fetches the blobs given a block ID.
func (s *Service) BeaconBlockBlobs(ctx context.Context,
	opts *api.BeaconBlockBlobsOpts,
) (
	*api.Response[[]*deneb.BlobSidecar],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.Block == "" {
		return nil, errors.New("no block specified")
	}

	httpResponse, err := s.get2(ctx, fmt.Sprintf("/eth/v1/beacon/blob_sidecars/%s", opts.Block))
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*deneb.BlobSidecar{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*deneb.BlobSidecar]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
