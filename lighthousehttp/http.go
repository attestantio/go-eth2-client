// Copyright © 2020 Attestant Limited.
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

package lighthousehttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// get sends an HTTP get request and returns the body.
func (s *Service) get(ctx context.Context, endpoint string) (io.ReadCloser, error) {
	log.Trace().Str("endpoint", endpoint).Msg("GET request")

	reference, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()
	log.Trace().Str("url", url).Msg("GET request")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to lighthouse HTTP endpoint")
	}
	statusFamily := resp.StatusCode / 100
	if statusFamily != 2 {
		if err := resp.Body.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close response body")
		}
		return nil, fmt.Errorf("HTTP call to lighthouse endpoint failed with status %d", resp.StatusCode)
	}

	if e := log.Trace(); e.Enabled() {
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Trace().Str("response", string(data)).Msg("GET response")
			if err := resp.Body.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close response body")
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(data))
		}
	}

	return resp.Body, nil
}

// post sends an HTTP post request and returns the body.
func (s *Service) post(ctx context.Context, endpoint string, body io.Reader) (io.ReadCloser, error) {
	reference, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()

	if e := log.Trace(); e.Enabled() {
		if body == nil {
			e.Str("url", url).Msg("POST request")
		} else {
			bodyBytes, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, errors.New("failed to read request body")
			}
			body = bytes.NewReader(bodyBytes)
			e.Str("url", url).Str("body", string(bodyBytes)).Msg("POST request")
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create POST request")
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to lighthouse HTTP endpoint")
	}
	statusFamily := resp.StatusCode / 100
	if statusFamily != 2 {
		if err := resp.Body.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close response body")
		}
		return nil, fmt.Errorf("HTTP call to lighthouse endpoint failed with status %d", resp.StatusCode)
	}

	if e := log.Trace(); e.Enabled() {
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Trace().Str("response", string(data)).Msg("POST response")
			if err := resp.Body.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close response body")
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(data))
		}
	}

	return resp.Body, nil
}
