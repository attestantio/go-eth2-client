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

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Domain provides a domain for a given domain type at a given epoch.
func (s *Service) Domain(ctx context.Context,
	domainType phase0.DomainType,
	epoch phase0.Epoch,
) (
	phase0.Domain,
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		domain, err := client.(consensusclient.DomainProvider).Domain(ctx, domainType, epoch)
		if err != nil {
			return nil, err
		}
		return domain, nil
	}, nil)
	if err != nil {
		return phase0.Domain{}, err
	}
	if res == nil {
		return phase0.Domain{}, err
	}
	return res.(phase0.Domain), nil
}
