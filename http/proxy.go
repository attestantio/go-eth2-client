// Copyright © 2020 - 2025 Attestant Limited.
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
	"fmt"
	"net/http"

	"github.com/attestantio/go-eth2-client/api"
)

// Proxy performs an HTTP proxy request and returns the response.
func (s *Service) Proxy(ctx context.Context, req *http.Request) (*http.Response, error) {
	return s.proxy(ctx, req)
}

// proxy performs an HTTP proxy request using a reverse proxy and returns the response.
//
//nolint:revive
func (s *Service) proxy(ctx context.Context, req *http.Request) (*http.Response, error) {
	endpoint := req.URL.Path
	query := req.URL.Query().Encode()

	var httpResponse *httpResponse
	var err error
	switch req.Method {
	case http.MethodGet:
		httpResponse, err = s.get(ctx, endpoint, query, &api.CommonOpts{}, false)
	case http.MethodPost:
		headers := make(map[string]string)
		for k, v := range req.Header {
			headers[k] = v[0]
		}
		httpResponse, err = s.post(ctx, endpoint, query, &api.CommonOpts{}, req.Body, ContentTypeJSON, headers)
	default:
		err = fmt.Errorf("unsupported method %s for proxy", req.Method)
	}

	if err != nil {
		return nil, err
	}

	return &httpResponse.raw, nil
}
