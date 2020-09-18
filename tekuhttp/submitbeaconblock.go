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
	"fmt"
	"io/ioutil"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

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
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	var resp []byte
	if respBodyReader != nil {
		resp, err = ioutil.ReadAll(respBodyReader)
		if err != nil {
			resp = nil
		}
	}
	resp, err = ioutil.ReadAll(respBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to obtain error message for beacon block proposal")
	}
	if len(resp) > 0 {
		return fmt.Errorf("failed to submit beacon block proposal: %s", string(resp))
	}

	return nil
}
