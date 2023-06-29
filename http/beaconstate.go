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
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0BeaconStateJSON struct {
	Data *phase0.BeaconState `json:"data"`
}

type altairBeaconStateJSON struct {
	Data *altair.BeaconState `json:"data"`
}

type bellatrixBeaconStateJSON struct {
	Data *bellatrix.BeaconState `json:"data"`
}

type capellaBeaconStateJSON struct {
	Data *capella.BeaconState `json:"data"`
}

type denebBeaconStateJSON struct {
	Data *deneb.BeaconState `json:"data"`
}

// BeaconState fetches a beacon state.
// N.B if the requested beacon state is not available this will return nil without an error.
func (s *Service) BeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	res, err := s.get2(ctx, fmt.Sprintf("/eth/v2/debug/beacon/states/%s", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon state")
	}
	if res.statusCode == http.StatusNotFound {
		return nil, nil
	}

	started := time.Now()
	defer fmt.Printf("Took %v\n", time.Since(started))
	switch res.contentType {
	case ContentTypeSSZ:
		return s.beaconStateFromSSZ(res)
	case ContentTypeJSON:
		return s.beaconStateFromJSON(res)
	default:
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}
}

func (s *Service) beaconStateFromSSZ(res *httpResponse) (*spec.VersionedBeaconState, error) {
	state := &spec.VersionedBeaconState{
		Version: res.consensusVersion,
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		state.Phase0 = &phase0.BeaconState{}
		if err := state.Phase0.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 beacon state")
		}
	case spec.DataVersionAltair:
		state.Altair = &altair.BeaconState{}
		if err := state.Altair.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair beacon state")
		}
	case spec.DataVersionBellatrix:
		state.Bellatrix = &bellatrix.BeaconState{}
		if err := state.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix beacon state")
		}
	case spec.DataVersionCapella:
		state.Capella = &capella.BeaconState{}
		if err := state.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella beacon state")
		}
	case spec.DataVersionDeneb:
		state.Deneb = &deneb.BeaconState{}
		if err := state.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb beacon state")
		}
	default:
		return nil, fmt.Errorf("unhandled state version %s", res.consensusVersion)
	}

	return state, nil
}

func (s *Service) beaconStateFromJSON(res *httpResponse) (*spec.VersionedBeaconState, error) {
	state := &spec.VersionedBeaconState{
		Version: res.consensusVersion,
	}

	reader := bytes.NewBuffer(res.body)

	switch state.Version {
	case spec.DataVersionPhase0:
		var resp phase0BeaconStateJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase 0 beacon state")
		}
		state.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairBeaconStateJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair beacon state")
		}
		state.Altair = resp.Data
	case spec.DataVersionBellatrix:
		var resp bellatrixBeaconStateJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix beacon state")
		}
		state.Bellatrix = resp.Data
	case spec.DataVersionCapella:
		var resp capellaBeaconStateJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella beacon state")
		}
		state.Capella = resp.Data
	case spec.DataVersionDeneb:
		var resp denebBeaconStateJSON
		if err := json.NewDecoder(reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse deneb beacon state")
		}
		state.Deneb = resp.Data
	}

	return state, nil
}
