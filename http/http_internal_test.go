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
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rs/zerolog"
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

func TestParseBasicAuth(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		headers    map[string]string
		expHeaders map[string]string
	}{
		{
			name:       "No Scheme",
			address:    "127.0.0.1:5051",
			expHeaders: nil,
		},
		{
			name:       "With Scheme",
			address:    "http://127.0.0.1:5051",
			expHeaders: nil,
		},
		{
			name:       "Simple",
			address:    "http://user:pass@foo.com",
			expHeaders: map[string]string{"Authorization": "Basic dXNlcjpwYXNz"},
		},
		{
			name:       "Missing user",
			address:    "http://:pass@foo.com",
			expHeaders: map[string]string{"Authorization": "Basic OnBhc3M="},
		},
		{
			name:       "Missing pass",
			address:    "http://user:@foo.com",
			expHeaders: map[string]string{"Authorization": "Basic dXNlcjo="},
		},
		{
			name:       "Missing user and pass",
			address:    "http://:@foo.com",
			expHeaders: map[string]string{"Authorization": "Basic Og=="},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			base, _, err := parseAddress(test.address)
			require.NoError(t, err)
			newHeaders := parseBasicAuth(base, test.headers)
			require.NoError(t, err)
			require.Equal(t, test.expHeaders, newHeaders)
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

func TestProxy(t *testing.T) {
	tests := []struct {
		name           string
		baseURL        string
		extraHeaders   map[string]string
		backendHandler http.HandlerFunc
		requestPath    string
		requestMethod  string
		requestBody    string
		requestHeaders map[string]string
		expectedStatus int
		expectedBody   string
		expectedHeader string
		expectError    bool
	}{
		{
			name:    "SimpleGET",
			baseURL: "", // Will be set to test server URL
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"result":"success"}`))
			},
			requestPath:    "/api/v1/test",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"result":"success"}`,
		},
		{
			name:    "POST",
			baseURL: "",
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"received":"` + string(body) + `"}`))
			},
			requestPath:    "/api/v1/create",
			requestMethod:  "POST",
			requestBody:    `{"data":"test"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"received":"{"data":"test"}"}`,
		},
		{
			name:    "CustomHeaders",
			baseURL: "",
			extraHeaders: map[string]string{
				"X-Custom-Header": "custom-value",
			},
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				// Verify custom header was forwarded
				if r.Header.Get("X-Custom-Header") != "custom-value" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("headers-ok"))
			},
			requestPath:    "/api/v1/headers",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "headers-ok",
		},
		{
			name:    "ResponseHeaders",
			baseURL: "",
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Response-Header", "response-value")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			requestPath:    "/api/v1/response",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
			expectedHeader: "response-value",
		},
		{
			name:    "ErrorResponse",
			baseURL: "",
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":"not found"}`))
			},
			requestPath:    "/api/v1/notfound",
			requestMethod:  "GET",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		},
		{
			name:    "BasicAuth",
			baseURL: "http://testuser:testpass@", // Will append server host
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				user, pass, ok := r.BasicAuth()
				if !ok || user != "testuser" || pass != "testpass" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("auth-ok"))
			},
			requestPath:    "/api/v1/auth",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "auth-ok",
		},
		{
			name:    "RequestHeaders",
			baseURL: "",
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Request-Header") != "request-value" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("request-headers-ok"))
			},
			requestPath:   "/api/v1/reqheaders",
			requestMethod: "GET",
			requestHeaders: map[string]string{
				"X-Request-Header": "request-value",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "request-headers-ok",
		},
		{
			name:    "LargeBody",
			baseURL: "",
			backendHandler: func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(body)
			},
			requestPath:    "/api/v1/large",
			requestMethod:  "POST",
			requestBody:    string(make([]byte, 10000)), // 10KB of zeros
			expectedStatus: http.StatusOK,
			expectedBody:   string(make([]byte, 10000)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create test backend server
			backend := httptest.NewServer(test.backendHandler)
			defer backend.Close()

			// Setup base URL
			baseURL := test.baseURL
			if baseURL == "" {
				baseURL = backend.URL
			} else if baseURL == "http://testuser:testpass@" {
				// Extract host from backend URL and add auth
				backendURL, _ := url.Parse(backend.URL)
				baseURL = "http://testuser:testpass@" + backendURL.Host
			}

			parsedBase, err := url.Parse(baseURL)
			require.NoError(t, err)

			// Create service
			service := &Service{
				log:          zerolog.Nop(),
				base:         parsedBase,
				address:      baseURL,
				client:       backend.Client(),
				extraHeaders: test.extraHeaders,
			}

			// Create request
			var body io.Reader
			if test.requestBody != "" {
				body = bytes.NewBufferString(test.requestBody)
			}
			req, err := http.NewRequest(test.requestMethod, test.requestPath, body)
			require.NoError(t, err)

			// Add request headers
			for k, v := range test.requestHeaders {
				req.Header.Set(k, v)
			}

			// Call proxy
			ctx := context.Background()
			resp, err := service.proxy(ctx, req)

			// Check error expectation
			if test.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Check status code
			require.Equal(t, test.expectedStatus, resp.StatusCode)

			// Check body
			if test.expectedBody != "" {
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Equal(t, test.expectedBody, string(bodyBytes))
				_ = resp.Body.Close()
			}

			// Check response header if specified
			if test.expectedHeader != "" {
				require.Equal(t, test.expectedHeader, resp.Header.Get("X-Response-Header"))
			}
		})
	}
}
