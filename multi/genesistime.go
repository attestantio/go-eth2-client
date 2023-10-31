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
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

// GenesisTime provides the genesis time of the chain.
//
// Deprecated: use Genesis().
func (s *Service) GenesisTime(ctx context.Context) (time.Time, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		genesisTime, err := client.(consensusclient.GenesisTimeProvider).GenesisTime(ctx)
		if err != nil {
			return nil, err
		}
		if genesisTime.IsZero() {
			return nil, errors.New("zero genesis time not a valid response")
		}

		return genesisTime, nil
	}, nil)
	if err != nil {
		return time.Time{}, err
	}

	return res.(time.Time), nil
}
