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
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0BeaconStateJSON struct {
	Data *phase0.BeaconState `json:"data"`
}

type altairBeaconStateJSON struct {
	Data *altair.BeaconState `json:"data"`
}

// BeaconState fetches a beacon state.
// N.B if the requested beacon state is not available this will return nil without an error.
func (s *Service) BeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	if s.supportsV2BeaconState {
		return s.beaconStateV2(ctx, stateID)
	}
	return s.beaconStateV1(ctx, stateID)
}

// beaconStateV1 fetches a beacon state from the V1 endpoint.
func (s *Service) beaconStateV1(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	url := fmt.Sprintf("/eth/v2/debug/beacon/states/%s", stateID)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon state")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp phase0BeaconStateJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon state")
	}

	return &spec.VersionedBeaconState{
		Version: spec.DataVersionPhase0,
		Phase0:  resp.Data,
	}, nil
}

// beaconStateV2 fetches a beacon state from the V2 endpoint.
func (s *Service) beaconStateV2(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	url := fmt.Sprintf("/eth/v2/debug/beacon/states/%s", stateID)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon state")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var dataBodyReader bytes.Buffer
	metadataReader := io.TeeReader(respBodyReader, &dataBodyReader)
	var metadata responseMetadata
	if err := json.NewDecoder(metadataReader).Decode(&metadata); err != nil {
		return nil, errors.Wrap(err, "failed to parse response")
	}
	res := &spec.VersionedBeaconState{
		Version: metadata.Version,
	}

	switch metadata.Version {
	case spec.DataVersionPhase0:
		var resp phase0BeaconStateJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase 0 beacon state")
		}
		res.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairBeaconStateJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair beacon state")
		}
		res.Altair = resp.Data
	}

	return res, nil
}
