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
	lhReader, err := specToLH(ctx, bytes.NewReader(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to convert spec JSON to lighthouse JSON")
	}

	log.Trace().Msg("Sending to /validator/block")
	respBodyReader, cancel, err := s.post(ctx, "/validator/block", lhReader)
	if err != nil {
		return errors.Wrap(err, "failed to send to /validator/block")
	}
	defer cancel()

	var resp []byte
	if respBodyReader != nil {
		resp, err = ioutil.ReadAll(respBodyReader)
		if err != nil {
			resp = nil
		}
	}

	if err != nil {
		log.Debug().Err(err).Str("resp", string(resp)).Msg("POST to /validator/block failed")
		return errors.Wrap(err, fmt.Sprintf("failed to submit beacon block proposal: %s", resp))
	}

	return nil
}
