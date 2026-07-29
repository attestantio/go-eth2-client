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

// A provider interface is only useful if every implementation carries it, but
// nothing in this repository forces that: the implementations are reached by
// type assertion at run time, so a missing method in multi, mock or a test
// client compiles cleanly and fails only when a caller asserts on it.
//
// These declarations turn that into a build failure.  They currently cover the
// payload timeliness committee surface, which the older providers predate; new
// providers belong in the block below rather than in a file of their own, so
// that the registry grows in one place.
var (
	_ client.PTCDutiesProvider = (*http.Service)(nil)
	_ client.PTCDutiesProvider = (*multi.Service)(nil)
	_ client.PTCDutiesProvider = (*mock.Service)(nil)
	_ client.PTCDutiesProvider = (*testclients.Erroring)(nil)
	_ client.PTCDutiesProvider = (*testclients.Sleepy)(nil)

	_ client.PayloadAttestationDataProvider = (*http.Service)(nil)
	_ client.PayloadAttestationDataProvider = (*multi.Service)(nil)
	_ client.PayloadAttestationDataProvider = (*mock.Service)(nil)
	_ client.PayloadAttestationDataProvider = (*testclients.Erroring)(nil)
	_ client.PayloadAttestationDataProvider = (*testclients.Sleepy)(nil)

	_ client.PayloadAttestationPoolProvider = (*http.Service)(nil)
	_ client.PayloadAttestationPoolProvider = (*multi.Service)(nil)
	_ client.PayloadAttestationPoolProvider = (*mock.Service)(nil)
	_ client.PayloadAttestationPoolProvider = (*testclients.Erroring)(nil)
	_ client.PayloadAttestationPoolProvider = (*testclients.Sleepy)(nil)

	_ client.PayloadAttestationMessagesSubmitter = (*http.Service)(nil)
	_ client.PayloadAttestationMessagesSubmitter = (*multi.Service)(nil)
	_ client.PayloadAttestationMessagesSubmitter = (*mock.Service)(nil)
	_ client.PayloadAttestationMessagesSubmitter = (*testclients.Erroring)(nil)
	_ client.PayloadAttestationMessagesSubmitter = (*testclients.Sleepy)(nil)
)
