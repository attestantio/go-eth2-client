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
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Service) SignedBeaconBlock(ctx context.Context,
	opts *api.SignedBeaconBlockOpts,
) (
	*api.Response[*spec.VersionedSignedBeaconBlock],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	url := fmt.Sprintf("/eth/v2/beacon/blocks/%s", opts.Block)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	var response *api.Response[*spec.VersionedSignedBeaconBlock]
	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.signedBeaconBlockFromSSZ(httpResponse)
	case ContentTypeJSON:
		response, err = s.signedBeaconBlockFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) signedBeaconBlockFromSSZ(res *httpResponse) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	response := &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.SignedBeaconBlock{}
		if err := response.Data.Phase0.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 signed beacon block")
		}
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.SignedBeaconBlock{}
		if err := response.Data.Altair.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair signed beacon block")
		}
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &bellatrix.SignedBeaconBlock{}
		if err := response.Data.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix signed beacon block")
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.SignedBeaconBlock{}
		if err := response.Data.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella signed beacon block")
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.SignedBeaconBlock{}
		if err := response.Data.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb signed block contents")
		}
	default:
		return nil, fmt.Errorf("unhandled block version %s", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) signedBeaconBlockFromJSON(res *httpResponse) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	response := &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &phase0.SignedBeaconBlock{})
	case spec.DataVersionAltair:
		response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.SignedBeaconBlock{})
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &bellatrix.SignedBeaconBlock{})
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.SignedBeaconBlock{})
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.SignedBeaconBlock{})
	default:
		return nil, fmt.Errorf("unhandled version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
