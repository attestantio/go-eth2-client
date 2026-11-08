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

// preferencesErrorClient rejects every submission with a fixed caller error.
type preferencesErrorClient struct {
	name  string
	calls int
	err   error
}

func (c *preferencesErrorClient) Name() string    { return c.name }
func (c *preferencesErrorClient) Address() string { return c.name }
func (*preferencesErrorClient) IsActive() bool    { return true }
func (*preferencesErrorClient) IsSynced() bool    { return true }

func (c *preferencesErrorClient) SubmitProposerPreferences(context.Context, *api.SubmitProposerPreferencesOpts) error {
	c.calls++

	return c.err
}

// serviceOnlyClient is a connected client that does not implement the
// proposer-preferences submitter.
type serviceOnlyClient struct {
	name string
}

func (c *serviceOnlyClient) Name() string    { return c.name }
func (c *serviceOnlyClient) Address() string { return c.name }
func (*serviceOnlyClient) IsActive() bool    { return true }
func (*serviceOnlyClient) IsSynced() bool    { return true }

type countingPreferencesClient struct {
	name  string
	calls int
}

func (c *countingPreferencesClient) Name() string    { return c.name }
func (c *countingPreferencesClient) Address() string { return c.name }
func (*countingPreferencesClient) IsActive() bool    { return true }
func (*countingPreferencesClient) IsSynced() bool    { return true }

func (c *countingPreferencesClient) SubmitProposerPreferences(context.Context, *api.SubmitProposerPreferencesOpts) error {
	c.calls++

	return nil
}

func TestSubmitProposerPreferencesSkipsAClientWithoutTheInterface(t *testing.T) {
	ctx := context.Background()

	serving := &countingPreferencesClient{name: "serving"}
	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			&serviceOnlyClient{name: "no-interface"},
			serving,
		}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, &api.SubmitProposerPreferencesOpts{})
	require.NoError(t, err)
	require.Equal(t, 1, serving.calls)
	require.Equal(t, "no-interface", multiClient.Address())
}

func TestSubmitProposerPreferencesReportsAClientWithoutTheInterface(t *testing.T) {
	ctx := context.Background()

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{&serviceOnlyClient{name: "no-interface"}}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.ProposerPreferencesSubmitter).
		SubmitProposerPreferences(ctx, &api.SubmitProposerPreferencesOpts{})
	require.ErrorIs(t, err, multi.ErrCallNotSupported)
	require.Equal(t, "no-interface", multiClient.Address())
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

	callerErr := errors.Join(errors.New("too many proposer preferences"), consensusclient.ErrInvalidOptions)
	client1 := &preferencesErrorClient{name: "invalid 1", err: callerErr}
	client2 := &preferencesErrorClient{name: "invalid 2", err: callerErr}

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

	client1 := &preferencesErrorClient{name: "no options 1", err: consensusclient.ErrNoOptions}
	client2 := &preferencesErrorClient{name: "no options 2", err: consensusclient.ErrNoOptions}

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
