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
	"maps"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// PayloadAttestationPool obtains the payload attestation pool for the given options.
func (s *Service) PayloadAttestationPool(ctx context.Context,
	opts *api.PayloadAttestationPoolOpts,
) (
	*api.Response[[]*spec.VersionedPayloadAttestation],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := "/eth/v1/beacon/pool/payload_attestations"

	query := ""
	if opts.Slot != nil {
		query = fmt.Sprintf("slot=%d", *opts.Slot)
	}

	// The endpoint advertises SSZ, but as a bare SSZ list with no container to
	// decode into, so JSON is requested here.
	httpResponse, err := s.getWithResponseLimit(ctx, endpoint, query, &opts.Common, false, maxEPBSResponseSize)
	if err != nil {
		return nil, err
	}

	response, err := payloadAttestationPoolFromResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	if err := verifyPayloadAttestationPool(opts.Slot, response.Data); err != nil {
		return nil, errors.Join(err, client.ErrInconsistentResult)
	}

	return response, nil
}

// payloadAttestationPoolFromResponse decodes a fetched payload attestation
// pool response.  It is separate from the fetch so that decoding can be
// exercised without a beacon node.
func payloadAttestationPoolFromResponse(httpResponse *httpResponse) (
	*api.Response[[]*spec.VersionedPayloadAttestation],
	error,
) {
	if httpResponse.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("payload attestations not available for version %s", httpResponse.consensusVersion)
	}

	metadata := metadataFromHeaders(httpResponse.headers)

	var attestations []*gloas.PayloadAttestation

	switch httpResponse.contentType {
	case ContentTypeJSON:
		decoded, jsonMetadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*gloas.PayloadAttestation{})
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s payload attestations", httpResponse.consensusVersion),
				err,
			)
		}

		attestations = decoded

		maps.Copy(metadata, jsonMetadata)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	data := make([]*spec.VersionedPayloadAttestation, len(attestations))
	for i := range attestations {
		if attestations[i] == nil {
			return nil, errors.New("nil payload attestation in response")
		}

		data[i] = &spec.VersionedPayloadAttestation{
			Version: httpResponse.consensusVersion,
			Gloas:   attestations[i],
		}
	}

	return &api.Response[[]*spec.VersionedPayloadAttestation]{
		Data:     data,
		Metadata: metadata,
	}, nil
}

// verifyPayloadAttestationPool confirms that every attestation returned by the
// node holds an attestation to read, and is for the requested slot when one
// was requested.
func verifyPayloadAttestationPool(slot *phase0.Slot, attestations []*spec.VersionedPayloadAttestation) error {
	for _, attestation := range attestations {
		// Data rejects a missing version arm as well as a missing datum, and
		// is worth calling even with no slot to filter on: requesting the
		// whole pool waives the filter, not the contents of what comes back.
		data, err := attestation.Data()
		if err != nil {
			return err
		}

		if slot != nil && data.Slot != *slot {
			return fmt.Errorf("payload attestation for slot %d; expected %d", data.Slot, *slot)
		}
	}

	return nil
}
