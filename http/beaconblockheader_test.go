// Copyright © 2020 - 2023 Attestant Limited.
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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockHeader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.BeaconBlockHeaderOpts
		expected *apiv1.BeaconBlockHeader
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "Genesis",
			opts: &api.BeaconBlockHeaderOpts{Block: "0"},
			expected: &apiv1.BeaconBlockHeader{
				Root:      *mustParseRoot("0x4d611d5b93fdab69013a7f0a2f961caca0c853f87cfe9595fe50038163079360"),
				Canonical: true,
				Header: &phase0.SignedBeaconBlockHeader{
					Message: &phase0.BeaconBlockHeader{
						Slot:          0,
						ProposerIndex: 0,
						ParentRoot:    *mustParseRoot("0x0000000000000000000000000000000000000000000000000000000000000000"),
						StateRoot:     *mustParseRoot("0x7e76880eb67bbdc86250aa578958e9d0675e64e714337855204fb5abaaf82c2b"),
						BodyRoot:      *mustParseRoot("0xccb62460692be0ec813b56be97f68a82cf57abc102e27bf49ebf4190ff22eedd"),
					},
					Signature: *mustParseSignature("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
		},
		{
			name: "Head",
			opts: &api.BeaconBlockHeaderOpts{Block: "head"},
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BeaconBlockHeadersProvider).BeaconBlockHeader(ctx, test.opts)
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
				if test.expected != nil {
					require.Equal(t, test.expected, response.Data)
				}
			}
		})
	}
}
