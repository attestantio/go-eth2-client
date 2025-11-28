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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name   string
		epoch  phase0.Epoch
		domain phase0.DomainType
	}{
		{
			name:   "Good",
			epoch:  0,
			domain: [4]byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	service := testService(ctx, t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			signatureDomain, err := service.(client.DomainProvider).Domain(ctx, test.domain, test.epoch)
			require.NoError(t, err)
			require.NotNil(t, signatureDomain)
			assert.Len(t, signatureDomain, 32)
		})
	}
}
