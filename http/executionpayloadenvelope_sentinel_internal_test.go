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
	"net/http"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
)

// TestExecutionPayloadEnvelopeSentinel pins the shape of the error a 404
// produces, which carries more weight than it looks like it should.
//
// A multi-client's failover decides whether a node is broken by unwrapping the
// error to an *api.Error and reading its status: an unwrappable 4xx is the node's
// answer, anything else is the node failing. This endpoint's 404 is the answer
// every node but the producing one gives, so a sentinel that hides the api.Error
// turns an ordinary call into the whole set looking broken, and the set gets
// deactivated a node at a time. Both properties therefore have to hold at once,
// and neither is implied by the other — which is why they are asserted together
// rather than left to the caller-facing errors.Is check alone.
func TestExecutionPayloadEnvelopeSentinel(t *testing.T) {
	notFound := &api.Error{
		Method:     http.MethodGet,
		StatusCode: http.StatusNotFound,
		Endpoint:   "/eth/v1/validator/execution_payload_envelopes/1/0x00",
		Data:       []byte(`{"message":"execution payload envelope not found for slot 1","code":404}`),
	}

	err := notFoundToSentinel(notFound)

	// What a caller matches on.
	require.ErrorIs(t, err, client.ErrNoExecutionPayloadEnvelope)

	// What a multi-client's failover matches on.
	var apiErr *api.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	require.Equal(t, 4, statusCodeFamily(apiErr.StatusCode),
		"a 4xx is what stops a multi-client treating this as the node failing")

	// And the body survives, which is the only way to tell an empty cache apart
	// from a node that does not serve this route at all — both answer 404.
	require.Contains(t, string(apiErr.Data), "not found for slot")

	// Any other status is not this endpoint's defined outcome and must stay a
	// plain error, so a real fault is not reported as an empty cache.
	for _, statusCode := range []int{http.StatusBadRequest, http.StatusForbidden, http.StatusInternalServerError} {
		other := notFoundToSentinel(&api.Error{StatusCode: statusCode})
		require.NotErrorIs(t, other, client.ErrNoExecutionPayloadEnvelope)
	}
}
