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

package prysmgrpc

import (
	"context"
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// Domain provides a domain for a given domain type at a given epoch.
func (s *Service) Domain(ctx context.Context, domainType spec.DomainType, epoch spec.Epoch) (spec.Domain, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling DomainData()")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.DomainData(opCtx, &ethpb.DomainRequest{
		Epoch:  uint64(epoch),
		Domain: domainType[:],
	})
	cancel()
	if err != nil {
		return spec.Domain{}, errors.Wrap(err, "failed to obtain domain")
	}

	var res spec.Domain
	copy(res[:], resp.SignatureDomain)
	log.Trace().
		Uint64("epoch", uint64(epoch)).
		Str("domain", fmt.Sprintf("%#x", domainType)).
		Str("signature_domain", fmt.Sprintf("%#x", res)).
		Msg("Signature domain obtained")
	return res, nil
}
