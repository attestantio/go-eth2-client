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
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type submitBeaconBlockResponseJSON struct {
	Empty   bool `json:"empty"`
	Present bool `json:"present"`
}

// SubmitBeaconBlock submits a beacon block.
func (s *Service) SubmitBeaconBlock(ctx context.Context, specBlock *spec.SignedBeaconBlock) error {
	specJSON, err := json.Marshal(specBlock)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	log.Trace().Msg("Sending to /validator/block")
	respBodyReader, err := s.post(ctx, "/validator/block", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to send to /validator/block")
	}

	var response submitBeaconBlockResponseJSON
	if err := json.NewDecoder(respBodyReader).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to parse submit beacon block response")
	}

	if response.Empty {
		return errors.New("beacon block proposal rejected as empty")
	}

	if response.Present {
		// This means that teku already knows about this block; that's fine.
		log.Debug().Msg("Beacon block already known by this node")
	}

	return nil
}
