// Copyright © 2021 - 2026 Attestant Limited.
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

package multi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/api/v1"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// proposerDutiesV1Only is a client that serves v1 proposer duties but does not
// implement the v2 provider.
type proposerDutiesV1Only struct {
	consensusclient.Service
}

func (s *proposerDutiesV1Only) ProposerDuties(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*v1.ProposerDuty],
	error,
) {
	return s.Service.(consensusclient.ProposerDutiesProvider).ProposerDuties(ctx, opts)
}

// TestProposerDutiesV2SkipsUnsupportedClientWithoutDeactivatingIt covers a client
// that does not implement the v2 provider at all.  v2 support is a static
// property of a node, so the call must move past it without evicting it: a
// deactivated client is out of the rotation for every other call until the next
// recheck, and duties are fetched every epoch, so it would flap permanently.
func TestProposerDutiesV2SkipsUnsupportedClientWithoutDeactivatingIt(t *testing.T) {
	ctx := context.Background()

	unsupportedClient, err := mock.New(ctx, mock.WithName("unsupported"))
	require.NoError(t, err)
	supportedClient, err := mock.New(ctx, mock.WithName("supported"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			&proposerDutiesV1Only{Service: unsupportedClient},
			supportedClient,
		}),
	)
	require.NoError(t, err)

	response, err := multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.NotNil(t, response)

	// Address() reports activeClients[0], so the v1-only client still leading it
	// is what says it was skipped rather than evicted.
	require.Equal(t, "unsupported", multiClient.Address())
	require.True(t, multiClient.IsActive())

	// And it is still the one serving the call it can serve.
	v1Response, err := multiClient.(consensusclient.ProposerDutiesProvider).ProposerDuties(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.NotNil(t, v1Response)
	require.Equal(t, "unsupported", multiClient.Address())
}

// TestProposerDutiesV2FailsOverFromANodeWithoutTheEndpoint covers the case that
// actually arises in production: *http.Service implements the provider whether
// or not the node serves the endpoint, so a node without v2 answers 404.  doCall
// returns 4xx to the caller without failing over, so this has to be classified
// as "cannot serve" for a mixed cluster to work at all.
func TestProposerDutiesV2FailsOverFromANodeWithoutTheEndpoint(t *testing.T) {
	ctx := context.Background()

	duties := []*v1.ProposerDuty{{ValidatorIndex: 7, Slot: 65}}

	noEndpointClient, err := mock.New(ctx, mock.WithName("no-endpoint"))
	require.NoError(t, err)
	noEndpointClient.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{
			Method:     "GET",
			Endpoint:   "/eth/v2/validator/duties/proposer/1",
			StatusCode: 404,
		}
	}
	servingClient, err := mock.New(ctx, mock.WithName("serving"))
	require.NoError(t, err)
	servingClient.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return &api.Response[[]*v1.ProposerDuty]{Data: duties, Metadata: map[string]any{}}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{noEndpointClient, servingClient}),
	)
	require.NoError(t, err)

	response, err := multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.Equal(t, duties, response.Data)

	// The 404 is not a health signal either, so that client keeps its place.
	require.Equal(t, "no-endpoint", multiClient.Address())
}

// TestProposerDutiesV2WithNoSupportingClient verifies the error a caller gets
// when nothing in the cluster can serve the call.
func TestProposerDutiesV2WithNoSupportingClient(t *testing.T) {
	ctx := context.Background()

	unsupportedClient, err := mock.New(ctx, mock.WithName("unsupported"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{&proposerDutiesV1Only{Service: unsupportedClient}}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.ErrorIs(t, err, multi.ErrCallNotSupported)
	require.True(t, multiClient.IsActive())
}

func TestProposerDutiesV2StatusClassification(t *testing.T) {
	ctx := context.Background()
	duties := []*v1.ProposerDuty{{ValidatorIndex: 7, Slot: 65}}

	tests := []struct {
		name     string
		status   int
		endpoint string
		failOver bool
	}{
		{name: "MethodNotAllowed", status: 405, endpoint: "/eth/v2/validator/duties/proposer/1", failOver: true},
		{name: "BadRequest", status: 400, endpoint: "/eth/v2/validator/duties/proposer/1", failOver: false},
		{name: "SpecNotFound", status: 404, endpoint: "/eth/v1/config/spec", failOver: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var secondCalls int
			first, err := mock.New(ctx, mock.WithName("first"))
			require.NoError(t, err)
			first.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
				return nil, &api.Error{Method: "GET", Endpoint: test.endpoint, StatusCode: test.status}
			}
			second, err := mock.New(ctx, mock.WithName("second"))
			require.NoError(t, err)
			second.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
				secondCalls++

				return &api.Response[[]*v1.ProposerDuty]{Data: duties, Metadata: map[string]any{}}, nil
			}

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]consensusclient.Service{first, second}),
			)
			require.NoError(t, err)

			response, err := multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
			if test.failOver {
				require.NoError(t, err)
				require.Equal(t, duties, response.Data)
				require.Equal(t, 1, secondCalls)
			} else {
				require.Error(t, err)
				require.NotErrorIs(t, err, multi.ErrCallNotSupported)
				var apiErr *api.Error
				require.ErrorAs(t, err, &apiErr)
				require.Equal(t, test.status, apiErr.StatusCode)
				require.Equal(t, 0, secondCalls)
			}
			require.Equal(t, "first", multiClient.Address())
			require.True(t, multiClient.IsActive())
		})
	}
}

func TestProposerDutiesV2WithNoServingEndpoint(t *testing.T) {
	ctx := context.Background()

	client, err := mock.New(ctx, mock.WithName("no-endpoint"))
	require.NoError(t, err)
	client.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{
			Method:     "GET",
			Endpoint:   "/eth/v2/validator/duties/proposer/1",
			StatusCode: 404,
			Data:       []byte(`{"message":"not found"}`),
		}
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.ErrorIs(t, err, multi.ErrCallNotSupported)
	require.NotContains(t, err.Error(), "\n")
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, 404, apiErr.StatusCode)
	require.True(t, multiClient.IsActive())
	require.Equal(t, "no-endpoint", multiClient.Address())
}

func TestProposerDutiesV2KeepsARealFailureAfterAnUnsupportedClient(t *testing.T) {
	ctx := context.Background()

	skipped, err := mock.New(ctx, mock.WithName("skipped"))
	require.NoError(t, err)
	skipped.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{Method: "GET", Endpoint: "/eth/v2/validator/duties/proposer/1", StatusCode: 404}
	}
	failed, err := mock.New(ctx, mock.WithName("failed"))
	require.NoError(t, err)
	failed.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{Method: "GET", Endpoint: "/eth/v2/validator/duties/proposer/1", StatusCode: 503}
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{skipped, failed}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NotErrorIs(t, err, multi.ErrCallNotSupported)
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, 503, apiErr.StatusCode)
	require.Equal(t, "skipped", multiClient.Address())
}

func TestProposerDutiesV2KeepsARealFailureAheadOfAnUnsupportedClient(t *testing.T) {
	ctx := context.Background()

	failed, err := mock.New(ctx, mock.WithName("failed"))
	require.NoError(t, err)
	failed.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{Method: "GET", Endpoint: "/eth/v2/validator/duties/proposer/1", StatusCode: 503}
	}
	skipped, err := mock.New(ctx, mock.WithName("skipped"))
	require.NoError(t, err)
	skipped.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return nil, &api.Error{Method: "GET", Endpoint: "/eth/v2/validator/duties/proposer/1", StatusCode: 404}
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{failed, skipped}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NotErrorIs(t, err, multi.ErrCallNotSupported)
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, 503, apiErr.StatusCode)
	require.Equal(t, "skipped", multiClient.Address())
	require.True(t, multiClient.IsActive())
}

func TestProposerDutiesV2DoesNotDeactivateAWrapperAroundAV1Client(t *testing.T) {
	ctx := context.Background()

	inner, err := mock.New(ctx, mock.WithName("v1"))
	require.NoError(t, err)
	wrapped, err := testclients.NewErroring(ctx, 0, &proposerDutiesV1Only{Service: inner})
	require.NoError(t, err)
	serving, err := mock.New(ctx, mock.WithName("serving"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{wrapped, serving}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.Equal(t, wrapped.Address(), multiClient.Address())
}

func TestProposerDutiesV2DeactivatesAnErroringClient(t *testing.T) {
	ctx := context.Background()

	inner, err := mock.New(ctx, mock.WithName("erroring"))
	require.NoError(t, err)
	erroring, err := testclients.NewErroring(ctx, 1, inner)
	require.NoError(t, err)
	serving, err := mock.New(ctx, mock.WithName("serving"))
	require.NoError(t, err)
	serving.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return &api.Response[[]*v1.ProposerDuty]{
			Data:     []*v1.ProposerDuty{{ValidatorIndex: 7, Slot: 65}},
			Metadata: map[string]any{},
		}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{erroring, serving}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.Equal(t, "serving", multiClient.Address())
}

func TestProposerDutiesV2(t *testing.T) {
	ctx := context.Background()
	root := phase0.Root{0xaa}
	duties := []*v1.ProposerDuty{{ValidatorIndex: 7, Slot: 65}}

	mockClient, err := mock.New(ctx)
	require.NoError(t, err)
	mockClient.ProposerDutiesV2Func = func(_ context.Context, _ *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return &api.Response[[]*v1.ProposerDuty]{
			Data: duties,
			Metadata: map[string]any{
				"dependent_root":       root,
				"execution_optimistic": true,
			},
		}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{mockClient}),
	)
	require.NoError(t, err)

	provider, supported := multiClient.(consensusclient.ProposerDutiesV2Provider)
	require.True(t, supported)

	response, err := provider.ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{Epoch: 2})
	require.NoError(t, err)
	require.Same(t, duties[0], response.Data[0])
	require.Equal(t, root, response.Metadata["dependent_root"])
	require.Equal(t, true, response.Metadata["execution_optimistic"])
}

func TestProposerDutiesV2NilOptionsDoesNotCallTheNode(t *testing.T) {
	ctx := context.Background()
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/eth/v1/node/version":
			fmt.Fprint(w, `{"data":{"version":"stub"}}`)
		case "/eth/v1/node/syncing":
			fmt.Fprint(w, `{"data":{"head_slot":"64","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`)
		default:
			hits++
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	first, err := eth2http.New(ctx, eth2http.WithAddress(server.URL), eth2http.WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	second, err := eth2http.New(ctx, eth2http.WithAddress(server.URL), eth2http.WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{first, second}),
	)
	require.NoError(t, err)

	hits = 0
	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, nil)
	require.ErrorIs(t, err, consensusclient.ErrNoOptions)
	require.Equal(t, 0, hits)
	require.Equal(t, first.Address(), multiClient.Address())
}

func TestProposerDutiesV2FailsOverFromTheRouteRealHTTPEmits(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/eth/v1/node/version":
			fmt.Fprint(w, `{"data":{"version":"stub"}}`)
		case "/eth/v1/node/syncing":
			fmt.Fprint(w, `{"data":{"head_slot":"64","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	missing, err := eth2http.New(ctx,
		eth2http.WithAddress(server.URL),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)
	var calls int
	serving, err := mock.New(ctx, mock.WithName("serving"))
	require.NoError(t, err)
	serving.ProposerDutiesV2Func = func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		calls++

		return &api.Response[[]*v1.ProposerDuty]{Data: []*v1.ProposerDuty{{ValidatorIndex: 7}}, Metadata: map[string]any{}}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{missing, serving}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.Equal(t, 1, calls)
	require.Equal(t, missing.Address(), multiClient.Address())
}

func TestProposerDuties(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	erroringClient1, err := testclients.NewErroring(ctx, 0.1, client1)
	require.NoError(t, err)
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)
	erroringClient2, err := testclients.NewErroring(ctx, 0.1, client2)
	require.NoError(t, err)
	client3, err := mock.New(ctx, mock.WithName("mock 3"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			erroringClient1,
			erroringClient2,
			client3,
		}),
	)
	require.NoError(t, err)

	for i := 0; i < 128; i++ {
		res, err := multiClient.(consensusclient.ProposerDutiesProvider).ProposerDuties(ctx, &api.ProposerDutiesOpts{})
		require.NoError(t, err)
		require.NotNil(t, res)
	}
	// At this point we expect mock 3 to be in active (unless probability hates us).
	require.Equal(t, "mock 3", multiClient.Address())
}
