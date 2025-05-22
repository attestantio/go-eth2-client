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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AttestationPool fetches the attestation pool for the given slot.
func (*Service) AttestationPool(_ context.Context,
	_ *api.AttestationPoolOpts,
) (
	*api.Response[[]*spec.VersionedAttestation],
	error,
) {
	data := make([]*spec.VersionedAttestation, 5)
	for i := 0; i < 5; i++ {
		data[i] = &spec.VersionedAttestation{
			Version: spec.DataVersionPhase0,
			Phase0: &phase0.Attestation{
				Data: &phase0.AttestationData{
					Source: &phase0.Checkpoint{},
					Target: &phase0.Checkpoint{},
				},
			},
		}
	}

	return &api.Response[[]*spec.VersionedAttestation]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
