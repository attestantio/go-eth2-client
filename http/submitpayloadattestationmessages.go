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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// SubmitPayloadAttestationMessages submits payload attestation messages.
func (s *Service) SubmitPayloadAttestationMessages(ctx context.Context,
	opts *api.SubmitPayloadAttestationMessagesOpts,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if len(opts.Messages) == 0 {
		return errors.Join(errors.New("no payload attestation messages supplied"), client.ErrInvalidOptions)
	}

	unversioned, version, err := createUnversionedPayloadAttestationMessages(opts.Messages)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(unversioned)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v1/beacon/pool/payload_attestations"

	// The endpoint requires the consensus version of the messages being
	// submitted, and rejects the request outright without it.
	headers := map[string]string{
		"Eth-Consensus-Version": strings.ToLower(version.String()),
	}

	if _, err := s.post(ctx,
		endpoint,
		"",
		&opts.Common,
		bytes.NewReader(reqBody),
		ContentTypeJSON,
		headers,
	); err != nil {
		return errors.Join(errors.New("failed to submit payload attestation messages"), err)
	}

	return nil
}

// createUnversionedPayloadAttestationMessages strips the version wrapper from
// the messages, returning the bare messages to place on the wire along with
// the version they share.
func createUnversionedPayloadAttestationMessages(messages []*spec.VersionedPayloadAttestationMessage) (
	[]any,
	spec.DataVersion,
	error,
) {
	var (
		version     spec.DataVersion
		unversioned = make([]any, 0, len(messages))
	)

	for i := range messages {
		if messages[i] == nil {
			return nil, spec.DataVersionUnknown,
				errors.Join(errors.New("nil payload attestation message supplied"), client.ErrInvalidOptions)
		}

		// Ensure consistent versioning.
		if version == spec.DataVersionUnknown {
			version = messages[i].Version
		} else if version != messages[i].Version {
			return nil, spec.DataVersionUnknown, errors.Join(
				errors.New("payload attestation messages must all be of the same version"),
				client.ErrInvalidOptions,
			)
		}

		switch messages[i].Version {
		case spec.DataVersionGloas:
			if messages[i].Gloas == nil {
				// Appending an unchecked arm here would put a null in the
				// request body rather than a message.
				return nil, spec.DataVersionUnknown,
					errors.Join(errors.New("no gloas payload attestation message supplied"), client.ErrInvalidOptions)
			}

			unversioned = append(unversioned, messages[i].Gloas)
		default:
			return nil, spec.DataVersionUnknown, errors.Join(
				fmt.Errorf("unsupported payload attestation message version %s", messages[i].Version),
				client.ErrInvalidOptions,
			)
		}
	}

	return unversioned, version, nil
}
