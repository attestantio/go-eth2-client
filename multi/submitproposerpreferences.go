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

package multi

import (
	"context"
	"errors"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
)

// SubmitProposerPreferences submits signed proposer preferences.
func (s *Service) SubmitProposerPreferences(ctx context.Context, opts *api.SubmitProposerPreferencesOpts) error {
	_, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		if err := client.(consensusclient.ProposerPreferencesSubmitter).SubmitProposerPreferences(ctx, opts); err != nil {
			return nil, err
		}

		return true, nil
	}, func(_ context.Context, _ consensusclient.Service, err error) (bool, error) {
		// Invalid input is invalid for every client, so neither fail over nor
		// deactivate a healthy client because the caller supplied bad
		// preferences.  ErrNoOptions is a separate sentinel from
		// ErrInvalidOptions rather than joined with it, so both have to be
		// named here: a nil opts would otherwise deactivate every provider in
		// turn and leave the multi client itself inactive.
		callerError := errors.Is(err, consensusclient.ErrInvalidOptions) ||
			errors.Is(err, consensusclient.ErrNoOptions)

		return !callerError, err
	})

	return err
}
