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

package http_test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	clienthttp "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSubmitBuilderPreferencesPosts(t *testing.T) {
	entries := []*gloas.BuilderPreferencesEntry{
		builderPreferencesEntry(0x01, "https://builder-a.example.com", 1000000000),
		builderPreferencesEntry(0x02, "https://builder-b.example.com/path", 7),
	}

	tests := []struct {
		name   string
		params []clienthttp.Parameter
		json   bool
	}{
		{name: "SSZ"},
		{name: "CustomSpecSSZ", params: []clienthttp.Parameter{clienthttp.WithCustomSpecSupport(true)}},
		{name: "JSON", params: []clienthttp.Parameter{clienthttp.WithEnforceJSON(true)}, json: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, r *nethttp.Request) {
				received = true
				require.Equal(t, nethttp.MethodPost, r.Method)
				require.Equal(t, "gloas", r.Header.Get("Eth-Consensus-Version"))
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				var decoded []*gloas.BuilderPreferencesEntry
				if test.json {
					require.Equal(t, "application/json", r.Header.Get("Content-Type"))
					require.NoError(t, json.Unmarshal(body, &decoded))
				} else {
					require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
					decoded = decodeBuilderPreferencesSSZ(t, body)
				}
				require.Equal(t, entries, decoded)
				w.WriteHeader(nethttp.StatusOK)
			})
			defer server.Close()

			service, err := clienthttp.New(context.Background(), append(test.params, clienthttp.WithAddress(server.URL))...)
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).SubmitBuilderPreferences(context.Background(),
				&api.SubmitBuilderPreferencesOpts{Entries: entries})
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

func TestSubmitBuilderPreferencesAllowsEmptyList(t *testing.T) {
	for _, test := range []struct {
		name   string
		params []clienthttp.Parameter
		body   string
	}{
		{name: "SSZ", body: ""},
		{name: "JSON", params: []clienthttp.Parameter{clienthttp.WithEnforceJSON(true)}, body: "[]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, r *nethttp.Request) {
				received = true
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, test.body, string(body))
				w.WriteHeader(nethttp.StatusOK)
			})
			defer server.Close()

			service, err := clienthttp.New(context.Background(), append(test.params, clienthttp.WithAddress(server.URL))...)
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).SubmitBuilderPreferences(context.Background(),
				&api.SubmitBuilderPreferencesOpts{})
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

// TestSubmitBuilderPreferencesReportsRejectedEntries covers a node that reports,
// by index, the entries it failed to submit.  The
// bodies are shaped as Prysm, Lighthouse and Teku answered on 2026-10-20; Teku
// sends the index as a string rather than the spec's number.
func TestSubmitBuilderPreferencesReportsRejectedEntries(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
	}{
		{name: "Prysm", body: `{"message":"One or more messages failed validation","code":400,"failures":[{"index":1,"message":"no such host"}]}`},
		{name: "Lighthouse", body: `{"code":400,"message":"BAD_REQUEST: error submitting builder preferences","failures":[{"index":1,"message":"no such host"}]}`},
		{name: "Teku", body: `{"code":400,"message":"Some items failed to publish, refer to errors for details","failures":[{"index":"1","message":"no such host"}]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, _ *nethttp.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(nethttp.StatusBadRequest)
				_, _ = w.Write([]byte(test.body))
			})
			defer server.Close()

			err := submitTwoBuilderPreferences(t, server.URL)
			require.ErrorContains(t, err, "failed to submit builder preferences")

			var indexedErr *api.IndexedError
			require.ErrorAs(t, err, &indexedErr)
			require.Equal(t, []api.IndexedFailure{{Index: 1, Message: "no such host"}}, indexedErr.Failures)
			var apiErr *api.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, nethttp.StatusBadRequest, apiErr.StatusCode)
		})
	}
}

// TestSubmitBuilderPreferencesReportsUnindexedErrors covers failures that are not
// per-entry rejections, such as a node that does not serve the endpoint, which
// must not look like a rejected entry.
func TestSubmitBuilderPreferencesReportsUnindexedErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "NotFound", status: nethttp.StatusNotFound},
		{name: "BadRequestWithoutFailures", status: nethttp.StatusBadRequest, body: `{"code":400,"message":"No data submitted"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, _ *nethttp.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			})
			defer server.Close()

			err := submitTwoBuilderPreferences(t, server.URL)
			var apiErr *api.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, test.status, apiErr.StatusCode)
			var indexedErr *api.IndexedError
			require.False(t, errors.As(err, &indexedErr))
		})
	}
}

func TestSubmitBuilderPreferencesRequiresExactStatusOK(t *testing.T) {
	server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNoContent)
	})
	defer server.Close()

	err := submitTwoBuilderPreferences(t, server.URL)
	require.EqualError(t, err, "failed to submit builder preferences\nunexpected status code 204")
}

func TestSubmitBuilderPreferencesRejectsInvalidOptions(t *testing.T) {
	noAuth := builderPreferencesEntry(0x01, "https://builder.example.com", 1)
	noAuth.Auth = nil
	noAuthMessage := builderPreferencesEntry(0x01, "https://builder.example.com", 1)
	noAuthMessage.Auth.Message = nil
	oversizedAuthData := builderPreferencesEntry(0x01, "https://builder.example.com", 1)
	oversizedAuthData.Auth.Message.Data = make([]byte, 4097)
	emptyAuthData := builderPreferencesEntry(0x01, "https://builder.example.com", 1)
	emptyAuthData.Auth.Message.Data = nil

	for _, test := range []struct {
		name string
		opts *api.SubmitBuilderPreferencesOpts
		err  error
		msg  string
	}{
		{name: "NilOpts", err: client.ErrNoOptions},
		{
			name: "NilEntry",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{nil}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 is nil",
		},
		{
			name: "EmptyURL",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{builderPreferencesEntry(0x01, "", 1)}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 URL length 0 is invalid",
		},
		{
			name: "OversizedURL",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{builderPreferencesEntry(0x01, "https://"+strings.Repeat("a", 2041), 1)}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 URL length 2049 is invalid",
		},
		{
			name: "OversizedAuthData",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{oversizedAuthData}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 authorization data length 4097 is invalid",
		},
		{
			name: "EmptyAuthData",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{emptyAuthData}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 authorization data length 0 is invalid",
		},
		{
			name: "NilAuth",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{noAuth}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 has no authorization",
		},
		{
			name: "NilAuthMessage",
			opts: &api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{noAuthMessage}},
			err:  client.ErrInvalidOptions,
			msg:  "builder preferences entry 0 has no authorization",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, _ *nethttp.Request) {
				received = true
				w.WriteHeader(nethttp.StatusOK)
			})
			defer server.Close()

			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).SubmitBuilderPreferences(context.Background(), test.opts)
			require.ErrorIs(t, err, test.err)
			if test.msg != "" {
				require.ErrorContains(t, err, test.msg)
			}
			require.False(t, received)
		})
	}
}

// TestSubmitBuilderPreferencesEnforcesLimit covers the limit of
// MAX_BUILDER_ENTRIES * (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH entries, with
// the node's spec giving 64 * (1 + 1) * 2 = 256, not the mainnet 4096.
func TestSubmitBuilderPreferencesEnforcesLimit(t *testing.T) {
	for _, test := range []struct {
		name  string
		count int
		err   string
	}{
		{name: "AtLimit", count: 256},
		{name: "OneOverLimit", count: 257, err: "too many builder preferences entries"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := builderPreferencesServer(t, func(w nethttp.ResponseWriter, _ *nethttp.Request) {
				received = true
				w.WriteHeader(nethttp.StatusOK)
			})
			defer server.Close()

			entries := make([]*gloas.BuilderPreferencesEntry, test.count)
			for i := range entries {
				entries[i] = builderPreferencesEntry(0x01, "https://builder.example.com", 1)
			}
			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.BuilderPreferencesSubmitter).SubmitBuilderPreferences(context.Background(),
				&api.SubmitBuilderPreferencesOpts{Entries: entries})
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				require.ErrorIs(t, err, client.ErrInvalidOptions)
				require.False(t, received)
			} else {
				require.NoError(t, err)
				require.True(t, received)
			}
		})
	}
}

func builderPreferencesEntry(pubkey byte, url string, maxExecutionPayment phase0.Gwei) *gloas.BuilderPreferencesEntry {
	return &gloas.BuilderPreferencesEntry{
		ProposerPubkey: phase0.BLSPubKey{pubkey},
		URL:            []byte(url),
		Auth: &gloas.SignedBuilderRequestAuth{
			Message:   &gloas.BuilderRequestAuth{Data: []byte("builder.example.com"), Slot: 123},
			Signature: phase0.BLSSignature{0x02},
		},
		MaxExecutionPayment: maxExecutionPayment,
	}
}

func submitTwoBuilderPreferences(t *testing.T, address string) error {
	t.Helper()

	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(address))
	require.NoError(t, err)

	return service.(client.BuilderPreferencesSubmitter).SubmitBuilderPreferences(context.Background(),
		&api.SubmitBuilderPreferencesOpts{Entries: []*gloas.BuilderPreferencesEntry{
			builderPreferencesEntry(0x01, "https://builder-a.example.com", 1),
			builderPreferencesEntry(0x01, "https://builder-b.example.com", 1),
		}})
}

// decodeBuilderPreferencesSSZ decodes an SSZ list of variable-size elements: a
// 4-byte offset per element, then the elements.
func decodeBuilderPreferencesSSZ(t *testing.T, body []byte) []*gloas.BuilderPreferencesEntry {
	t.Helper()

	require.GreaterOrEqual(t, len(body), 4)
	count := int(binary.LittleEndian.Uint32(body[0:4])) / 4
	offsets := make([]int, count+1)
	for i := range count {
		offsets[i] = int(binary.LittleEndian.Uint32(body[i*4 : i*4+4]))
	}
	offsets[count] = len(body)

	entries := make([]*gloas.BuilderPreferencesEntry, count)
	for i := range count {
		entries[i] = &gloas.BuilderPreferencesEntry{}
		require.NoError(t, entries[i].UnmarshalSSZ(body[offsets[i]:offsets[i+1]]))
	}

	return entries
}

func builderPreferencesServer(t *testing.T, handler nethttp.HandlerFunc) *httptest.Server {
	t.Helper()

	return httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/config/spec":
			_, _ = w.Write([]byte(`{"data":{"MIN_SEED_LOOKAHEAD":"1","SLOTS_PER_EPOCH":"2"}}`))
		case "/eth/v1/validator/builder_preferences":
			handler(w, r)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
}
