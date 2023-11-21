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

package multi

import (
	"context"
	"strings"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// SubmitAttestations submits attestations.
func (s *Service) SubmitAttestations(ctx context.Context,
	attestations []*phase0.Attestation,
) error {
	_, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		err := client.(consensusclient.AttestationsSubmitter).SubmitAttestations(ctx, attestations)
		if err != nil {
			return nil, err
		}

		return true, nil
	}, func(ctx context.Context, client consensusclient.Service, err error) (bool, error) {
		// We have received an error, decide if it requires us to fail over or not.
		provider := s.providerInfo(ctx, client)
		switch {
		case provider == "lighthouse" && strings.Contains(err.Error(), "PriorAttestationKnown"):
			// Lighthouse rejects duplicate attestations.  It is possible that an attestation sent
			// to another node already propagated to this node, or the caller is attempting to resend
			// an existing attestation, but either way it is not a failover-worthy error.
			log := s.log.With().Logger()
			log.Trace().Msg("Lighthouse rejected submission as it already knew about it")

			return false /* failover */, err
		case provider == "lighthouse" && strings.Contains(err.Error(), "UnknownHeadBlock"):
			// Lighthouse rejects an attestation for a block  that is not its current head.  We assume that
			// the request is valid and it is the node that it is somehow out of sync, so failover.
			log := s.log.With().Logger()
			log.Trace().Err(err).Msg("Lighthouse rejected submission as it did not know about the relevant head block")

			return true /* failover */, err
		default:
			// Any other error should result in a failover.

			return true /* failover */, err
		}
	})

	return err
}
