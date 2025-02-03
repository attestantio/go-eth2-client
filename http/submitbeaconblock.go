// Copyright © 2020, 2024 Attestant Limited.
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
	"errors"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// SubmitBeaconBlock submits a beacon block.
//
// Deprecated: this will not work from the deneb hard-fork onwards.  Use SubmitProposal() instead.
func (s *Service) SubmitBeaconBlock(ctx context.Context, block *spec.VersionedSignedBeaconBlock) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	var specJSON []byte
	var err error

	if block == nil {
		return errors.Join(errors.New("no block supplied"), client.ErrInvalidOptions)
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
	case spec.DataVersionDeneb:
		specJSON, err = json.Marshal(block.Deneb)
	case spec.DataVersionElectra:
		specJSON, err = json.Marshal(block.Electra)
	default:
		err = errors.New("unknown block version")
	}
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v1/beacon/blocks"
	query := ""

	if _, err := s.post(ctx,
		endpoint,
		query,
		&api.CommonOpts{},
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		map[string]string{},
	); err != nil {
		return errors.Join(errors.New("failed to submit beacon block"), err)
	}

	return nil
}
