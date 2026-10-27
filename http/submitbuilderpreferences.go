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
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	nethttp "net/http"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

type builderPreferencesList []*gloas.BuilderPreferencesEntry

// MAX_BUILDER_ENTRIES and MAX_BUILDER_URL_SIZE from beacon-APIs
// types/gloas/builder_entry.yaml, and MAX_BUILDER_AUTH_DATA_SIZE from
// builder-specs gloas/builder.md.  beacon-APIs also makes a zero-length url
// (types/gloas/builder_preferences_entry.yaml) and zero-length auth data
// (types/gloas/request_auth.yaml) invalid.
const (
	maxBuilderEntries      = 64
	maxBuilderURLSize      = 2048
	maxBuilderAuthDataSize = 4096
)

// MarshalSSZ encodes the list as an SSZ list of variable-size elements: a
// 4-byte offset per element, then the elements' own encodings.
func (p builderPreferencesList) MarshalSSZ() ([]byte, error) {
	size := 4 * len(p)
	for _, entry := range p {
		size += entry.SizeSSZ()
	}

	body := make([]byte, 0, size)
	offset := 4 * len(p)
	for _, entry := range p {
		body = binary.LittleEndian.AppendUint32(body, uint32(offset))
		offset += entry.SizeSSZ()
	}
	for _, entry := range p {
		var err error
		if body, err = entry.MarshalSSZTo(body); err != nil {
			return nil, err
		}
	}

	return body, nil
}

// SubmitBuilderPreferences submits builder preferences, which the node forwards
// to each entry's builder.  When the node answers 400 with failures, the
// returned error wraps an *api.IndexedError as well as the *api.Error.  It holds
// only the failures whose entry decodes with an index within opts.Entries, so
// every Index is safe to use on opts.Entries; if none qualifies there is no
// *api.IndexedError, and the raw body is always on the *api.Error.  A node that
// does not serve the endpoint returns an *api.Error with status 404 and no
// *api.IndexedError.
func (s *Service) SubmitBuilderPreferences(ctx context.Context, opts *api.SubmitBuilderPreferencesOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	// The limit is MAX_BUILDER_ENTRIES times the proposer lookahead length, which
	// proposerPreferencesLimit returns, compared by division so the product
	// cannot overflow.
	if (uint64(len(opts.Entries))+maxBuilderEntries-1)/maxBuilderEntries > s.proposerPreferencesLimit(ctx) {
		return errors.Join(errors.New("too many builder preferences entries"), client.ErrInvalidOptions)
	}
	// Invalid entries fail as ErrInvalidOptions so that multi does not fail over.
	for i, entry := range opts.Entries {
		if entry == nil {
			return errors.Join(fmt.Errorf("builder preferences entry %d is nil", i), client.ErrInvalidOptions)
		}
		if len(entry.URL) == 0 || len(entry.URL) > maxBuilderURLSize {
			return errors.Join(
				fmt.Errorf("builder preferences entry %d URL length %d is invalid", i, len(entry.URL)),
				client.ErrInvalidOptions,
			)
		}
		// SSZ marshaling substitutes zero values for nil, which would submit an
		// unsigned authorization.
		if entry.Auth == nil || entry.Auth.Message == nil {
			return errors.Join(fmt.Errorf("builder preferences entry %d has no authorization", i), client.ErrInvalidOptions)
		}
		if len(entry.Auth.Message.Data) == 0 || len(entry.Auth.Message.Data) > maxBuilderAuthDataSize {
			return errors.Join(
				fmt.Errorf("builder preferences entry %d authorization data length %d is invalid", i, len(entry.Auth.Message.Data)),
				client.ErrInvalidOptions,
			)
		}
	}

	requestEntries := builderPreferencesList(opts.Entries)
	if requestEntries == nil {
		requestEntries = builderPreferencesList{}
	}
	body, contentType, err := s.marshalRequestBody(ctx, requestEntries)
	if err != nil {
		return err
	}

	response, err := s.post(ctx,
		"/eth/v1/validator/builder_preferences",
		"",
		&opts.Common,
		bytes.NewReader(body),
		contentType,
		map[string]string{"Eth-Consensus-Version": strings.ToLower(spec.DataVersionGloas.String())},
	)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.StatusCode == nethttp.StatusBadRequest {
			if indexedErr := decodeIndexedError(apiErr.Data, len(opts.Entries)); indexedErr != nil {
				return errors.Join(errors.New("failed to submit builder preferences"), indexedErr, err)
			}
		}

		return errors.Join(errors.New("failed to submit builder preferences"), err)
	}
	if response.statusCode != nethttp.StatusOK {
		return errors.Join(
			errors.New("failed to submit builder preferences"),
			fmt.Errorf("unexpected status code %d", response.statusCode),
		)
	}

	return nil
}

// decodeIndexedError decodes an IndexedErrorMessage body, keeping only the
// failures with an index within the submitted list, so that callers can index
// their entries with it.  It returns nil if no failure is usable.
func decodeIndexedError(data []byte, entryCount int) *api.IndexedError {
	// A body that is not an IndexedErrorMessage leaves no failures.
	var decoded api.IndexedError
	_ = json.Unmarshal(data, &decoded)

	indexedErr := &api.IndexedError{Message: decoded.Message}
	for _, failure := range decoded.Failures {
		if failure.Index >= 0 && failure.Index < entryCount {
			indexedErr.Failures = append(indexedErr.Failures, failure)
		}
	}
	if len(indexedErr.Failures) == 0 {
		return nil
	}

	return indexedErr
}
