// Copyright © 2023, 2024 Attestant Limited.
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
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// SubmitBlindedProposal submits a blinded proposal.
func (s *Service) SubmitBlindedProposal(ctx context.Context,
	opts *api.SubmitBlindedProposalOpts,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}
	if opts == nil {
		return client.ErrNoOptions
	}
	if opts.Proposal == nil {
		return errors.Join(errors.New("no proposal supplied"), client.ErrInvalidOptions)
	}

	var specJSON []byte
	var err error

	switch opts.Proposal.Version {
	case spec.DataVersionPhase0:
		err = errors.New("blinded phase0 proposals not supported")
	case spec.DataVersionAltair:
		err = errors.New("blinded altair proposals not supported")
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(opts.Proposal.Bellatrix)
	case spec.DataVersionCapella:
		specJSON, err = json.Marshal(opts.Proposal.Capella)
	case spec.DataVersionDeneb:
		specJSON, err = json.Marshal(opts.Proposal.Deneb)
	case spec.DataVersionElectra:
		specJSON, err = json.Marshal(opts.Proposal.Electra)
	default:
		err = errors.New("unknown proposal version")
	}
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v2/beacon/blinded_blocks"
	query := ""
	if opts.BroadcastValidation != nil {
		query = "broadcast_validation=" + opts.BroadcastValidation.String()
	}

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(opts.Proposal.Version.String())
	_, err = s.post(ctx, endpoint, query, &opts.Common, bytes.NewBuffer(specJSON), ContentTypeJSON, headers)
	if err != nil {
		return errors.Join(errors.New("failed to submit blinded proposal"), err)
	}

	return nil
}
