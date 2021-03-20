// Copyright © 2021 Attestant Limited.
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

package v2_test

import (
	"context"
	"os"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain(t *testing.T) {
	tests := []struct {
		name   string
		epoch  spec.Epoch
		domain spec.DomainType
	}{
		{
			name:   "Good",
			epoch:  0,
			domain: [4]byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
		standardhttp.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			signatureDomain, err := service.Domain(context.Background(), test.domain, test.epoch)
			require.NoError(t, err)
			require.NotNil(t, signatureDomain)
			assert.Len(t, signatureDomain, 32)
		})
	}
}
