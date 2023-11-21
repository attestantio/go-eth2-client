// Copyright © 2021, 2022 Attestant Limited.
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
	"github.com/pkg/errors"
)

// SlotsPerEpoch provides the slots per epoch of the chain.
//
// Deprecated: use Spec().
func (s *Service) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		slotsPerEpoch, err := client.(consensusclient.SlotsPerEpochProvider).SlotsPerEpoch(ctx)
		if err != nil {
			return nil, err
		}
		if slotsPerEpoch == 0 {
			return nil, errors.New("zero value not a valid response")
		}

		return slotsPerEpoch, nil
	}, nil)
	if err != nil {
		return 0, err
	}

	return res.(uint64), nil
}
