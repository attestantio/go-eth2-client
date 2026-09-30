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
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// invalidOptionsPreferencesClient rejects every submission with ErrInvalidOptions,
// as the HTTP service does when the caller supplies too many or nil preferences.
type invalidOptionsPreferencesClient struct {
	name  string
	calls int
}

func (c *invalidOptionsPreferencesClient) Name() string    { return c.name }
func (c *invalidOptionsPreferencesClient) Address() string { return c.name }
func (*invalidOptionsPreferencesClient) IsActive() bool    { return true }
func (*invalidOptionsPreferencesClient) IsSynced() bool    { return true }

func (c *invalidOptionsPreferencesClient) SubmitProposerPreferences(_ context.Context,
	_ *api.SubmitProposerPreferencesOpts,
) error {
	c.calls++

	return errors.Join(errors.New("too many proposer preferences"), consensusclient.ErrInvalidOptions)
}

// noOptionsPreferencesClient rejects every submission with ErrNoOptions, as the
// HTTP service does for a nil opts.  That is a separate sentinel from
// ErrInvalidOptions rather than one joined with it, so it needs its own double.
type noOptionsPreferencesClient struct {
	name  string
	calls int
}

func (c *noOptionsPreferencesClient) Name() string    { return c.name }
func (c *noOptionsPreferencesClient) Address() string { return c.name }
func (*noOptionsPreferencesClient) IsActive() bool    { return true }
func (*noOptionsPreferencesClient) IsSynced() bool    { return true }

func (c *noOptionsPreferencesClient) SubmitProposerPreferences(_ context.Context,
	_ *api.SubmitProposerPreferencesOpts,
) error {
	c.calls++

	return consensusclient.ErrNoOptions
}

func TestSubmitProposerPreferences(t *testing.T) {
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

	err = multiClient.(consensusclient.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, &api.SubmitProposerPreferencesOpts{
			Preferences: []*gloas.SignedProposerPreferences{{Message: &gloas.ProposerPreferences{}}},
		})
	require.NoError(t, err)
}

// TestSubmitProposerPreferencesDoesNotFailOverOnInvalidOptions ensures that a
// caller-side error is not treated as a provider failure: invalid input is
// invalid everywhere, so neither client should be tried twice or deactivated.
func TestSubmitProposerPreferencesDoesNotFailOverOnInvalidOptions(t *testing.T) {
	ctx := context.Background()

	client1 := &invalidOptionsPreferencesClient{name: "invalid 1"}
	client2 := &invalidOptionsPreferencesClient{name: "invalid 2"}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client1, client2}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, &api.SubmitProposerPreferencesOpts{
			Preferences: []*gloas.SignedProposerPreferences{{Message: &gloas.ProposerPreferences{}}},
		})
	require.ErrorIs(t, err, consensusclient.ErrInvalidOptions)
	require.Equal(t, 1, client1.calls)
	require.Equal(t, 0, client2.calls)
	// The first client is still active, so it remains the multi client's address.
	require.Equal(t, "invalid 1", multiClient.Address())
}

// TestSubmitProposerPreferencesDoesNotFailOverOnNoOptions is the ErrNoOptions
// half of the same property.  A nil opts is a caller error, not a provider
// fault, so it must not walk the cluster deactivating healthy clients: doing so
// leaves the multi client itself inactive and makes every subsequent call on it,
// not just this endpoint, pay a recheck first.
func TestSubmitProposerPreferencesDoesNotFailOverOnNoOptions(t *testing.T) {
	ctx := context.Background()

	client1 := &noOptionsPreferencesClient{name: "no options 1"}
	client2 := &noOptionsPreferencesClient{name: "no options 2"}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client1, client2}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, nil)
	require.ErrorIs(t, err, consensusclient.ErrNoOptions)
	require.Equal(t, 1, client1.calls)
	require.Equal(t, 0, client2.calls)
	require.True(t, multiClient.IsActive())
	require.Equal(t, "no options 1", multiClient.Address())
}
