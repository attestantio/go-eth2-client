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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
)

// SubmitBlindedBeaconBlock submits a blinded beacon block.
//
// Deprecated: this will not work from the deneb hard-fork onwards.  Use SubmitBlindedProposal() instead.
func (s *Service) SubmitBlindedBeaconBlock(ctx context.Context, block *api.VersionedSignedBlindedBeaconBlock) error {
	var specJSON []byte
	var err error

	if block == nil {
		return errors.New("no blinded block supplied")
	}

	switch block.Version {
	case spec.DataVersionPhase0:
		err = errors.New("blinded phase0 blocks not supported")
	case spec.DataVersionAltair:
		err = errors.New("blinded altair blocks not supported")
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(block.Bellatrix)
	case spec.DataVersionCapella:
		specJSON, err = json.Marshal(block.Capella)
	case spec.DataVersionDeneb:
		specJSON, err = json.Marshal(block.Deneb)
	default:
		err = errors.New("unknown block version")
	}
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	_, err = s.post(ctx, "/eth/v1/beacon/blinded_blocks", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit blinded beacon block")
	}

	return nil
}
