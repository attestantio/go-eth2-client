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
	"net/http"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

// PayloadAttestationData obtains payload attestation data for the given options.
//
// A node only produces this for its current slot; if it has seen no block for
// that slot it answers 204, which surfaces as client.ErrNoPayloadAttestationData
// and means the validator must not cast a payload attestation.
func (s *Service) PayloadAttestationData(ctx context.Context,
	opts *api.PayloadAttestationDataOpts,
) (
	*api.Response[*spec.VersionedPayloadAttestationData],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/payload_attestation_data/%d", opts.Slot)
	query := ""

	httpResponse, err := s.getWithResponseLimit(ctx, endpoint, query, &opts.Common, true, maxEPBSResponseSize)
	if err != nil {
		return nil, err
	}

	response, err := payloadAttestationDataFromResponse(httpResponse)
	if err != nil {
		return nil, err
	}

	slot, err := response.Data.Slot()
	if err != nil {
		return nil, err
	}

	if slot != opts.Slot {
		return nil, errors.Join(
			fmt.Errorf("payload attestation data for slot %d; expected %d", slot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}

	return response, nil
}

// payloadAttestationDataFromResponse decodes a fetched payload attestation
// data response.  It is separate from the fetch so that the paths a live node
// will not produce on demand — a 204, or an SSZ body — remain testable.
func payloadAttestationDataFromResponse(httpResponse *httpResponse) (
	*api.Response[*spec.VersionedPayloadAttestationData],
	error,
) {
	if httpResponse.statusCode == http.StatusNoContent {
		// The node has seen no block for this slot, so there is nothing to
		// attest to.  This is a defined outcome of the endpoint, not a
		// failure, but it must not be mistaken for a signable datum.
		return nil, client.ErrNoPayloadAttestationData
	}

	if httpResponse.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("payload attestation data not available for version %s", httpResponse.consensusVersion)
	}

	// Payload attestation data is a gloas-onwards container, so the wire bytes
	// always parse into *gloas.PayloadAttestationData.
	data := &gloas.PayloadAttestationData{}
	metadata := metadataFromHeaders(httpResponse.headers)

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		// The container is a fixed-size root, slot and two flags, with no
		// preset-derived sizes, so the generated codec decodes it correctly at
		// any preset and no dynamic SSZ codec is required.
		if err := data.UnmarshalSSZ(httpResponse.body); err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s payload attestation data", httpResponse.consensusVersion),
				err,
			)
		}
	case ContentTypeJSON:
		// The decoder only writes the keys it finds, so seeding it with a nil
		// pointer rather than the empty datum above is what makes a body with
		// no data key in it distinguishable from one carrying a zero-valued
		// datum: only the former leaves the seed as it was.
		decoded, jsonMetadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), (*gloas.PayloadAttestationData)(nil))
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s payload attestation data", httpResponse.consensusVersion),
				err,
			)
		}

		if decoded == nil {
			// Neither an absent data key nor an explicit null leaves anything
			// to attest to.  Returning a success here would hand back a datum
			// whose accessors all answer zero without error, which a caller
			// cannot tell apart from a real vote for the zero root.
			return nil, fmt.Errorf("no %s payload attestation data in response", httpResponse.consensusVersion)
		}

		data = decoded

		maps.Copy(metadata, jsonMetadata)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	return &api.Response[*spec.VersionedPayloadAttestationData]{
		Data: &spec.VersionedPayloadAttestationData{
			Version: httpResponse.consensusVersion,
			Gloas:   data,
		},
		Metadata: metadata,
	}, nil
}
