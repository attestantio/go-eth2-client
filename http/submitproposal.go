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
	var body []byte
	var contentType ContentType
	var err error

	if s.enforceJSON {
		contentType = ContentTypeJSON
		body, err = s.submitProposalJSON(ctx, proposal)
	} else {
		contentType = ContentTypeSSZ
		body, err = s.submitProposalSSZ(ctx, proposal)
	}

	if err != nil {
		return nil, ContentTypeUnknown, err
	}

	return body, contentType, nil
}

func (*Service) submitProposalJSON(_ context.Context,
	proposal *api.VersionedSignedProposal,
) (
	[]byte,
	error,
) {
	var specJSON []byte
	var err error

	if err := proposal.AssertPresent(); err != nil {
		return nil, err
	}

	switch proposal.Version {
	case spec.DataVersionPhase0:
		specJSON, err = json.Marshal(proposal.Phase0)
	case spec.DataVersionAltair:
		specJSON, err = json.Marshal(proposal.Altair)
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(proposal.Bellatrix)
	case spec.DataVersionCapella:
		specJSON, err = json.Marshal(proposal.Capella)
	case spec.DataVersionDeneb:
		specJSON, err = json.Marshal(proposal.Deneb)
	case spec.DataVersionElectra:
		specJSON, err = json.Marshal(proposal.Electra)
	default:
		err = errors.New("unknown proposal version")
	}
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal JSON"), err)
	}

	return specJSON, nil
}

func (*Service) submitProposalSSZ(_ context.Context,
	proposal *api.VersionedSignedProposal,
) (
	[]byte,
	error,
) {
	var specSSZ []byte
	var err error

	if err := proposal.AssertPresent(); err != nil {
		return nil, err
	}

	switch proposal.Version {
	case spec.DataVersionPhase0:
		specSSZ, err = proposal.Phase0.MarshalSSZ()
	case spec.DataVersionAltair:
		specSSZ, err = proposal.Altair.MarshalSSZ()
	case spec.DataVersionBellatrix:
		specSSZ, err = proposal.Bellatrix.MarshalSSZ()
	case spec.DataVersionCapella:
		specSSZ, err = proposal.Capella.MarshalSSZ()
	case spec.DataVersionDeneb:
		specSSZ, err = proposal.Deneb.MarshalSSZ()
	case spec.DataVersionElectra:
		specSSZ, err = proposal.Electra.MarshalSSZ()
	default:
		err = errors.New("unknown proposal version")
	}
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal SSZ"), err)
	}

	return specSSZ, nil
}
