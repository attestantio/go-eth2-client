// Copyright © 2022 Attestant Limited.
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
	"errors"
	"math"
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestSubmitValidatorRegistrations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name          string
		registrations []*api.VersionedSignedValidatorRegistration
		expectErr     error
	}{
		{
			name: "InvalidVersion",
			registrations: []*api.VersionedSignedValidatorRegistration{
				{
					Version: 99999,
					V1: &v1.SignedValidatorRegistration{
						Message: &v1.ValidatorRegistration{},
					},
				},
			},
			expectErr: errors.New("unknown validator registration version\ninvalid options"),
		},
		{
			name: "InconsistentVersioning",
			registrations: []*api.VersionedSignedValidatorRegistration{
				{
					Version: spec.BuilderVersionV1,
					V1: &v1.SignedValidatorRegistration{
						Message: &v1.ValidatorRegistration{},
					},
				},
				{
					Version: math.MaxInt / 2,
					V1: &v1.SignedValidatorRegistration{
						Message: &v1.ValidatorRegistration{},
					},
				},
			},
			expectErr: errors.New("registrations must all be of the same version\ninvalid options"),
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.(client.ValidatorRegistrationsSubmitter).
				SubmitValidatorRegistrations(ctx, test.registrations)
			require.Equal(t, test.expectErr.Error(), err.Error())
		})
	}
}
