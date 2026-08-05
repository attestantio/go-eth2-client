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

package client_test

import (
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/testclients"
)

// Provider assertions keep implementations in parity.
var (
	_ client.EPBSProposalProvider = (*http.Service)(nil)
	_ client.EPBSProposalProvider = (*multi.Service)(nil)
	_ client.EPBSProposalProvider = (*mock.Service)(nil)
	_ client.EPBSProposalProvider = (*testclients.Erroring)(nil)
	_ client.EPBSProposalProvider = (*testclients.Sleepy)(nil)

	_ client.ExecutionPayloadEnvelopeProvider = (*http.Service)(nil)
	_ client.ExecutionPayloadEnvelopeProvider = (*multi.Service)(nil)
	_ client.ExecutionPayloadEnvelopeProvider = (*mock.Service)(nil)
	_ client.ExecutionPayloadEnvelopeProvider = (*testclients.Erroring)(nil)
	_ client.ExecutionPayloadEnvelopeProvider = (*testclients.Sleepy)(nil)

	_ client.ExecutionPayloadProvider = (*http.Service)(nil)
	_ client.ExecutionPayloadProvider = (*multi.Service)(nil)
	_ client.ExecutionPayloadProvider = (*mock.Service)(nil)
	_ client.ExecutionPayloadProvider = (*testclients.Erroring)(nil)
	_ client.ExecutionPayloadProvider = (*testclients.Sleepy)(nil)

	_ client.ExecutionPayloadEnvelopeSubmitter = (*http.Service)(nil)
	_ client.ExecutionPayloadEnvelopeSubmitter = (*multi.Service)(nil)
	_ client.ExecutionPayloadEnvelopeSubmitter = (*mock.Service)(nil)
	_ client.ExecutionPayloadEnvelopeSubmitter = (*testclients.Erroring)(nil)
	_ client.ExecutionPayloadEnvelopeSubmitter = (*testclients.Sleepy)(nil)

	_ client.ExecutionPayloadBidSubmitter = (*http.Service)(nil)
	_ client.ExecutionPayloadBidSubmitter = (*multi.Service)(nil)
	_ client.ExecutionPayloadBidSubmitter = (*mock.Service)(nil)
	_ client.ExecutionPayloadBidSubmitter = (*testclients.Erroring)(nil)
	_ client.ExecutionPayloadBidSubmitter = (*testclients.Sleepy)(nil)
)
