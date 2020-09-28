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

package lighthousehttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// ProposerDuties obtains proposer duties.
func (s *Service) ProposerDuties(ctx context.Context, epoch uint64, validators []client.ValidatorIDProvider) ([]*api.ProposerDuty, error) {
	if len(validators) == 0 {
		return s.proposerDuties(ctx, epoch)
	}

	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`{"epoch":`); err != nil {
		return nil, errors.Wrap(err, "failed to write header")
	}
	if _, err := reqBodyReader.WriteString(fmt.Sprintf("%d", epoch)); err != nil {
		return nil, errors.Wrap(err, "failed to write epoch")
	}
	if validators != nil {
		if _, err := reqBodyReader.WriteString(`,"pubkeys":[`); err != nil {
			return nil, errors.Wrap(err, "failed to write pubkeys")
		}
		for i := range validators {
			pubKey, err := validators[i].PubKey(ctx)
			if err != nil {
				// Warn but continue.
				log.Warn().Err(err).Msg("Failed to obtain public key for validator; skipping")
				continue
			}
			if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%#x"`, pubKey)); err != nil {
				return nil, errors.Wrap(err, "failed to write pubkey")
			}
			if i != len(validators)-1 {
				if _, err := reqBodyReader.WriteString(`,`); err != nil {
					return nil, errors.Wrap(err, "failed to write separator")
				}
			}
		}
		if _, err := reqBodyReader.WriteString(`]`); err != nil {
			return nil, errors.Wrap(err, "failed to write end of pubkeys")
		}
	}
	if _, err := reqBodyReader.WriteString(`}`); err != nil {
		return nil, errors.Wrap(err, "failed to write footer")
	}

	respBodyReader, err := s.post(ctx, "/validator/duties", &reqBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request proposer duties")
	}

	var resp []*dutyJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse proposer duties response")
	}

	proposerDuties := make([]*api.ProposerDuty, 0, len(resp))
	for i := range resp {
		for _, slot := range resp[i].BlockProposalSlots {
			proposerDuties = append(proposerDuties, &api.ProposerDuty{
				ValidatorIndex: resp[i].ValidatorIndex,
				Slot:           slot,
			})
		}
	}

	return proposerDuties, nil
}

func (s *Service) proposerDuties(ctx context.Context, epoch uint64) ([]*api.ProposerDuty, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/validator/duties/all?epoch=%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request proposer duties")
	}

	var resp []*dutyJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse proposer duties response")
	}

	proposerDuties := make([]*api.ProposerDuty, 0, len(resp))
	for i := range resp {
		for _, slot := range resp[i].BlockProposalSlots {
			proposerDuties = append(proposerDuties, &api.ProposerDuty{
				ValidatorIndex: resp[i].ValidatorIndex,
				Slot:           slot,
			})
		}
	}

	return proposerDuties, nil
}
