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

package http

import (
	"context"
	"time"
)

// SlotDuration provides the duration of a slot for the chain.
func (s *Service) SlotDuration(ctx context.Context) (time.Duration, error) {
	spec, err := s.Spec(ctx)
	if err != nil {
		return 0, err
	}

	return spec["SECONDS_PER_SLOT"].(time.Duration), nil
}
