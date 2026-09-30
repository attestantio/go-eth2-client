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
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

type proposerPreferencesList []*gloas.SignedProposerPreferences

// staticProposerPreferencesLimit is the mainnet-preset proposer lookahead length,
// used whenever the chain's own values are not available.
const staticProposerPreferencesLimit uint64 = 64

// proposerPreferencesLimit returns the maximum number of proposer preferences that
// may be submitted at once: the proposer lookahead length,
// (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH.
//
// This is a client-side sanity bound and the node remains the authority on what it
// accepts, so every way of failing to derive the chain's own value falls back to the
// mainnet preset rather than failing the submission.  That is what allows the limit
// to be derived unconditionally: Spec() is cached, and a node that does not publish
// the spec, or publishes it without these keys, costs a fallback rather than an error.
func (s *Service) proposerPreferencesLimit(ctx context.Context) uint64 {
	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return staticProposerPreferencesLimit
	}

	minSeedLookahead, ok := response.Data["MIN_SEED_LOOKAHEAD"].(uint64)
	if !ok {
		return staticProposerPreferencesLimit
	}
	slotsPerEpoch, ok := response.Data["SLOTS_PER_EPOCH"].(uint64)
	if !ok {
		return staticProposerPreferencesLimit
	}
	if minSeedLookahead == ^uint64(0) || slotsPerEpoch > ^uint64(0)/(minSeedLookahead+1) {
		return staticProposerPreferencesLimit
	}

	return (minSeedLookahead + 1) * slotsPerEpoch
}

// MarshalSSZ encodes the list as an SSZ list of fixed-size elements, which is the
// concatenation of the elements' own encodings.
func (p proposerPreferencesList) MarshalSSZ() ([]byte, error) {
	body := make([]byte, 0, len(p)*(&gloas.SignedProposerPreferences{}).SizeSSZ())
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
func (s *Service) SubmitProposerPreferences(ctx context.Context, opts *api.SubmitProposerPreferencesOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if uint64(len(opts.Preferences)) > s.proposerPreferencesLimit(ctx) {
		return errors.Join(errors.New("too many proposer preferences"), client.ErrInvalidOptions)
	}
	for _, preference := range opts.Preferences {
		if preference == nil {
			return errors.Join(errors.New("nil proposer preference supplied"), client.ErrInvalidOptions)
		}
		// SSZ marshaling substitutes a zero value for a nil message, so without this
		// check the same input would silently submit zeroed preferences over SSZ and
		// be rejected by the node as a null message over JSON.
		if preference.Message == nil {
			return errors.Join(errors.New("nil proposer preference message supplied"), client.ErrInvalidOptions)
		}
	}

	requestPreferences := proposerPreferencesList(opts.Preferences)
	if requestPreferences == nil {
		requestPreferences = proposerPreferencesList{}
	}
	body, contentType, err := s.marshalRequestBody(ctx, requestPreferences)
	if err != nil {
		return err
	}

	// post() already rejects anything outside the 2xx family, and every other
	// submitter in this package takes that as success.  beacon-APIs documents
	// 200 as the only success code here, but insisting on exactly 200 is
	// asymmetric: a node that answers 202 or 204 on a successful
	// store-and-broadcast would produce an error that is neither an api.Error
	// 4xx nor ErrInvalidOptions, which multi reads as a provider fault and
	// deactivates the client for, on every submission.
	_, err = s.post(ctx,
		"/eth/v1/validator/proposer_preferences",
		"",
		&opts.Common,
		bytes.NewReader(body),
		contentType,
		map[string]string{"Eth-Consensus-Version": strings.ToLower(proposerPreferencesConsensusVersion.String())},
	)
	if err != nil {
		return errors.Join(errors.New("failed to submit proposer preferences"), err)
	}

	return nil
}

// proposerPreferencesConsensusVersion is the fork named in the
// Eth-Consensus-Version header on a proposer preferences submission.
//
// beacon-APIs describes that header as "the active consensus version to which
// the signed proposer preferences being submitted belongs".
// SignedProposerPreferences is unversioned in this client because the type is
// Gloas-only, but the header is about the chain's active fork, not the type: a
// node validating it against its own active version will reject submissions
// once the fork after Gloas activates, even though the message encoding has not
// changed.
//
// It is a constant rather than a value derived from the node's fork schedule
// because Gloas is the newest fork this library has a DataVersion for, so
// deriving it cannot produce a different answer today -- it would only add a
// genesis lookup to the submission path.  Whoever adds the DataVersion for the
// next fork has to replace this with that derivation: the active version then
// depends on the chain's clock, not on which forks the library knows.
const proposerPreferencesConsensusVersion = spec.DataVersionGloas
