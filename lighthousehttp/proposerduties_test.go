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

// mockValidatorIDProvider implements ValidatorIDProvider.
type mockValidatorIDProvider struct {
	index  uint64
	pubKey []byte
}

func (m *mockValidatorIDProvider) Index(ctx context.Context) (uint64, error) {
	return m.index, nil
}
func (m *mockValidatorIDProvider) PubKey(ctx context.Context) ([]byte, error) {
	return m.pubKey, nil
}

func TestProposerDuties(t *testing.T) {
	tests := []struct {
		name       string
		epoch      uint64
		validators []client.ValidatorIDProvider
		expected   int
	}{
		{
			name:     "Good",
			epoch:    4092,
			expected: 32,
		},
		{
			name:  "GoodWithValidators",
			epoch: 4092,
			validators: []client.ValidatorIDProvider{
				&mockValidatorIDProvider{
					index: 16056,
					pubKey: []byte{
						0x95, 0x53, 0xa6, 0x3a, 0x58, 0xd3, 0xa7, 0x76, 0xa2, 0x48, 0x31, 0x84, 0xe5, 0xaf, 0x37, 0xae,
						0xdf, 0x13, 0x1b, 0x82, 0xef, 0x1e, 0x0b, 0xcb, 0xa7, 0xb3, 0xc0, 0x18, 0x18, 0xf4, 0x90, 0x37,
						0x1a, 0xac, 0x0c, 0x6f, 0x9a, 0x32, 0x7f, 0xb7, 0xeb, 0x89, 0x19, 0x0a, 0xf7, 0xb0, 0x85, 0xa5,
					},
				},
				&mockValidatorIDProvider{
					index: 35476,
					pubKey: []byte{
						0x92, 0x16, 0x09, 0x1f, 0x3e, 0x4f, 0xe0, 0xb0, 0x56, 0x2a, 0x6c, 0x5b, 0xf6, 0xe8, 0xc3, 0x5c,
						0xf0, 0xc3, 0xb3, 0x21, 0xb6, 0xf4, 0x15, 0xde, 0x66, 0x31, 0xd7, 0xd1, 0x2e, 0x58, 0x60, 0x3e,
						0x1e, 0x23, 0xc8, 0xd7, 0x8f, 0x44, 0x9b, 0x60, 0x1f, 0x8d, 0x24, 0x4d, 0x26, 0xf7, 0x0a, 0xa7,
					},
				},
			},
			expected: 2,
		},
	}

	service, err := lighthousehttp.New(context.Background(), lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")))
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.ProposerDuties(context.Background(), test.epoch, test.validators)
			require.NoError(t, err)
			require.NotNil(t, duties)
			require.Equal(t, test.expected, len(duties))
		})
	}
}
