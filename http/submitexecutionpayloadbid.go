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

// SubmitExecutionPayloadBid submits a signed execution payload bid for
// gossip broadcast.
func (s *Service) SubmitExecutionPayloadBid(ctx context.Context,
	opts *api.SubmitExecutionPayloadBidOpts,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if opts.SignedExecutionPayloadBid == nil {
		return errors.Join(errors.New("no bid supplied"), client.ErrInvalidOptions)
	}

	versioned := opts.SignedExecutionPayloadBid

	var bid any

	switch versioned.Version {
	case spec.DataVersionGloas:
		if versioned.Gloas == nil {
			return errors.Join(errors.New("no gloas bid supplied"), client.ErrInvalidOptions)
		}

		bid = versioned.Gloas
	default:
		return errors.Join(
			fmt.Errorf("unsupported bid version %s", versioned.Version),
			client.ErrInvalidOptions,
		)
	}

	return s.postExecutionPayloadBid(ctx, &opts.Common, versioned.Version, bid)
}

// postExecutionPayloadBid marshals the bid to the negotiated content type
// (SSZ unless JSON is enforced) and performs the POST.
func (s *Service) postExecutionPayloadBid(ctx context.Context,
	common *api.CommonOpts,
	consensusVersion spec.DataVersion,
	bid any,
) error {
	var (
		body        []byte
		contentType ContentType
		err         error
	)

	if s.enforceJSON {
		contentType = ContentTypeJSON

		body, err = json.Marshal(bid)
		if err != nil {
			return errors.Join(errors.New("failed to marshal JSON"), err)
		}
	} else {
		contentType = ContentTypeSSZ

		ds, dsErr := s.dynSSZForRequest(ctx)
		if dsErr != nil {
			return dsErr
		}

		body, err = ds.MarshalSSZ(bid)
		if err != nil {
			return errors.Join(errors.New("failed to marshal SSZ"), err)
		}
	}

	endpoint := "/eth/v1/beacon/execution_payload_bids"

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(consensusVersion.String())

	if _, err := s.post(ctx, endpoint, "", common, bytes.NewBuffer(body), contentType, headers); err != nil {
		return errors.Join(errors.New("failed to submit execution payload bid"), err)
	}

	return nil
}
