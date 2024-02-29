package http_test

import (
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	defaultStatus := nethttp.StatusTeapot
	defaultResponse := []byte("data")
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"8504736","sync_distance":"0"}}`))
		default:
			w.WriteHeader(nethttp.StatusTeapot)
			_, _ = w.Write([]byte("data"))
		}
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := http.New(ctx, http.WithAddress(srv.URL))
	require.NoError(t, err)

	_, err = svc.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.EqualError(t, err, "failed to request genesis\nGET failed with status 418: data")

	var apiError *api.Error
	require.True(t, errors.As(err, &apiError))
	require.Equal(t, defaultStatus, apiError.StatusCode)
	require.Equal(t, defaultResponse, apiError.Data)
	require.Equal(t, nethttp.MethodGet, apiError.Method)
	require.Equal(t, "/eth/v1/beacon/genesis", apiError.Endpoint)
}

func TestClientShouldSendExtraHeadersWhenProvided(t *testing.T) {
	authorizationHeader := "Authorization"
	authorizationToken := "Bearer token"

	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Header.Get(authorizationHeader) != authorizationToken {
			w.WriteHeader(nethttp.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/eth/v1/node/version":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"8504736","sync_distance":"0"}}`))
		default:
			w.WriteHeader(nethttp.StatusTeapot)
			_, _ = w.Write([]byte("data"))
		}
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := http.New(ctx,
		http.WithAddress(srv.URL),
		http.WithExtraHeaders(map[string]string{authorizationHeader: authorizationToken}),
	)
	require.NoError(t, err)

	_, err = svc.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.Error(t, err)
	var apiError *api.Error
	require.True(t, errors.As(err, &apiError))
	require.Equal(t, nethttp.StatusTeapot, apiError.StatusCode)
}
