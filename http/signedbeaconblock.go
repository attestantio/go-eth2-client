// Copyright © 2020, 2021 Attestant Limited.
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
	"fmt"
	"io"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0SignedBeaconBlockJSON struct {
	Data *phase0.SignedBeaconBlock `json:"data"`
}

type altairSignedBeaconBlockJSON struct {
	Data *altair.SignedBeaconBlock `json:"data"`
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
// N.B if a signed beacon block for the block ID is not available this will return nil without an error.
func (s *Service) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	if s.supportsV2BeaconBlocks {
		return s.signedBeaconBlockV2(ctx, blockID)
	}
	return s.signedBeaconBlockV1(ctx, blockID)
}

// signedBeaconBlockV1 fetches a signed beacon block from the V1 endpoint.
func (s *Service) signedBeaconBlockV1(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/blocks/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp phase0SignedBeaconBlockJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon block")
	}

	return &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionPhase0,
		Phase0:  resp.Data,
	}, nil
}

// signedBeaconBlockV2 fetches a signed beacon block from the V2 endpoint.
func (s *Service) signedBeaconBlockV2(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v2/beacon/blocks/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var dataBodyReader bytes.Buffer
	metadataReader := io.TeeReader(respBodyReader, &dataBodyReader)
	var metadata responseMetadata
	if err := json.NewDecoder(metadataReader).Decode(&metadata); err != nil {
		return nil, errors.Wrap(err, "failed to parse response")
	}
	res := &spec.VersionedSignedBeaconBlock{
		Version: metadata.Version,
	}

	switch metadata.Version {
	case spec.DataVersionPhase0:
		var resp phase0SignedBeaconBlockJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase 0 signed beacon block")
		}
		res.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairSignedBeaconBlockJSON
		if err := json.NewDecoder(&dataBodyReader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair signed beacon block")
		}
		res.Altair = resp.Data
	}

	return res, nil
}
