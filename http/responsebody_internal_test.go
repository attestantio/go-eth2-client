// Copyright © 2026 Attestant Limited.
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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestResponseLimitFollowsResponseEncoding checks that GET and POST pick the
// response limit from the encoding the node returned, not the one requested.
func TestResponseLimitFollowsResponseEncoding(t *testing.T) {
	limit := responseLimit{ssz: 4, json: 8}

	tests := []struct {
		name          string
		contentType   string
		contentLength string
		body          string
		unbounded     bool
		err           string
	}{
		{
			name:        "SSZAtLimit",
			contentType: "application/octet-stream",
			body:        "1234",
		},
		{
			name:        "SSZExceedsLimit",
			contentType: "application/octet-stream",
			body:        "12345",
			err:         "response body exceeds 4 bytes",
		},
		{
			name:        "JSONAboveSSZLimit",
			contentType: "application/json",
			body:        `{"a":1}`,
		},
		{
			name:        "JSONAtLimit",
			contentType: "application/json",
			body:        `{"a":12}`,
		},
		{
			name:        "JSONExceedsLimit",
			contentType: "application/json",
			body:        `{"a":123}`,
			err:         "response body exceeds 8 bytes",
		},
		{
			name: "MissingContentTypeUsesJSONLimit",
			body: `{"a":1}`,
		},
		{
			// The declared length is rejected before the body is read, so the
			// short body never surfaces as an unexpected EOF.
			name:          "DeclaredLengthExceedsLimit",
			contentType:   "application/octet-stream",
			contentLength: "1000",
			err:           "response body exceeds 4 bytes",
		},
		{
			name:        "UnboundedWhenZero",
			contentType: "application/octet-stream",
			unbounded:   true,
			body:        strings.Repeat("1", 1024),
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if test.contentType != "" {
				w.Header().Set("Content-Type", test.contentType)
			} else {
				// Stop net/http sniffing a content type for the body.
				w.Header()["Content-Type"] = nil
			}
			if test.contentLength != "" {
				w.Header().Set("Content-Length", test.contentLength)
			}
			_, _ = w.Write([]byte(test.body))
		}))
		t.Cleanup(server.Close)

		base, err := url.Parse(server.URL)
		require.NoError(t, err)
		s := &Service{log: zerolog.Nop(), base: base, client: server.Client(), timeout: time.Second}

		testLimit := limit
		if test.unbounded {
			testLimit = responseLimit{}
		}

		requests := map[string]func() (*httpResponse, error){
			"GET": func() (*httpResponse, error) {
				return s.getWithResponseLimit(context.Background(), "/test", "", &api.CommonOpts{}, true, testLimit)
			},
			"POST": func() (*httpResponse, error) {
				return s.postWithResponseLimit(context.Background(), "/test", "", &api.CommonOpts{},
					strings.NewReader("{}"), ContentTypeJSON, nil, testLimit)
			},
		}

		for method, request := range requests {
			t.Run(test.name+"/"+method, func(t *testing.T) {
				res, err := request()
				if test.err != "" {
					require.ErrorContains(t, err, test.err)
				} else {
					require.NoError(t, err)
					require.Equal(t, []byte(test.body), res.body)
				}
			})
		}
	}
}

// TestEPBSProposalResponseLimitsAreSurvivable pins the property that makes the
// block-production limits a defence rather than a comment.  For a body without
// a declared length, readResponseBody reaches its length check only after
// io.ReadAll has buffered the whole body, and ReadAll grows its buffer
// geometrically, so peak allocation is close to twice the limit.  A limit taken from the protocol's theoretical maximum
// (~2GiB of SSZ, ~5GiB of JSON hex) is therefore never reached -- the process
// dies first -- and a body admitted just under one is equally fatal.
//
// The ceiling below is an allocation budget, not a protocol bound: every live
// preset produces a payload-included response orders of magnitude smaller, so
// the JSON limit sitting exactly at half the budget is the deliberate ceiling
// rather than a coincidence.  Raising either limit past it is the regression
// this test exists to catch.
func TestEPBSProposalResponseLimitsAreSurvivable(t *testing.T) {
	const maxProposalResponsePeakAllocation = 1024 * 1024 * 1024

	tests := []struct {
		name  string
		limit int
	}{
		{name: "SSZ", limit: maxEPBSProposalResponseSize},
		{name: "JSON", limit: maxEPBSProposalJSONResponseSize},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.LessOrEqual(t, 2*test.limit, maxProposalResponsePeakAllocation)
		})
	}
}
