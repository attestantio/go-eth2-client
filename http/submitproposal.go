// Copyright © 2023 - 2026 Attestant Limited.
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
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// SubmitProposal submits a proposal.
func (s *Service) SubmitProposal(ctx context.Context,
	opts *api.SubmitProposalOpts,
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

	body, contentType, err := s.submitProposalData(ctx, opts.Proposal)
	if err != nil {
		return err
	}

	endpoint := "/eth/v2/beacon/blocks"

	query := ""
	if opts.BroadcastValidation != nil {
		query = "broadcast_validation=" + opts.BroadcastValidation.String()
	}

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(opts.Proposal.Version.String())
	if opts.BuilderURL != "" {
		headers["Eth-Builder-Url"] = opts.BuilderURL
	}

	_, err = s.post(ctx, endpoint, query, &opts.Common, bytes.NewBuffer(body), contentType, headers)
	if err != nil {
		return errors.Join(errors.New("failed to submit proposal"), err)
	}

	return nil
}

func (s *Service) submitProposalData(ctx context.Context,
	proposal *api.VersionedSignedProposal,
) (
	[]byte,
	ContentType,
	error,
) {
	if err := proposal.AssertPresent(); err != nil {
		return nil, ContentTypeUnknown, err
	}

	var container any

	switch proposal.Version {
	case spec.DataVersionPhase0:
		container = proposal.Phase0
	case spec.DataVersionAltair:
		container = proposal.Altair
	case spec.DataVersionBellatrix:
		container = proposal.Bellatrix
	case spec.DataVersionCapella:
		container = proposal.Capella
	case spec.DataVersionDeneb:
		container = proposal.Deneb
	case spec.DataVersionElectra:
		container = proposal.Electra
	case spec.DataVersionFulu:
		container = proposal.Fulu
	case spec.DataVersionGloas:
		// A plain signed block, unlike the contents wrappers Deneb through Fulu
		// publish: post-Gloas the blobs travel in the execution payload
		// envelope, which is published through its own endpoint.
		container = proposal.Gloas
	default:
		// Unreachable: AssertPresent above rejects every version this switch
		// does not name.  Kept so the two cannot drift apart silently.
		return nil, ContentTypeUnknown, errors.New("unknown proposal version")
	}

	return s.marshalRequestBody(ctx, container)
}
