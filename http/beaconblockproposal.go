// Copyright © 2020 - 2023 Attestant Limited.
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
	"net/http"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

type denebBeaconBlockProposalJSON struct {
	Data *deneb.BeaconBlock `json:"data"`
}

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*spec.VersionedBeaconBlock, error) {
	// Graffiti should be 32 bytes.
	var fixedGraffiti [32]byte
	copy(fixedGraffiti[:], graffiti)

	return s.beaconBlockProposal(ctx, slot, randaoReveal, fixedGraffiti)
}

//nolint:gocyclo
func (s *Service) beaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti [32]byte) (*spec.VersionedBeaconBlock, error) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "beaconBlockProposal")
	defer span.End()

	url := fmt.Sprintf("/eth/v2/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, graffiti)

	res, err := s.get2(ctx, url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to request beacon block proposal")
		return nil, errors.Wrap(err, "failed to request beacon block proposal")
	}
	if res.statusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "Client returned 404")
		return nil, nil
	}

	var block *spec.VersionedBeaconBlock
	switch res.contentType {
	case ContentTypeSSZ:
		block, err = s.beaconBlockProposalFromSSZ(res)
	case ContentTypeJSON:
		block, err = s.beaconBlockProposalFromJSON(res)
	default:
		span.SetStatus(codes.Error, fmt.Sprintf("Unhandled content type %s", res.contentType))
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode body")
		return nil, err
	}

	// Ensure the data returned to us is as expected given our input.
	blockSlot, err := block.Slot()
	if err != nil {
		return nil, err
	}
	if blockSlot != slot {
		span.SetStatus(codes.Error, fmt.Sprintf("Proposal slot %d; expected %d", blockSlot, slot))
		return nil, errors.New("beacon block proposal not for requested slot")
	}

	// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := block.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], randaoReveal[:]) {
			span.SetStatus(codes.Error, fmt.Sprintf("Proposal RANDAO reveal %#x; expected requested %#x", blockRandaoReveal[:], randaoReveal[:]))
			return nil, fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], randaoReveal[:])
		}

		blockGraffiti, err := block.Graffiti()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockGraffiti[:], graffiti[:]) {
			span.SetStatus(codes.Error, fmt.Sprintf("Proposal graffiti %#x; expected %#x", blockGraffiti[:], graffiti[:]))
			return nil, fmt.Errorf("beacon block proposal has graffiti %#x; expected %#x", blockGraffiti[:], graffiti[:])
		}
	}

	return block, nil
}

func (s *Service) beaconBlockProposalFromSSZ(res *httpResponse) (*spec.VersionedBeaconBlock, error) {
	block := &spec.VersionedBeaconBlock{
		Version: res.consensusVersion,
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		block.Phase0 = &phase0.BeaconBlock{}
		if err := block.Phase0.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 beacon block proposal")
		}
	case spec.DataVersionAltair:
		block.Altair = &altair.BeaconBlock{}
		if err := block.Altair.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair beacon block proposal")
		}
	case spec.DataVersionBellatrix:
		block.Bellatrix = &bellatrix.BeaconBlock{}
		if err := block.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix beacon block proposal")
		}
	case spec.DataVersionCapella:
		block.Capella = &capella.BeaconBlock{}
		if err := block.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella beacon block proposal")
		}
	case spec.DataVersionDeneb:
		block.Deneb = &deneb.BeaconBlock{}
		if err := block.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb beacon block proposal")
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	return block, nil
}

func (s *Service) beaconBlockProposalFromJSON(res *httpResponse) (*spec.VersionedBeaconBlock, error) {
	block := &spec.VersionedBeaconBlock{
		Version: res.consensusVersion,
	}

	reader := bytes.NewBuffer(res.body)
	switch block.Version {
	case spec.DataVersionPhase0:
		var resp phase0BeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase0 beacon block proposal")
		}
		block.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair beacon block proposal")
		}
		block.Altair = resp.Data
	case spec.DataVersionBellatrix:
		var resp bellatrixBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix beacon block proposal")
		}
		block.Bellatrix = resp.Data
	case spec.DataVersionCapella:
		var resp capellaBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella beacon block proposal")
		}
		block.Capella = resp.Data
	case spec.DataVersionDeneb:
		var resp denebBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse deneb beacon block proposal")
		}
		block.Deneb = resp.Data
	default:
		return nil, fmt.Errorf("unsupported block version %s", res.consensusVersion)
	}

	return block, nil
}
