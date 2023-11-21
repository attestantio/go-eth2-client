// Copyright © 2020, 2021 Attestant Limited.
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
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestFinality(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name    string
		opts    *api.FinalityOpts
		err     string
		errCode int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "Invalid",
			opts: &api.FinalityOpts{
				State: "current",
			},
			errCode: 400,
		},
		{
			name: "Genesis",
			opts: &api.FinalityOpts{
				State: "genesis",
			},
		},
		{
			name: "Zero",
			opts: &api.FinalityOpts{
				State: "0",
			},
		},
		{
			name: "Head",
			opts: &api.FinalityOpts{
				State: "head",
			},
		},
		{
			name: "Finalized",
			opts: &api.FinalityOpts{
				State: "finalized",
			},
		},
		{
			name: "Justified",
			opts: &api.FinalityOpts{
				State: "justified",
			},
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.FinalityProvider).Finality(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
			}
		})
	}
}
