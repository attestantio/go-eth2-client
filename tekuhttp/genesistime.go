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

package tekuhttp

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// GenesisTime provides the genesis time of the chain.
func (s *Service) GenesisTime(ctx context.Context) (time.Time, error) {
	if s.genesisTime == nil {
		respBodyReader, cancel, err := s.get(ctx, "/node/genesis_time")
		if err != nil {
			return time.Now(), errors.Wrap(err, "failed to obtain genesis time")
		}
		defer cancel()

		genesisTimeBytes, err := ioutil.ReadAll(respBodyReader)
		if err != nil {
			return time.Now(), errors.Wrap(err, "failed to read genesis time")
		}
		genesisTimeBytes = bytes.Trim(genesisTimeBytes, "\"")
		genesisTimeInt, err := strconv.ParseInt(string(genesisTimeBytes), 10, 64)
		if err != nil {
			return time.Now(), errors.Wrap(err, "failed to parse genesis time")
		}
		genesisTime := time.Unix(genesisTimeInt, 0)
		s.genesisTime = &genesisTime
	}
	return *s.genesisTime, nil
}
