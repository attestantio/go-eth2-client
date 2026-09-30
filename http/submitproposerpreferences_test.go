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

// TestSubmitProposerPreferencesSkipsLimitWithoutSpec covers a node that does not
// publish the spec at all.  The bound is a client-side sanity check and the node
// is the authority, so a chain whose own lookahead cannot be derived gets no
// client-side bound rather than the mainnet preset: substituting 64 fails open
// on a shorter-lookahead chain, but rejects a submission a longer-lookahead
// chain would have accepted, without ever reaching the node.
func TestSubmitProposerPreferencesSkipsLimitWithoutSpec(t *testing.T) {
	for _, test := range []struct {
		name  string
		count int
	}{
		{name: "AtMainnetLimit", count: 64},
		{name: "OverMainnetLimit", count: 65},
		{name: "OverALongerChainsLimit", count: 129},
	} {
		t.Run(test.name, func(t *testing.T) {
			received := false
			server := proposerPreferencesServer(t, nethttp.StatusOK, &received)
			defer server.Close()
			service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
			require.NoError(t, err)
			err = service.(client.ProposerPreferencesSubmitter).
				SubmitProposerPreferences(context.Background(), preferencesOpts(makePreferences(test.count)))
			require.NoError(t, err)
			require.True(t, received)
		})
	}
}

func TestSubmitProposerPreferencesChecksLimitBeforeNilElements(t *testing.T) {
	received := false
	// A chain whose limit is derivable, so that there is a limit to hit first.
	server := proposerPreferencesServerWithSpec(t, nethttp.StatusOK, &received, 1, 2)
	defer server.Close()

	preferences := makePreferences(5)
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

// TestSubmitProposerPreferencesSkipsLimitWhenSpecIncomplete covers a node that
// publishes the spec but omits one of the two keys the limit is derived from.
// The chain's own lookahead is no more knowable than for a node serving no spec
// at all, so there is no client-side bound: a count that mainnet would refuse
// still reaches the node, which is the authority on it.
func TestSubmitProposerPreferencesSkipsLimitWhenSpecIncomplete(t *testing.T) {
	received := false
	server := preferencesServer(t, nethttp.StatusOK, &received, `{"data":{"SLOTS_PER_EPOCH":"32"}}`)
	defer server.Close()

	service, err := clienthttp.New(context.Background(), clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	err = service.(client.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(context.Background(), preferencesOpts(makePreferences(65)))
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

// countingPreferencesServer serves the submission endpoint and counts requests
// to the spec endpoint, so that what reaches the network on a submission can be
// asserted rather than inferred.
func countingPreferencesServer(t *testing.T, specResponse string, specRequests *int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		case "/eth/v1/config/spec":
			*specRequests++
			if specResponse == "" {
				w.WriteHeader(nethttp.StatusNotFound)

				return
			}
			_, _ = w.Write([]byte(specResponse))
		case "/eth/v1/validator/proposer_preferences":
			w.WriteHeader(nethttp.StatusOK)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
}

// TestSubmitProposerPreferencesDerivesTheLimitOnce pins the limit being cached.
//
// Spec() assigns s.spec only on success and clearStaticValues nils it every 5
// minutes, while preferences are submitted about once per epoch -- so consulting
// it per submission put a blocking spec round trip on the publish path roughly
// every time, holding specMutex across it and stalling every other spec consumer
// with it.  The lookahead is a property of the chain, so one derivation is
// enough for the life of the service.
func TestSubmitProposerPreferencesDerivesTheLimitOnce(t *testing.T) {
	ctx := context.Background()
	specRequests := 0
	server := countingPreferencesServer(t,
		`{"data":{"MIN_SEED_LOOKAHEAD":"1","SLOTS_PER_EPOCH":"2"}}`, &specRequests)
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	submitter := service.(client.ProposerPreferencesSubmitter)

	// Whatever the connection setup asked for, no submission adds to it after
	// the first that has to derive the limit.
	for range 5 {
		require.NoError(t, submitter.SubmitProposerPreferences(ctx, preferencesOpts(makePreferences(1))))
	}
	afterFirst := specRequests

	for range 5 {
		require.NoError(t, submitter.SubmitProposerPreferences(ctx, preferencesOpts(makePreferences(1))))
	}
	require.Equal(t, afterFirst, specRequests)

	// And the cached limit is still enforced.
	err = submitter.SubmitProposerPreferences(ctx, preferencesOpts(makePreferences(5)))
	require.ErrorContains(t, err, "too many proposer preferences")
}

// TestSubmitProposerPreferencesSkipsTheLimitForAnEmptyList covers the common
// "nothing to publish" call.  Go evaluates the right-hand operand of > without
// regard to the left, so an unguarded comparison made an empty submission pay
// the spec lookup to check a bound it cannot exceed.
func TestSubmitProposerPreferencesSkipsTheLimitForAnEmptyList(t *testing.T) {
	ctx := context.Background()
	specRequests := 0
	// A node with no spec endpoint, so each lookup is a failed round trip that
	// is not cached and is repeated on every submission.
	server := countingPreferencesServer(t, "", &specRequests)
	defer server.Close()

	service, err := clienthttp.New(ctx, clienthttp.WithAddress(server.URL))
	require.NoError(t, err)
	submitter := service.(client.ProposerPreferencesSubmitter)

	before := specRequests
	for range 5 {
		require.NoError(t, submitter.SubmitProposerPreferences(ctx, preferencesOpts(nil)))
	}
	require.Equal(t, before, specRequests)
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
