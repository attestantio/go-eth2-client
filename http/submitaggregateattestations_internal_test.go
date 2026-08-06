// Copyright © 2026 Attestant Limited.
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

package http

import (
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// TestCreateUnversionedAggregatesGloas verifies that a Gloas-versioned aggregate and proof is
// unwrapped for submission. Gloas carries its own container (gloas.Attestation merkleizes its
// aggregation bits as a progressive bitlist), so the arm must forward the Gloas-typed value
// unchanged; without the arm the switch hits its default and errors.
func TestCreateUnversionedAggregatesGloas(t *testing.T) {
	aggregate := &gloas.SignedAggregateAndProof{
		Message: &gloas.AggregateAndProof{AggregatorIndex: 42},
	}
	aggregateAndProofs := []*spec.VersionedSignedAggregateAndProof{
		{Version: spec.DataVersionGloas, Gloas: aggregate},
	}

	unversioned, err := createUnversionedAggregates(aggregateAndProofs)
	require.NoError(t, err)
	require.Len(t, unversioned, 1)
	require.Equal(t, aggregate, unversioned[0])
}

func TestCreateUnversionedAggregatesGloasNil(t *testing.T) {
	_, err := createUnversionedAggregates([]*spec.VersionedSignedAggregateAndProof{{Version: spec.DataVersionGloas}})
	require.ErrorIs(t, err, client.ErrInvalidOptions)
	require.EqualError(t, err, "nil gloas aggregate and proof supplied\ninvalid options")
}
