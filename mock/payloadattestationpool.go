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

package mock

import (
	"context"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// PayloadAttestationPool obtains the payload attestation pool for the given options.
func (*Service) PayloadAttestationPool(_ context.Context,
	opts *api.PayloadAttestationPoolOpts,
) (
	*api.Response[[]*spec.VersionedPayloadAttestation],
	error,
) {
	var slot phase0.Slot
	if opts.Slot != nil {
		slot = *opts.Slot
	}

	data := make([]*spec.VersionedPayloadAttestation, 2)
	for i := range data {
		data[i] = &spec.VersionedPayloadAttestation{
			Version: spec.DataVersionGloas,
			Gloas: &gloas.PayloadAttestation{
				AggregationBits: bitfield.NewBitvector512(),
				Data: &gloas.PayloadAttestationData{
					Slot:              slot,
					PayloadPresent:    true,
					BlobDataAvailable: true,
				},
			},
		}
	}

	return &api.Response[[]*spec.VersionedPayloadAttestation]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
