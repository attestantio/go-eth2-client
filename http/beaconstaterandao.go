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
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type stateRandaoJSON struct {
	Data *stateRandaoDataJSON `json:"data"`
}

type stateRandaoDataJSON struct {
	Randao string `json:"randao"`
}

// BeaconStateRandao fetches a beacon state RANDAO given a state ID.
func (s *Service) BeaconStateRandao(ctx context.Context, stateID string) (*phase0.Root, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/randao", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request state RANDAO")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var data stateRandaoJSON
	if err := json.NewDecoder(respBodyReader).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "failed to parse state RANDAO")
	}

	bytes, err := hex.DecodeString(strings.TrimPrefix(data.Data.Randao, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse state RANDAO value")
	}
	var stateRandao phase0.Root
	copy(stateRandao[:], bytes)

	return &stateRandao, nil
}
