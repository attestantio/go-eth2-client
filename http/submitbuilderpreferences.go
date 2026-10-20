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

// maxBuilderEntries is MAX_BUILDER_ENTRIES, the number of builders per proposer.
const maxBuilderEntries = 64

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
// to each entry's builder.  If some builders reject their entries the returned
// error wraps an *api.IndexedError naming them by index; any other failure,
// including a node that does not serve the endpoint, wraps only an *api.Error.
func (s *Service) SubmitBuilderPreferences(ctx context.Context, opts *api.SubmitBuilderPreferencesOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	// The limit is MAX_BUILDER_ENTRIES times the proposer lookahead length,
	// compared by division so the product cannot overflow.
	if (uint64(len(opts.Entries))+maxBuilderEntries-1)/maxBuilderEntries > s.proposerPreferencesLimit(ctx) {
		return errors.Join(errors.New("too many builder preferences entries"), client.ErrInvalidOptions)
	}
	for i, entry := range opts.Entries {
		if entry == nil {
			return errors.Join(errors.New("nil builder preferences entry supplied"), client.ErrInvalidOptions)
		}
		if len(entry.URL) == 0 {
			return errors.Join(fmt.Errorf("builder preferences entry %d has no URL", i), client.ErrInvalidOptions)
		}
		// SSZ marshaling substitutes zero values for nil, which would submit an
		// unsigned authorization.
		if entry.Auth == nil || entry.Auth.Message == nil {
			return errors.Join(fmt.Errorf("builder preferences entry %d has no authorization", i), client.ErrInvalidOptions)
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
			indexedErr := &api.IndexedError{}
			if json.Unmarshal(apiErr.Data, indexedErr) == nil && len(indexedErr.Failures) > 0 {
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
