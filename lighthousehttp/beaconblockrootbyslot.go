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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// BeaconBlockRootBySlot fetches a block's root given its slot.
func (s *Service) BeaconBlockRootBySlot(ctx context.Context, slot uint64) ([]byte, error) {
	url := fmt.Sprintf("/beacon/block_root?slot=%d", slot)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon block root")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	rootBytes, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read block root")
	}
	root, err := hex.DecodeString(strings.TrimPrefix(strings.Trim(string(rootBytes), `"`), "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse block root")
	}

	return root, nil
}
