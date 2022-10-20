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

package http

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
)

// SubmitBeaconBlock submits a beacon block.
func (s *Service) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	var specJSON []byte
	var err error

	if block == nil {
		return errors.New("no block supplied")
	}

	switch block.Version {
	case spec.DataVersionPhase0:
		specJSON, err = json.Marshal(block.Phase0)
	case spec.DataVersionAltair:
		specJSON, err = json.Marshal(block.Altair)
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(block.Bellatrix)
	case spec.DataVersionCapella:
		specJSON, err = json.Marshal(block.Capella)
	default:
		err = errors.New("unknown block version")
	}
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	_, err = s.post(ctx, "/eth/v1/beacon/blocks", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit beacon block")
	}

	return nil
}
