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
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/pkg/errors"
)

type beaconBlockBlobsJSON struct {
	Data []*deneb.BlobSidecar `json:"data"`
}

// BeaconBlockBlobs fetches the blobs given a block ID.
func (s *Service) BeaconBlockBlobs(ctx context.Context, blockID string) ([]*deneb.BlobSidecar, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/blob_sidecars/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request blobs")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp beaconBlockBlobsJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse blobs")
	}

	// Data is not guaranteed to be returned in indx order, so fix that.
	sort.Slice(resp.Data, func(i int, j int) bool {
		return resp.Data[i].Index < resp.Data[j].Index
	})

	return resp.Data, nil
}
