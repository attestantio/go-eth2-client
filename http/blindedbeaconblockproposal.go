// Copyright © 2022, 2023 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type bellatrixBlindedBeaconBlockProposalJSON struct {
	Data *apiv1bellatrix.BlindedBeaconBlock `json:"data"`
}

type capellaBlindedBeaconBlockProposalJSON struct {
	Data *apiv1capella.BlindedBeaconBlock `json:"data"`
}

type denebBlindedBeaconBlockProposalJSON struct {
	Data *apiv1deneb.BlindedBeaconBlock `json:"data"`
}

// BlindedBeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BlindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti []byte) (*api.VersionedBlindedBeaconBlock, error) {
	// Graffiti should be 32 bytes.
	var fixedGraffiti [32]byte
	copy(fixedGraffiti[:], graffiti)

	return s.blindedBeaconBlockProposal(ctx, slot, randaoReveal, fixedGraffiti)
}

// blindedBeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) blindedBeaconBlockProposal(ctx context.Context, slot phase0.Slot, randaoReveal phase0.BLSSignature, graffiti [32]byte) (*api.VersionedBlindedBeaconBlock, error) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "blindedBeaconBlockProposal")
	defer span.End()

	res, err := s.get2(ctx, fmt.Sprintf("/eth/v1/validator/blinded_blocks/%d?randao_reveal=%#x&graffiti=%#x", slot, randaoReveal, graffiti))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to request beacon block proposal")
		return nil, errors.Wrap(err, "failed to request blinded beacon block proposal")
	}
	if res.statusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "Client returned 404")
		return nil, nil
	}

	var block *api.VersionedBlindedBeaconBlock
	switch res.contentType {
	case ContentTypeSSZ:
		block, err = s.blindedBeaconBlockProposalFromSSZ(res)
	case ContentTypeJSON:
		block, err = s.blindedBeaconBlockProposalFromJSON(res)
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
		return nil, errors.New("blinded beacon block proposal not for requested slot")
	}

	// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := block.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], randaoReveal[:]) {
			span.SetStatus(codes.Error, fmt.Sprintf("Proposal RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], randaoReveal[:]))
			return nil, fmt.Errorf("blinded beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], randaoReveal[:])
		}

		blockGraffiti, err := block.Graffiti()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockGraffiti[:], graffiti[:]) {
			span.SetStatus(codes.Error, fmt.Sprintf("Proposal graffiti %#x; expected %#x", blockGraffiti[:], graffiti[:]))
			return nil, fmt.Errorf("blinded beacon block proposal has graffiti %#x; expected %#x", blockGraffiti[:], graffiti[:])
		}
	}

	return block, nil
}

func (s *Service) blindedBeaconBlockProposalFromSSZ(res *httpResponse) (*api.VersionedBlindedBeaconBlock, error) {
	block := &api.VersionedBlindedBeaconBlock{
		Version: res.consensusVersion,
	}

	switch res.consensusVersion {
	case spec.DataVersionBellatrix:
		block.Bellatrix = &apiv1bellatrix.BlindedBeaconBlock{}
		if err := block.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix blinded beacon block proposal")
		}
	case spec.DataVersionCapella:
		block.Capella = &apiv1capella.BlindedBeaconBlock{}
		if err := block.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella blinded beacon block proposal")
		}
	case spec.DataVersionDeneb:
		block.Deneb = &apiv1deneb.BlindedBeaconBlock{}
		if err := block.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb blinded beacon block proposal")
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	return block, nil
}

func (s *Service) blindedBeaconBlockProposalFromJSON(res *httpResponse) (*api.VersionedBlindedBeaconBlock, error) {
	block := &api.VersionedBlindedBeaconBlock{
		Version: res.consensusVersion,
	}

	reader := bytes.NewBuffer(res.body)
	switch block.Version {
	case spec.DataVersionBellatrix:
		var resp bellatrixBlindedBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix blinded beacon block proposal")
		}
		block.Bellatrix = resp.Data
	case spec.DataVersionCapella:
		var resp capellaBlindedBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella blinded beacon block proposal")
		}
		block.Capella = resp.Data
	case spec.DataVersionDeneb:
		var resp denebBlindedBeaconBlockProposalJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse deneb blinded beacon block proposal")
		}
		block.Deneb = resp.Data
	default:
		return nil, fmt.Errorf("unsupported block version %s", res.consensusVersion)
	}

	return block, nil
}
