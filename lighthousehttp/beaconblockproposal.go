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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot uint64, randaoReveal []byte, graffiti []byte) (*spec.BeaconBlock, error) {
	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	url := fmt.Sprintf("/validator/block?slot=%d&randao_reveal=%#02x&graffiti=%#02x", slot, randaoReveal, fixedGraffiti)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon block proposal")
	}

	specReader, err := lhToSpec(ctx, respBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert lighthouse response to spec response")
	}
	var block *spec.BeaconBlock
	if err := json.NewDecoder(specReader).Decode(&block); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon block proposal")
	}

	// Ensure the data returned to us is as expected given our input.
	if block.Slot != slot {
		return nil, errors.New("beacon block proposal not for requested slot")
	}
	if !bytes.Equal(block.Body.RANDAOReveal, randaoReveal) {
		return nil, errors.New("beacon block proposal has incorrect RANDAO reveal")
	}
	if !bytes.Equal(block.Body.Graffiti, fixedGraffiti) {
		return nil, errors.New("beacon block proposal has incorrect graffiti")
	}

	return block, nil
}
