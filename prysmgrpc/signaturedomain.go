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

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SignatureDomain provides a signature domain for a given domain at a given epoch.
func (s *Service) SignatureDomain(ctx context.Context, domain []byte, epoch uint64) ([]byte, error) {
	if len(domain) != 4 {
		return nil, errors.New("invalid domain supplied")
	}

	client := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling DomainData()")
	resp, err := client.DomainData(ctx, &ethpb.DomainRequest{
		Epoch:  epoch,
		Domain: domain,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain signature domain")
	}

	log.Trace().
		Uint64("epoch", epoch).
		Str("domain", fmt.Sprintf("%#x", domain)).
		Str("signature_domain", fmt.Sprintf("%#x", resp.SignatureDomain)).
		Msg("Signature domain obtained")
	return resp.SignatureDomain, nil
}
