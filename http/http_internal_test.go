// Copyright © 2024 Attestant Limited.
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
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		base    string
		address string
		err     string
	}{
		{
			name:    "Simple",
			input:   "foo",
			base:    "http://foo",
			address: "http://foo",
		},
		{
			name:    "Port",
			input:   "foo:12345",
			base:    "http://foo:12345",
			address: "http://foo:12345",
		},
		{
			name:    "Scheme",
			input:   "https://foo",
			base:    "https://foo",
			address: "https://foo",
		},
		{
			name:    "Query",
			input:   "http://foo.com?a=1&b=2",
			base:    "http://foo.com?a=1&b=2",
			address: "http://foo.com?a=xxxxx&b=xxxxx",
		},
		{
			name:    "User",
			input:   "http://user@foo.com?a=1&b=2",
			base:    "http://user@foo.com?a=1&b=2",
			address: "http://user@foo.com?a=xxxxx&b=xxxxx",
		},
		{
			name:    "Pass",
			input:   "http://user:pass@foo.com?a=1&b=2",
			base:    "http://user:pass@foo.com?a=1&b=2",
			address: "http://user:xxxxx@foo.com?a=xxxxx&b=xxxxx",
		},
		{
			name:    "Path",
			input:   "http://user:pass@foo.com/path?a=1&b=2",
			base:    "http://user:pass@foo.com/path?a=1&b=2",
			address: "http://user:xxxxx@foo.com/xxxxx?a=xxxxx&b=xxxxx",
		},
		{
			name:    "PathTrailingSlash",
			input:   "http://user:pass@foo.com/path/?a=1&b=2",
			base:    "http://user:pass@foo.com/path?a=1&b=2",
			address: "http://user:xxxxx@foo.com/xxxxx?a=xxxxx&b=xxxxx",
		},
		{
			name:  "Invalid",
			input: "http:// foo",
			err:   "invalid URL\nparse \"http:// foo\": invalid character \" \" in host name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			base, address, err := parseAddress(test.input)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.Equal(t, test.base, base.String())
				require.Equal(t, test.address, address.String())
			}
		})
	}
}

func mustParseURL(input string) *url.URL {
	base, _, err := parseAddress(input)
	if err != nil {
		panic(err)
	}

	return base
}

func TestURLForCall(t *testing.T) {
	tests := []struct {
		name     string
		base     *url.URL
		endpoint string
		query    string
		expected string
	}{
		{
			name:     "Simple",
			base:     mustParseURL("http://example.com"),
			endpoint: "/foo",
			expected: "http://example.com/foo",
		},
		{
			name:     "WithQuery",
			base:     mustParseURL("http://example.com"),
			endpoint: "/foo",
			query:    "bar=3",
			expected: "http://example.com/foo?bar=3",
		},
		{
			name:     "WithBaseQuery",
			base:     mustParseURL("http://example.com/?bar=3"),
			endpoint: "/foo",
			expected: "http://example.com/foo?bar=3",
		},
		{
			name:     "Complex",
			base:     mustParseURL("http://user:pass@foo.com/path?a=1&b=2"),
			endpoint: "/foo",
			query:    "bar=3",
			expected: "http://user:pass@foo.com/path/foo?a=1&b=2&bar=3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := urlForCall(test.base, test.endpoint, test.query)
			require.Equal(t, test.expected, url.String())
		})
	}
}
