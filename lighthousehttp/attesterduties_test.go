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
				&mockValidatorIDProvider{
					index:  0,
					pubKey: []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c},
				},
				&mockValidatorIDProvider{
					index:  1,
					pubKey: []byte{0xb8, 0x9b, 0xeb, 0xc6, 0x99, 0x76, 0x97, 0x26, 0xa3, 0x18, 0xc8, 0xe9, 0x97, 0x1b, 0xd3, 0x17, 0x12, 0x97, 0xc6, 0x1a, 0xea, 0x4a, 0x65, 0x78, 0xa7, 0xa4, 0xf9, 0x4b, 0x54, 0x7d, 0xcb, 0xa5, 0xba, 0xc1, 0x6a, 0x89, 0x10, 0x8b, 0x6b, 0x6a, 0x1f, 0xe3, 0x69, 0x5d, 0x1a, 0x87, 0x4a, 0x0b},
				},
			},
		},
	}

	service, err := lighthousehttp.New(context.Background(), lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")))
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.AttesterDuties(context.Background(), test.epoch, test.validators)
			require.NoError(t, err)
			require.NotNil(t, duties)
		})
	}
}
