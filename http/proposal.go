// Copyright © 2020 - 2024 Attestant Limited.
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
	"errors"
	"fmt"
	"math/big"
	"strings"

	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"go.opentelemetry.io/otel"
)

// Proposal fetches a potential beacon block for signing.
func (s *Service) Proposal(ctx context.Context,
	opts *api.ProposalOpts,
) (
	*api.Response[*api.VersionedProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "Proposal")
	defer span.End()

	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.Slot == 0 {
		return nil, errors.Join(errors.New("no slot specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v3/validator/blocks/%d", opts.Slot)
	query := fmt.Sprintf("randao_reveal=%#x&graffiti=%#x", opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return nil, errors.Join(errors.New("randao reveal must be point at infinity if skip randao verification is set"), client.ErrInvalidOptions)
		}
		query = fmt.Sprintf("%s&skip_randao_verification", query)
	}

	if opts.BuilderBoostFactor == nil {
		query += "&builder_boost_factor=100"
	} else {
		query = fmt.Sprintf("%s&builder_boost_factor=%d", query, *opts.BuilderBoostFactor)
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request beacon block proposal"), err)
	}

	var response *api.Response[*api.VersionedProposal]
	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.beaconBlockProposalFromSSZ(httpResponse)
	case ContentTypeJSON:
		response, err = s.beaconBlockProposalFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
	if err != nil {
		return nil, err
	}

	// Ensure the data returned to us is as expected given our input.
	blockSlot, err := response.Data.Slot()
	if err != nil {
		return nil, err
	}
	if blockSlot != opts.Slot {
		return nil, errors.Join(fmt.Errorf("beacon block proposal for slot %d; expected %d", blockSlot, opts.Slot), client.ErrInconsistentResult)
	}

	// Only check the RANDAO reveal if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, errors.Join(fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], opts.RandaoReveal[:]), client.ErrInconsistentResult)
		}
	}

	return response, nil
}

func (s *Service) beaconBlockProposalFromSSZ(res *httpResponse) (*api.Response[*api.VersionedProposal], error) {
	response := &api.Response[*api.VersionedProposal]{
		Data: &api.VersionedProposal{
			Version:        res.consensusVersion,
			ConsensusValue: big.NewInt(0),
			ExecutionValue: big.NewInt(0),
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	if err := s.populateProposalDataFromHeaders(response, res.headers); err != nil {
		return nil, err
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.BeaconBlock{}
		err = response.Data.Phase0.UnmarshalSSZ(res.body)
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.BeaconBlock{}
		err = response.Data.Altair.UnmarshalSSZ(res.body)
	case spec.DataVersionBellatrix:
		if response.Data.Blinded {
			response.Data.BellatrixBlinded = &apiv1bellatrix.BlindedBeaconBlock{}
			err = response.Data.BellatrixBlinded.UnmarshalSSZ(res.body)
		} else {
			response.Data.Bellatrix = &bellatrix.BeaconBlock{}
			err = response.Data.Bellatrix.UnmarshalSSZ(res.body)
		}
	case spec.DataVersionCapella:
		if response.Data.Blinded {
			response.Data.CapellaBlinded = &apiv1capella.BlindedBeaconBlock{}
			err = response.Data.CapellaBlinded.UnmarshalSSZ(res.body)
		} else {
			response.Data.Capella = &capella.BeaconBlock{}
			err = response.Data.Capella.UnmarshalSSZ(res.body)
		}
	case spec.DataVersionDeneb:
		if response.Data.Blinded {
			response.Data.DenebBlinded = &apiv1deneb.BlindedBeaconBlock{}
			err = response.Data.DenebBlinded.UnmarshalSSZ(res.body)
		} else {
			response.Data.Deneb = &apiv1deneb.BlockContents{}
			err = response.Data.Deneb.UnmarshalSSZ(res.body)
		}
	case spec.DataVersionElectra:
		if response.Data.Blinded {
			response.Data.ElectraBlinded = &apiv1electra.BlindedBeaconBlock{}
			err = response.Data.ElectraBlinded.UnmarshalSSZ(res.body)
		} else {
			response.Data.Electra = &apiv1electra.BlockContents{}
			err = response.Data.Electra.UnmarshalSSZ(res.body)
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to decode %v SSZ beacon block (blinded: %v)", res.consensusVersion, response.Data.Blinded), err)
	}

	return response, nil
}

func (s *Service) beaconBlockProposalFromJSON(res *httpResponse) (*api.Response[*api.VersionedProposal], error) {
	response := &api.Response[*api.VersionedProposal]{
		Data: &api.VersionedProposal{
			Version:        res.consensusVersion,
			ConsensusValue: big.NewInt(0),
			ExecutionValue: big.NewInt(0),
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	if err := s.populateProposalDataFromHeaders(response, res.headers); err != nil {
		return nil, err
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &phase0.BeaconBlock{})
	case spec.DataVersionAltair:
		response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.BeaconBlock{})
	case spec.DataVersionBellatrix:
		if response.Data.Blinded {
			response.Data.BellatrixBlinded, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1bellatrix.BlindedBeaconBlock{})
		} else {
			response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &bellatrix.BeaconBlock{})
		}
	case spec.DataVersionCapella:
		if response.Data.Blinded {
			response.Data.CapellaBlinded, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1capella.BlindedBeaconBlock{})
		} else {
			response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.BeaconBlock{})
		}
	case spec.DataVersionDeneb:
		if response.Data.Blinded {
			response.Data.DenebBlinded, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1deneb.BlindedBeaconBlock{})
		} else {
			response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1deneb.BlockContents{})
		}
	case spec.DataVersionElectra:
		if response.Data.Blinded {
			response.Data.ElectraBlinded, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1electra.BlindedBeaconBlock{})
		} else {
			response.Data.Electra, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1electra.BlockContents{})
		}
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to decode %v JSON beacon block (blinded: %v)", res.consensusVersion, response.Data.Blinded), err)
	}

	return response, nil
}

func (*Service) populateProposalDataFromHeaders(response *api.Response[*api.VersionedProposal],
	headers map[string]string,
) error {
	for k, v := range headers {
		switch {
		case strings.EqualFold(k, "Eth-Execution-Payload-Blinded"):
			response.Data.Blinded = strings.EqualFold(v, "true")
		case strings.EqualFold(k, "Eth-Execution-Payload-Value"):
			var success bool
			response.Data.ExecutionValue, success = new(big.Int).SetString(v, 10)
			if !success {
				return fmt.Errorf("proposal header Eth-Execution-Payload-Value %s not a valid integer", v)
			}
		case strings.EqualFold(k, "Eth-Consensus-Block-Value"):
			var success bool
			response.Data.ConsensusValue, success = new(big.Int).SetString(v, 10)
			if !success {
				return fmt.Errorf("proposal header Eth-Consensus-Block-Value %s not a valid integer", v)
			}
		}
	}

	return nil
}
