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

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
//
// Deprecated:  Use Spec().
func (s *Service) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		aggregators, err := client.(consensusclient.TargetAggregatorsPerCommitteeProvider).TargetAggregatorsPerCommittee(ctx)
		if err != nil {
			return nil, err
		}
		if aggregators == 0 {
			return nil, errors.New("zero value not a valid response")
		}

		return aggregators, nil
	}, nil)
	if err != nil {
		return 0, err
	}

	return res.(uint64), nil
}
