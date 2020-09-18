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

// mockValidatorIDProvider implementes ValidatorIDProvider.
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
	}{
		{
			name:  "GoodWithValidators",
			epoch: 4092,
			validators: []client.ValidatorIDProvider{
				&mockValidatorIDProvider{
					index:  3243,
					pubKey: []byte{0x82, 0x37, 0x68, 0xb7, 0xd7, 0x94, 0x94, 0xf9, 0x50, 0xe3, 0x5d, 0x6d, 0xe9, 0xef, 0x4c, 0x0b, 0x4b, 0x1a, 0x5f, 0xef, 0xa4, 0xff, 0xba, 0x5f, 0x64, 0x88, 0xde, 0x99, 0x6b, 0x64, 0x44, 0x0b, 0x04, 0xc6, 0x03, 0xae, 0x13, 0x2c, 0xc4, 0x1e, 0x98, 0xe6, 0xe6, 0x55, 0xf1, 0x83, 0xcc, 0xd6},
				},
				&mockValidatorIDProvider{
					index:  3284,
					pubKey: []byte{0xb8, 0x02, 0xF8, 0xE6, 0x6B, 0x5A, 0xB0, 0x35, 0xE8, 0x1C, 0xD2, 0x66, 0xB3, 0x1A, 0xB4, 0xCF, 0x04, 0x9F, 0x80, 0xAE, 0x14, 0x7F, 0x96, 0x41, 0xDe, 0x14, 0x47, 0x6f, 0x13, 0x0b, 0xfc, 0x1d, 0x3f, 0xfb, 0x8d, 0xc6, 0x43, 0x08, 0xc0, 0x64, 0x90, 0xe9, 0x7a, 0x73, 0xda, 0x78, 0x2f, 0xfa},
				},
			},
		},
	}

	service, err := lighthousehttp.New(context.Background(), lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")))
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.ProposerDuties(context.Background(), test.epoch, test.validators)
			require.NoError(t, err)
			require.NotNil(t, duties)
		})
	}
}
