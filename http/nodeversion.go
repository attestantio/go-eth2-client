// Copyright © 2020, 2023 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/pkg/errors"
)

type nodeVersionJSON struct {
	Version string `json:"version"`
}

// NodeVersion provides the version information of the node.
func (s *Service) NodeVersion(ctx context.Context,
	opts *api.NodeVersionOpts,
) (
	*api.Response[string],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	s.nodeVersionMutex.RLock()
	if s.nodeVersion != "" {
		defer s.nodeVersionMutex.RUnlock()

		return &api.Response[string]{
			Data:     s.nodeVersion,
			Metadata: make(map[string]any),
		}, nil
	}
	s.nodeVersionMutex.RUnlock()

	s.nodeVersionMutex.Lock()
	defer s.nodeVersionMutex.Unlock()
	if s.nodeVersion != "" {
		// Someone else fetched this whilst we were waiting for the lock.
		return &api.Response[string]{
			Data:     s.nodeVersion,
			Metadata: make(map[string]any),
		}, nil
	}

	// Up to us to fetch the information.
	url := "/eth/v1/node/version"
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), nodeVersionJSON{})
	if err != nil {
		return nil, err
	}

	s.nodeVersion = data.Version

	return &api.Response[string]{
		Metadata: metadata,
		Data:     data.Version,
	}, nil
}
