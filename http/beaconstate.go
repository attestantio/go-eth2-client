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

	"github.com/attestantio/go-eth2-client/spec/electra"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	dynssz "github.com/pk910/dynamic-ssz"
)

// BeaconState fetches a beacon state.
func (s *Service) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v2/debug/beacon/states/%s", opts.State)
	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		return s.beaconStateFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		return s.beaconStateFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

func (s *Service) beaconStateFromSSZ(ctx context.Context, res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
	response := &api.Response[*spec.VersionedBeaconState]{
		Data: &spec.VersionedBeaconState{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	var dynSSZ *dynssz.DynSsz
	if s.customSpecSupport {
		specs, err := s.Spec(ctx, &api.SpecOpts{})
		if err != nil {
			return nil, errors.Join(errors.New("failed to request specs"), err)
		}

		dynSSZ = dynssz.NewDynSsz(specs.Data)
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Phase0, res.body)
		} else {
			err = response.Data.Phase0.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode phase0 beacon state"), err)
		}
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Altair, res.body)
		} else {
			err = response.Data.Altair.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode altair beacon state"), err)
		}
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &bellatrix.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Bellatrix, res.body)
		} else {
			err = response.Data.Bellatrix.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode bellatrix beacon state"), err)
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Capella, res.body)
		} else {
			err = response.Data.Capella.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode capella beacon state"), err)
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Deneb, res.body)
		} else {
			err = response.Data.Deneb.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode deneb beacon state"), err)
		}
	case spec.DataVersionElectra:
		response.Data.Electra = &electra.BeaconState{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Electra, res.body)
		} else {
			err = response.Data.Electra.UnmarshalSSZ(res.body)
		}
		if err != nil {
			return nil, errors.Join(errors.New("failed to decode electra beacon state"), err)
		}
	default:
		return nil, fmt.Errorf("unhandled state version %s", res.consensusVersion)
	}

	return response, nil
}

func (*Service) beaconStateFromJSON(res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
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
	case spec.DataVersionElectra:
		response.Data.Electra, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &electra.BeaconState{})
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
