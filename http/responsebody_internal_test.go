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
		// duplicateType sends the Content-Type header twice.
		duplicateType bool
		body          string
		limit         *responseLimit
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
			name:        "JSONWithParametersUsesJSONLimit",
			contentType: "application/json; charset=utf-8",
			body:        `{"a":1}`,
		},
		{
			name:        "MixedCaseJSONUsesJSONLimit",
			contentType: "Application/JSON",
			body:        `{"a":1}`,
		},
		{
			name:        "SSZAboveJSONLimit",
			contentType: "application/octet-stream",
			limit:       &responseLimit{ssz: 8, json: 4},
			body:        "12345",
		},
		{
			// An unlabelled body could be either encoding, so it gets the
			// smaller limit to bound memory.
			name: "MissingContentTypeUsesSmallerLimit",
			body: "12345",
			err:  "response body exceeds 4 bytes",
		},
		{
			name:  "MissingContentTypeUsesSmallerJSONLimit",
			limit: &responseLimit{ssz: 8, json: 4},
			body:  "12345",
			err:   "response body exceeds 4 bytes",
		},
		{
			name:          "DuplicateContentTypeUsesSmallerLimit",
			contentType:   "application/json",
			duplicateType: true,
			body:          "12345",
			err:           "response body exceeds 4 bytes",
		},
		{
			name:        "UnrecognisedContentTypeUsesSmallerLimit",
			contentType: "text/plain",
			body:        "12345",
			err:         "response body exceeds 4 bytes",
		},
		{
			name:  "MissingContentTypeUsesOnlySetLimit",
			limit: &responseLimit{json: 4},
			body:  "12345",
			err:   "response body exceeds 4 bytes",
		},
		{
			name:  "MissingContentTypeUsesOnlySetSSZLimit",
			limit: &responseLimit{ssz: 4},
			body:  "12345",
			err:   "response body exceeds 4 bytes",
		},
		{
			name:  "MissingContentTypeUnboundedWhenZero",
			limit: &responseLimit{},
			body:  `{"a":"` + strings.Repeat("1", 1024) + `"}`,
		},
		{
			name:          "DeclaredLengthAtLimit",
			contentType:   "application/octet-stream",
			contentLength: "4",
			body:          "1234",
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
			limit:       &responseLimit{},
			body:        strings.Repeat("1", 1024),
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if test.contentType != "" {
				w.Header().Set("Content-Type", test.contentType)
				if test.duplicateType {
					w.Header().Add("Content-Type", test.contentType)
				}
			} else {
				// Stop net/http sniffing a content type for the body.
				w.Header()["Content-Type"] = nil
			}
			if test.contentLength != "" {
				w.Header().Set("Content-Length", test.contentLength)
			} else {
				// Flushing first sends the body chunked, without a declared
				// length, so the limit is enforced while reading the body.
				w.(http.Flusher).Flush()
			}
			_, _ = w.Write([]byte(test.body))
		}))
		t.Cleanup(server.Close)

		base, err := url.Parse(server.URL)
		require.NoError(t, err)
		s := &Service{log: zerolog.Nop(), base: base, client: server.Client(), timeout: time.Second}

		testLimit := limit
		if test.limit != nil {
			testLimit = *test.limit
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

// TestEPBSResponseLimits pins two properties of every ePBS limit pair.
//
// Survivable: for a body without a declared length, readResponseBody reaches
// its length check only after io.ReadAll has buffered the whole body.  Peak
// allocation is about twice the limit from Go 1.26, and up to about 2.25 times
// on Go 1.25, whose ReadAll grows one buffer by about a quarter at a time.  A
// limit taken from the protocol's theoretical maximum (~2GiB of SSZ, ~5GiB of
// JSON hex) is therefore never reached -- the process dies first -- and a body
// admitted just under one is equally fatal.  The ceiling below is an
// allocation budget, not a protocol bound: every live preset produces a
// payload-included response orders of magnitude smaller, so the proposal JSON
// limit sitting exactly at half the budget is the deliberate ceiling rather
// than a coincidence.  Raising any limit past it is a regression.
//
// One budget for the data: JSON hex costs slightly over two bytes per SSZ
// byte, so a JSON limit below twice its SSZ limit would accept less data when
// a node chooses to answer in JSON.
func TestEPBSResponseLimits(t *testing.T) {
	const maxResponsePeakAllocation = 1024 * 1024 * 1024

	tests := []struct {
		name  string
		limit responseLimit
	}{
		{name: "Fetch", limit: epbsResponseLimit},
		{name: "Proposal", limit: epbsProposalResponseLimit},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Zero would leave the encoding unbounded.
			require.Positive(t, test.limit.ssz)
			require.Positive(t, test.limit.json)
			require.LessOrEqual(t, 2*test.limit.ssz, maxResponsePeakAllocation)
			require.LessOrEqual(t, 2*test.limit.json, maxResponsePeakAllocation)
			require.GreaterOrEqual(t, test.limit.json, 2*test.limit.ssz)
		})
	}
}
