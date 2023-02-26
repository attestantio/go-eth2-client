// Copyright © 2020, 2021 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0BeaconBlockProposalJSON struct {
	Data *phase0.BeaconBlock `json:"data"`
}

type altairBeaconBlockProposalJSON struct {
	Data *altair.BeaconBlock `json:"data"`
}

type bellatrixBeaconBlockProposalJSON struct {
	Data *bellatrix.BeaconBlock `json:"data"`
}

type capellaBeaconBlockProposalJSON struct {
	Data *capella.BeaconBlock `json:"data"`
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	return s.beaconBlockProposal(ctx, slot, randaoReveal, fixedGraffiti)
}

func (s *Service) beaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	url := fmt.Sprintf("/eth/v2/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, graffiti)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
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
		res.Altair = resp.Data
	case spec.DataVersionBellatrix:
		var resp bellatrixBeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix beacon block proposal")
		}
		// Ensure the data returned to us is as expected given our input.
		if resp.Data.Slot != slot {
			return nil, errors.New("beacon block proposal not for requested slot")
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
		var resp capellaBeaconBlockProposalJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella beacon block proposal")
		}
		// Ensure the data returned to us is as expected given our input.
		if resp.Data.Slot != slot {
			return nil, errors.New("beacon block proposal not for requested slot")
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
