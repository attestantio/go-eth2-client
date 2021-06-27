// Copyright © 2020 Attestant Limited.
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

package v1

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type stateRootJSON struct {
	Data *stateRootDataJSON `json:"data"`
}

type stateRootDataJSON struct {
	Root string `json:"root"`
}

// BeaconStateRoot fetches a beacon state root given a state ID.
func (s *Service) BeaconStateRoot(ctx context.Context, stateID string) (*spec.Root, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/root", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request state root")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain state root")
	}

	var stateRootJSON stateRootJSON
	if err := json.NewDecoder(respBodyReader).Decode(&stateRootJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse state root")
	}

	bytes, err := hex.DecodeString(strings.TrimPrefix(stateRootJSON.Data.Root, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse state root value")
	}
	var stateRoot spec.Root
	copy(stateRoot[:], bytes)

	return &stateRoot, nil
}
