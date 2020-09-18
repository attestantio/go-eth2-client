// Copyright © 2019 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/lighthousehttp"
	"github.com/stretchr/testify/require"
)

func TestVoluntaryExitDomain(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Good",
		},
	}

	service, err := lighthousehttp.New(context.Background(), lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")))
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			voluntaryExitDomain, err := service.VoluntaryExitDomain(context.Background())
			require.NoError(t, err)
			require.NotNil(t, voluntaryExitDomain)
		})
	}
}
