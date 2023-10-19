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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

// BlobSidecars fetches the blob sidecars given options.
func (s *Service) BlobSidecars(ctx context.Context, opts *api.BlobSidecarsOpts) ([]*deneb.BlobSidecar, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		blobSidecars, err := client.(consensusclient.BlobSidecarsProvider).BlobSidecars(ctx, opts)
		if err != nil {
			return nil, err
		}

		return blobSidecars, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return res.([]*deneb.BlobSidecar), nil
}
