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

package lighthousehttp

import (
	"context"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

// SlotsPerEpoch provides the slots per epoch of the chain.
func (s *Service) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	if s.slotsPerEpoch == nil {
		respBodyReader, cancel, err := s.get(ctx, "/spec/slots_per_epoch")
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain slots per epoch")
		}
		defer cancel()

		slotsPerEpochBytes, err := ioutil.ReadAll(respBodyReader)
		if err != nil {
			return 0, errors.Wrap(err, "failed to read slots per epoch")
		}
		slotsPerEpoch, err := strconv.ParseUint(string(slotsPerEpochBytes), 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse slots per epoch")
		}

		s.slotsPerEpoch = &slotsPerEpoch
	}
	return *s.slotsPerEpoch, nil
}
