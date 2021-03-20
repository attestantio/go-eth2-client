// Copyright © 2021 Attestant Limited.
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

package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pkg/errors"
)

type beaconBlockProposalJSON struct {
	Data *spec.BeaconBlock `json:"data"`
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot spec.Slot, randaoReveal spec.BLSSignature, graffiti []byte) (*spec.BeaconBlock, error) {
	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	url := fmt.Sprintf("/eth/v1/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, fixedGraffiti)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon block proposal")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain beacon block proposal")
	}

	var resp beaconBlockProposalJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon block proposal")
	}

	// Ensure the data returned to us is as expected given our input.
	if resp.Data.Slot != slot {
		return nil, errors.New("beacon block proposal not for requested slot")
	}
	if !bytes.Equal(resp.Data.Body.RANDAOReveal[:], randaoReveal[:]) {
		return nil, errors.New("beacon block proposal has incorrect RANDAO reveal")
	}
	if !bytes.Equal(resp.Data.Body.Graffiti, fixedGraffiti) {
		return nil, errors.New("beacon block proposal has incorrect graffiti")
	}

	return resp.Data, nil
}
