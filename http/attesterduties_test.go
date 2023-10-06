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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestAttesterDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.AttesterDutiesOpts
		expected []*apiv1.AttesterDuty
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "Good",
			opts: &api.AttesterDutiesOpts{
				Epoch:   0,
				Indices: []phase0.ValidatorIndex{0, 1},
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
			response, err := service.(client.AttesterDutiesProvider).AttesterDuties(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				require.Equal(t, test.errCode, err.(api.Error).StatusCode)
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.expected != nil {
					require.Equal(t, response.Data, test.expected)
				}
			}
		})
	}
}
