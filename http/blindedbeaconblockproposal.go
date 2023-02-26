// Copyright © 2022 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type bellatrixBlindedBeaconBlockProposalJSON struct {
	Data *apiv1bellatrix.BlindedBeaconBlock `json:"data"`
}

type capellaBlindedBeaconBlockProposalJSON struct {
	Data *apiv1capella.BlindedBeaconBlock `json:"data"`
}

// BlindedBeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BlindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*api.VersionedBlindedBeaconBlock, error) {
	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	return s.blindedBeaconBlockProposal(ctx, slot, randaoReveal, fixedGraffiti)
}

// blindedBeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) blindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*api.VersionedBlindedBeaconBlock, error) {
	url := fmt.Sprintf("/eth/v1/validator/blinded_blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, graffiti)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request blinded beacon block proposal")
	}
	if respBodyReader == nil {
		return nil, errors.New("blinded beacon block proposal response empty")
	}

	var dataBodyReader bytes.Buffer
	metadataReader := io.TeeReader(respBodyReader, &dataBodyReader)
	var metadata responseMetadata
	if err := json.NewDecoder(metadataReader).Decode(&metadata); err != nil {
		return nil, errors.Wrap(err, "failed to parse response")
	}
	res := &api.VersionedBlindedBeaconBlock{
		Version: metadata.Version,
	}

	switch metadata.Version {
	case spec.DataVersionBellatrix:
		var resp bellatrixBlindedBeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix blinded beacon block proposal")
		}
		// Ensure the data returned to us is as expected given our input.
		if resp.Data.Slot != slot {
			return nil, errors.New("blinded beacon block proposal not for requested slot")
		}
		// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
		// as the returned values will be decided by the middleware.
		if !s.connectedToDVTMiddleware {
			if !bytes.Equal(resp.Data.Body.RANDAOReveal[:], randaoReveal[:]) {
				return nil, fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", resp.Data.Body.RANDAOReveal[:], randaoReveal[:])
			}
			if !bytes.Equal(resp.Data.Body.Graffiti[:], graffiti) {
				return nil, fmt.Errorf("beacon block proposal has graffiti %#x; expected %#x", resp.Data.Body.Graffiti[:], graffiti)
			}
		}
		res.Bellatrix = resp.Data
	case spec.DataVersionCapella:
		var resp capellaBlindedBeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella blinded beacon block proposal")
		}
		// Ensure the data returned to us is as expected given our input.
		if resp.Data.Slot != slot {
			return nil, errors.New("blinded beacon block proposal not for requested slot")
		}
		// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
		// as the returned values will be decided by the middleware.
		if !s.connectedToDVTMiddleware {
			if !bytes.Equal(resp.Data.Body.RANDAOReveal[:], randaoReveal[:]) {
				return nil, fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", resp.Data.Body.RANDAOReveal[:], randaoReveal[:])
			}
			if !bytes.Equal(resp.Data.Body.Graffiti[:], graffiti) {
				return nil, fmt.Errorf("beacon block proposal has graffiti %#x; expected %#x", resp.Data.Body.Graffiti[:], graffiti)
			}
		}
		res.Capella = resp.Data
	default:
		return nil, fmt.Errorf("unsupported block version %s", metadata.Version)
	}

	return res, nil
}
