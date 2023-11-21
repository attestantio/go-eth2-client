// Copyright © 2023 Attestant Limited.
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
	"strings"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
)

// SubmitBlindedProposal submits a blinded proposal.
func (s *Service) SubmitBlindedProposal(ctx context.Context, proposal *api.VersionedSignedBlindedProposal) error {
	var specJSON []byte
	var err error

	if proposal == nil {
		return errors.New("no blinded proposal supplied")
	}

	switch proposal.Version {
	case spec.DataVersionPhase0:
		err = errors.New("blinded phase0 proposals not supported")
	case spec.DataVersionAltair:
		err = errors.New("blinded altair proposals not supported")
	case spec.DataVersionBellatrix:
		specJSON, err = json.Marshal(proposal.Bellatrix)
	case spec.DataVersionCapella:
		specJSON, err = json.Marshal(proposal.Capella)
	case spec.DataVersionDeneb:
		specJSON, err = json.Marshal(proposal.Deneb)
	default:
		err = errors.New("unknown proposal version")
	}
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(proposal.Version.String())
	_, err = s.post2(ctx, "/eth/v2/beacon/blinded_blocks", bytes.NewBuffer(specJSON), ContentTypeJSON, headers)
	if err != nil {
		return errors.Wrap(err, "failed to submit blinded proposal")
	}

	return nil
}
