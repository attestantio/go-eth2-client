// Copyright © 2020 Attestant Limited.
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

package lighthousehttp_test

import (
	"context"
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/lighthousehttp"
	"github.com/stretchr/testify/require"
)

func TestAttesterDuties(t *testing.T) {
	tests := []struct {
		name       string
		epoch      uint64
		validators []client.ValidatorIDProvider
	}{
		{
			name:  "Good",
			epoch: 1,
			validators: []client.ValidatorIDProvider{
				&testValidatorIDProvider{
					index:  0,
					pubKey: "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c",
				},
				&testValidatorIDProvider{
					index:  1,
					pubKey: "0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b",
				},
			},
		},
	}

	service, err := lighthousehttp.New(context.Background(),
		lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")),
		lighthousehttp.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.AttesterDuties(context.Background(), test.epoch, test.validators)
			require.NoError(t, err)
			require.NotNil(t, duties)
		})
	}
}
