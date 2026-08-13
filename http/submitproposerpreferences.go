// Copyright © 2026 Attestant Limited.
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
	"fmt"
	nethttp "net/http"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

type proposerPreferencesList []*gloas.SignedProposerPreferences

const staticProposerPreferencesLimit uint64 = 64

func (s *Service) proposerPreferencesLimit(ctx context.Context) (uint64, error) {
	if !s.customSpecSupport {
		return staticProposerPreferencesLimit, nil
	}

	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return 0, err
	}

	minSeedLookahead, ok := response.Data["MIN_SEED_LOOKAHEAD"].(uint64)
	if !ok {
		return 0, ErrIncorrectType
	}
	slotsPerEpoch, ok := response.Data["SLOTS_PER_EPOCH"].(uint64)
	if !ok {
		return 0, ErrIncorrectType
	}
	if minSeedLookahead == ^uint64(0) || slotsPerEpoch > ^uint64(0)/(minSeedLookahead+1) {
		return 0, errors.New("proposer preferences limit overflows")
	}

	return (minSeedLookahead + 1) * slotsPerEpoch, nil
}

func (p proposerPreferencesList) MarshalSSZ() ([]byte, error) {
	var body []byte
	for _, preference := range p {
		encoded, err := preference.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		body = append(body, encoded...)
	}

	return body, nil
}

// SubmitProposerPreferences submits signed proposer preferences.
func (s *Service) SubmitProposerPreferences(ctx context.Context, preferences []*gloas.SignedProposerPreferences) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}
	limit, err := s.proposerPreferencesLimit(ctx)
	if err != nil {
		return err
	}
	if uint64(len(preferences)) > limit {
		return errors.Join(errors.New("too many proposer preferences"), client.ErrInvalidOptions)
	}
	for _, preference := range preferences {
		if preference == nil {
			return errors.Join(errors.New("nil proposer preference supplied"), client.ErrInvalidOptions)
		}
	}

	requestPreferences := proposerPreferencesList(preferences)
	if requestPreferences == nil {
		requestPreferences = proposerPreferencesList{}
	}
	body, contentType, err := s.marshalRequestBody(ctx, requestPreferences)
	if err != nil {
		return err
	}

	response, err := s.post(ctx,
		"/eth/v1/validator/proposer_preferences",
		"",
		&api.CommonOpts{},
		bytes.NewReader(body),
		contentType,
		map[string]string{"Eth-Consensus-Version": strings.ToLower(spec.DataVersionGloas.String())},
	)
	if err != nil {
		return errors.Join(errors.New("failed to submit proposer preferences"), err)
	}
	if response.statusCode != nethttp.StatusOK {
		return errors.Join(
			errors.New("failed to submit proposer preferences"),
			fmt.Errorf("unexpected status code %d", response.statusCode),
		)
	}

	return nil
}
