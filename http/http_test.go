package http_test

import (
	"context"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestError(t *testing.T) {
	status := nethttp.StatusTeapot
	data := []byte("data")
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(status)
		_, _ = w.Write(data)
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := http.New(ctx, http.WithAddress(srv.URL))

	require.NotNil(t, err)
	require.Equal(t, "failed to confirm node connection: failed to fetch genesis: failed to request genesis: GET failed with status 418: data", err.Error())

	var httpError http.Error
	require.True(t, errors.As(err, &httpError))
	require.Equal(t, status, httpError.StatusCode)
	require.Equal(t, data, httpError.Data)
	require.Equal(t, nethttp.MethodGet, httpError.Method)
	require.Equal(t, "/eth/v1/beacon/genesis", httpError.Endpoint)
}
