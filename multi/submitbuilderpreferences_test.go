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

package multi_test

import (
	"context"
	"errors"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// invalidOptionsBuilderPreferencesClient rejects every submission with a
// caller error: ErrNoOptions for nil options, else ErrInvalidOptions.
type invalidOptionsBuilderPreferencesClient struct {
	name  string
	calls int
}

func (c *invalidOptionsBuilderPreferencesClient) Name() string    { return c.name }
func (c *invalidOptionsBuilderPreferencesClient) Address() string { return c.name }
func (*invalidOptionsBuilderPreferencesClient) IsActive() bool    { return true }
func (*invalidOptionsBuilderPreferencesClient) IsSynced() bool    { return true }

func (c *invalidOptionsBuilderPreferencesClient) SubmitBuilderPreferences(_ context.Context,
	opts *api.SubmitBuilderPreferencesOpts,
) error {
	c.calls++
	if opts == nil {
		return consensusclient.ErrNoOptions
	}

	return errors.Join(errors.New("nil builder preferences entry supplied"), consensusclient.ErrInvalidOptions)
}

func TestSubmitBuilderPreferences(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client1, client2}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.BuilderPreferencesSubmitter).
		SubmitBuilderPreferences(ctx, &api.SubmitBuilderPreferencesOpts{})
	require.NoError(t, err)
}

// TestSubmitBuilderPreferencesDoesNotFailOverOnInvalidOptions ensures that a
// caller-side error neither reaches a second client nor deactivates the first.
func TestSubmitBuilderPreferencesDoesNotFailOverOnInvalidOptions(t *testing.T) {
	for _, test := range []struct {
		name string
		opts *api.SubmitBuilderPreferencesOpts
		err  error
	}{
		{name: "InvalidOptions", opts: &api.SubmitBuilderPreferencesOpts{}, err: consensusclient.ErrInvalidOptions},
		{name: "NoOptions", err: consensusclient.ErrNoOptions},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			client1 := &invalidOptionsBuilderPreferencesClient{name: "invalid 1"}
			client2 := &invalidOptionsBuilderPreferencesClient{name: "invalid 2"}

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]consensusclient.Service{client1, client2}),
			)
			require.NoError(t, err)

			err = multiClient.(consensusclient.BuilderPreferencesSubmitter).SubmitBuilderPreferences(ctx, test.opts)
			require.ErrorIs(t, err, test.err)
			require.Equal(t, 1, client1.calls)
			require.Equal(t, 0, client2.calls)
			require.Equal(t, "invalid 1", multiClient.Address())
		})
	}
}
