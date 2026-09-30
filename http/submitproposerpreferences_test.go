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
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	clienthttp "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestSubmitProposerPreferencesPosts(t *testing.T) {
	tests := []struct {
		name        string
		enforceJSON bool
	}{
		{name: "SSZ"},
		{name: "JSON", enforceJSON: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			received := false
			server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				switch r.URL.Path {
				case "/eth/v1/node/version":
					_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
				case "/eth/v1/node/syncing":
					_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
				case "/eth/v1/validator/proposer_preferences":
					received = true
					require.Equal(t, nethttp.MethodPost, r.Method)
					require.Equal(t, "gloas", r.Header.Get("Eth-Consensus-Version"))
					if test.enforceJSON {
						require.Equal(t, "application/json", r.Header.Get("Content-Type"))
						var body []json.RawMessage
						require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
						require.Len(t, body, 1)
					} else {
						require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
						body, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						require.Len(t, body, 172)
						decoded := &gloas.SignedProposerPreferences{}
						require.NoError(t, decoded.UnmarshalSSZ(body))
					}
					w.WriteHeader(nethttp.StatusOK)
				default:
					w.WriteHeader(nethttp.StatusNotFound)
				}
			}))
			defer server.Close()

			params := []clienthttp.Parameter{clienthttp.WithAddress(server.URL)}
			if test.enforceJSON {
				params = append(params, clienthttp.WithEnforceJSON(true))
			}
			service, err := clienthttp.New(ctx, params...)
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(ctx, preferencesOpts(onePreference()))
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

func TestSubmitProposerPreferencesReportsTransportError(t *testing.T) {
	ctx := context.Background()
	received := false
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/validator/proposer_preferences":
			received = true
			require.Equal(t, nethttp.MethodPost, r.Method)
			require.Equal(t, "gloas", r.Header.Get("Eth-Consensus-Version"))
			w.WriteHeader(nethttp.StatusBadRequest)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, preferencesOpts(onePreference()))
	require.ErrorContains(t, err, "failed to submit proposer preferences")
	require.True(t, received)
}

func TestSubmitProposerPreferencesRequiresOptions(t *testing.T) {
	received := false
	server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
	defer server.Close()

	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).SubmitProposerPreferences(context.Background(), nil)
	require.ErrorIs(t, err, client.ErrNoOptions)
	require.False(t, received)
}

func TestSubmitProposerPreferencesRequiresExactStatusOK(t *testing.T) {
	ctx := context.Background()
	server := proposerPreferencesServer(t, nethttp.StatusNoContent, nil)
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, preferencesOpts(onePreference()))
	require.EqualError(t, err, "failed to submit proposer preferences\nunexpected status code 204")
}

// TestSubmitProposerPreferencesEnforcesStaticLimit covers a node that does not
// publish the spec at all: the mainnet preset stands in for the chain's own
// lookahead length rather than the submission failing.
func TestSubmitProposerPreferencesEnforcesStaticLimit(t *testing.T) {
	for _, test := range []struct {
		name  string
		count int
		err   string
	}{
		{name: "AtLimit", count: 64},
		{name: "OneOverLimit", count: 65, err: "too many proposer preferences"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
			defer server.Close()
			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(context.Background(), preferencesOpts(makePreferences(test.count)))
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				require.False(t, received)
			} else {
				require.NoError(t, err)
				require.True(t, received)
			}
		})
	}
}

func TestSubmitProposerPreferencesChecksLimitBeforeNilElements(t *testing.T) {
	received := false
	server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
	defer server.Close()

	preferences := makePreferences(65)
	preferences[0] = nil
	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(context.Background(), preferencesOpts(preferences))
	require.ErrorContains(t, err, "too many proposer preferences")
	require.ErrorIs(t, err, client.ErrInvalidOptions)
	require.False(t, received)
}

// TestSubmitProposerPreferencesUsesSpecDerivedLimit covers a minimal-preset chain
// whose lookahead is shorter than mainnet's.  The limit comes from the spec
// without WithCustomSpecSupport, which selects SSZ codecs rather than governing
// whether the spec is available.
func TestSubmitProposerPreferencesUsesSpecDerivedLimit(t *testing.T) {
	for _, test := range []struct {
		name  string
		count int
		err   string
	}{
		{name: "AtLimit", count: 4},
		{name: "OneOverLimit", count: 5, err: "too many proposer preferences"},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := proposerPreferencesServerWithSpec(t, nethttp.StatusOK, &received, 1, 2)
			defer server.Close()
			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(context.Background(), preferencesOpts(makePreferences(test.count)))
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				require.False(t, received)
			} else {
				require.NoError(t, err)
				require.True(t, received)
			}
		})
	}
}

func TestSubmitProposerPreferencesRejectsNilElement(t *testing.T) {
	for _, test := range []struct {
		name        string
		enforceJSON bool
	}{
		{name: "SSZ"},
		{name: "JSON", enforceJSON: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
			defer server.Close()

			params := []clienthttp.Parameter{clienthttp.WithAddress(server.URL)}
			if test.enforceJSON {
				params = append(params, clienthttp.WithEnforceJSON(true))
			}
			service, err := clienthttp.New(context.Background(), params...)
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(context.Background(), preferencesOpts([]*gloas.SignedProposerPreferences{nil}))
			require.ErrorContains(t, err, "nil proposer preference supplied")
			require.ErrorIs(t, err, client.ErrInvalidOptions)
			require.False(t, received)
		})
	}
}

func TestSubmitProposerPreferencesRejectsNilMessage(t *testing.T) {
	received := false
	server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
	defer server.Close()

	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(context.Background(), preferencesOpts([]*gloas.SignedProposerPreferences{{}}))
	require.ErrorContains(t, err, "nil proposer preference message supplied")
	require.ErrorIs(t, err, client.ErrInvalidOptions)
	require.False(t, received)
}

// TestSubmitProposerPreferencesFallsBackWhenSpecIncomplete covers a node that
// publishes the spec but omits one of the two keys the limit is derived from.
func TestSubmitProposerPreferencesFallsBackWhenSpecIncomplete(t *testing.T) {
	received := false
	server := preferencesServer(t, nethttp.StatusOK, &received, `{"data":{"SLOTS_PER_EPOCH":"32"}}`)
	defer server.Close()

	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(context.Background(), preferencesOpts(makePreferences(64)))
	require.NoError(t, err)
	require.True(t, received)
}

func TestSubmitProposerPreferencesAllowsEmptyList(t *testing.T) {
	for _, test := range []struct {
		name        string
		enforceJSON bool
	}{
		{name: "SSZ"},
		{name: "JSON", enforceJSON: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := emptyPreferencesServer(t, test.enforceJSON, &received)
			defer server.Close()

			params := []clienthttp.Parameter{clienthttp.WithAddress(server.URL)}
			if test.enforceJSON {
				params = append(params, clienthttp.WithEnforceJSON(true))
			}
			service, err := clienthttp.New(context.Background(), params...)
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(context.Background(), &api.SubmitProposerPreferencesOpts{})
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

func preferencesOpts(preferences []*gloas.SignedProposerPreferences) *api.SubmitProposerPreferencesOpts {
	return &api.SubmitProposerPreferencesOpts{Preferences: preferences}
}

func onePreference() []*gloas.SignedProposerPreferences {
	return []*gloas.SignedProposerPreferences{{Message: &gloas.ProposerPreferences{}}}
}

func makePreferences(count int) []*gloas.SignedProposerPreferences {
	preferences := make([]*gloas.SignedProposerPreferences, count)
	for i := range preferences {
		preferences[i] = onePreference()[0]
	}

	return preferences
}

// proposerPreferencesServer serves a node that does not publish the spec, so the
// static limit applies.
func proposerPreferencesServer(t *testing.T, status int, received *bool) *httptest.Server {
	t.Helper()

	return preferencesServer(t, status, received, "")
}

// proposerPreferencesServerWithSpec serves a node publishing both of the keys the
// limit is derived from.
func proposerPreferencesServerWithSpec(t *testing.T,
	status int,
	received *bool,
	minSeedLookahead, slotsPerEpoch uint64,
) *httptest.Server {
	t.Helper()

	return preferencesServer(t, status, received,
		fmt.Sprintf(`{"data":{"MIN_SEED_LOOKAHEAD":"%d","SLOTS_PER_EPOCH":"%d"}}`, minSeedLookahead, slotsPerEpoch))
}

// preferencesServer serves the submission endpoint with the given status.  An
// empty specResponse serves a node with no spec endpoint at all.
func preferencesServer(t *testing.T, status int, received *bool, specResponse string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/config/spec":
			if specResponse == "" {
				w.WriteHeader(nethttp.StatusNotFound)

				return
			}
			_, _ = w.Write([]byte(specResponse))
		case "/eth/v1/validator/proposer_preferences":
			if received != nil {
				*received = true
			}
			w.WriteHeader(status)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
}

func emptyPreferencesServer(t *testing.T, enforceJSON bool, received *bool) *httptest.Server {
	t.Helper()

	return httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/validator/proposer_preferences":
			*received = true
			require.Equal(t, nethttp.MethodPost, r.Method)
			if enforceJSON {
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, "[]", string(body))
				var decoded []json.RawMessage
				require.NoError(t, json.Unmarshal(body, &decoded))
				require.Empty(t, decoded)
			} else {
				require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Empty(t, body)
			}
			w.WriteHeader(nethttp.StatusOK)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
}
