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

// BeaconState fetches a beacon state.
func (s *Service) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.State == "" {
		return nil, errors.New("no state specified")
	}

	url := fmt.Sprintf("/eth/v2/debug/beacon/states/%s", opts.State)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		return s.beaconStateFromSSZ(httpResponse)
	case ContentTypeJSON:
		return s.beaconStateFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

func (s *Service) beaconStateFromSSZ(res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
	response := &api.Response[*spec.VersionedBeaconState]{
		Data: &spec.VersionedBeaconState{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.BeaconState{}
		if err := response.Data.Phase0.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 beacon state")
		}
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.BeaconState{}
		if err := response.Data.Altair.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair beacon state")
		}
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &bellatrix.BeaconState{}
		if err := response.Data.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix beacon state")
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.BeaconState{}
		if err := response.Data.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella beacon state")
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.BeaconState{}
		if err := response.Data.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb beacon state")
		}
	default:
		return nil, fmt.Errorf("unhandled state version %s", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) beaconStateFromJSON(res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
	response := &api.Response[*spec.VersionedBeaconState]{
		Data: &spec.VersionedBeaconState{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &phase0.BeaconState{})
	case spec.DataVersionAltair:
		response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.BeaconState{})
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &bellatrix.BeaconState{})
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.BeaconState{})
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.BeaconState{})
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
