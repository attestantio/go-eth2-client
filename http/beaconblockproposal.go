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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0BeaconBlockProposalJSON struct {
	Data *phase0.BeaconBlock `json:"data"`
}

type altairBeaconBlockProposalJSON struct {
	Data *altair.BeaconBlock `json:"data"`
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	url := fmt.Sprintf("/eth/v2/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, fixedGraffiti)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon block proposal")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain beacon block proposal")
	}

	var dataBodyReader bytes.Buffer
	metadataReader := io.TeeReader(respBodyReader, &dataBodyReader)
	var metadata responseMetadata
	if err := json.NewDecoder(metadataReader).Decode(&metadata); err != nil {
		return nil, errors.Wrap(err, "failed to parse response")
	}
	res := &spec.VersionedBeaconBlock{
		Version: metadata.Version,
	}

	switch metadata.Version {
	case spec.DataVersionPhase0:
		var resp phase0BeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase 0 beacon block proposal")
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
		res.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairBeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair beacon block proposal")
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
		res.Altair = resp.Data
	}

	return res, nil
}
