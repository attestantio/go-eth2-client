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

// proposerPreferencesLimit returns the maximum number of proposer preferences
// that may be submitted at once -- the proposer lookahead length,
// (MIN_SEED_LOOKAHEAD + 1) * SLOTS_PER_EPOCH -- and whether it could be derived
// at all.
//
// The node is the authority on what it accepts; this is only a client-side
// sanity bound.  So when the chain's own values cannot be derived the bound is
// skipped rather than replaced with the mainnet preset.  Mainnet's 64 fails
// open on a chain with a shorter lookahead, where the extra entries are simply
// rejected by the node, but on one with a longer lookahead -- SLOTS_PER_EPOCH
// 64 gives a real limit of 128 -- it would reject a 65-to-128 entry submission
// the node would have accepted, with ErrInvalidOptions and no network call.
//
// The derived value is cached because the comment this replaces assumed Spec()
// already was.  It is not: Spec() assigns s.spec only on success, so against a
// node that does not serve /eth/v1/config/spec the fetch is repeated on every
// submission, and clearStaticValues nils it every 5 minutes, while preferences
// are submitted about once per epoch -- so in practice most submissions were
// the call that refilled the cache.  Spec() holds specMutex across the whole
// round trip, so that refill also blocked every other spec consumer while a
// time-sensitive submission waited on it.
func (s *Service) proposerPreferencesLimit(ctx context.Context) (uint64, bool) {
	if cached := s.proposerPreferencesLimitCache.Load(); cached != 0 {
		return cached, true
	}

	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return 0, false
	}

	minSeedLookahead, ok := response.Data["MIN_SEED_LOOKAHEAD"].(uint64)
	if !ok {
		return 0, false
	}
	slotsPerEpoch, ok := response.Data["SLOTS_PER_EPOCH"].(uint64)
	if !ok {
		return 0, false
	}
	if minSeedLookahead == ^uint64(0) || slotsPerEpoch > ^uint64(0)/(minSeedLookahead+1) {
		return 0, false
	}

	limit := (minSeedLookahead + 1) * slotsPerEpoch
	if limit == 0 {
		return 0, false
	}

	s.proposerPreferencesLimitCache.Store(limit)

	return limit, true
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

	// Guarded on a non-empty list: Go evaluates the right-hand operand of >
	// unconditionally, so without this the common "nothing to publish" call
	// pays the spec lookup to check a bound it cannot exceed.
	if len(opts.Preferences) > 0 {
		if limit, known := s.proposerPreferencesLimit(ctx); known && uint64(len(opts.Preferences)) > limit {
			return errors.Join(errors.New("too many proposer preferences"), client.ErrInvalidOptions)
		}
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

	response, err := s.post(ctx,
		"/eth/v1/validator/proposer_preferences",
		"",
		&opts.Common,
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
