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

package tekuhttp

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
// The cancel function must be called once use of the returned reader is complete.
func (s *Service) get(ctx context.Context, endpoint string) (io.Reader, context.CancelFunc, error) {
	log.Trace().Str("endpoint", endpoint).Msg("GET request")

	reference, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()
	log.Trace().Str("url", url).Msg("GET request")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	req, err := http.NewRequestWithContext(opCtx, http.MethodGet, url, nil)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to create GET request")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to connect to teku GET endpoint")
	}
	statusFamily := resp.StatusCode / 100
	if statusFamily != 2 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to read response body")
		}
		cancel()
		return nil, nil, fmt.Errorf("HTTP call to teku endpoint failed with status %d: %s", resp.StatusCode, string(data))
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

	return resp.Body, cancel, nil
}

// post sends an HTTP post request and returns the body.
// The cancel function must be called once use of the returned reader is complete.
func (s *Service) post(ctx context.Context, endpoint string, body io.Reader) (io.Reader, context.CancelFunc, error) {
	reference, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()

	if e := log.Trace(); e.Enabled() {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, nil, errors.New("failed to read request body")
		}
		body = bytes.NewReader(bodyBytes)

		e.Str("url", url).Str("body", string(bodyBytes)).Msg("POST request")
	}

	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	req, err := http.NewRequestWithContext(opCtx, http.MethodPost, url, body)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to create POST request")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to connect to teku HTTP endpoint")
	}
	statusFamily := resp.StatusCode / 100
	if statusFamily != 2 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to read response body")
		}
		cancel()
		return nil, nil, fmt.Errorf("HTTP call to teku endpoint failed with status %d: %s", resp.StatusCode, string(data))
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

	return resp.Body, cancel, nil
}
