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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type phase0SignedBeaconBlockJSON struct {
	Data *phase0.SignedBeaconBlock `json:"data"`
}

type altairSignedBeaconBlockJSON struct {
	Data *altair.SignedBeaconBlock `json:"data"`
}

type bellatrixSignedBeaconBlockJSON struct {
	Data *bellatrix.SignedBeaconBlock `json:"data"`
}

type capellaSignedBeaconBlockJSON struct {
	Data *capella.SignedBeaconBlock `json:"data"`
}

type denebSignedBeaconBlockJSON struct {
	Data *deneb.SignedBeaconBlock `json:"data"`
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
// N.B if a signed beacon block for the block ID is not available this will return nil without an error.
func (s *Service) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.VersionedSignedBeaconBlock, error) {
	res, err := s.get2(ctx, fmt.Sprintf("/eth/v2/beacon/blocks/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}
	defer res.reader.Close()
	if res.statusCode == http.StatusNotFound {
		return nil, nil
	}

	switch res.contentType {
	case ContentTypeSSZ:
		return s.signedBeaconBlockFromSSZ(res)
	case ContentTypeJSON:
		return s.signedBeaconBlockFromJSON(res)
	default:
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}
}

func (s *Service) signedBeaconBlockFromSSZ(res *httpResponse) (*spec.VersionedSignedBeaconBlock, error) {
	block := &spec.VersionedSignedBeaconBlock{
		Version: res.consensusVersion,
	}

	data, err := io.ReadAll(res.reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read block")
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		block.Phase0 = &phase0.SignedBeaconBlock{}
		if err := block.Phase0.UnmarshalSSZ(data); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 signed beacon block")
		}
	case spec.DataVersionAltair:
		block.Altair = &altair.SignedBeaconBlock{}
		if err := block.Altair.UnmarshalSSZ(data); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair signed beacon block")
		}
	case spec.DataVersionBellatrix:
		block.Bellatrix = &bellatrix.SignedBeaconBlock{}
		if err := block.Bellatrix.UnmarshalSSZ(data); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix signed beacon block")
		}
	case spec.DataVersionCapella:
		block.Capella = &capella.SignedBeaconBlock{}
		if err := block.Capella.UnmarshalSSZ(data); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella signed beacon block")
		}
	case spec.DataVersionDeneb:
		block.Deneb = &deneb.SignedBeaconBlock{}
		if err := block.Deneb.UnmarshalSSZ(data); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb signed beacon block")
		}
	default:
		return nil, fmt.Errorf("unhandled block version %s", res.consensusVersion)
	}

	return block, nil
}

func (s *Service) signedBeaconBlockFromJSON(res *httpResponse) (*spec.VersionedSignedBeaconBlock, error) {
	block := &spec.VersionedSignedBeaconBlock{
		Version: res.consensusVersion,
	}

	switch block.Version {
	case spec.DataVersionPhase0:
		var resp phase0SignedBeaconBlockJSON
		if err := json.NewDecoder(res.reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse phase 0 signed beacon block")
		}
		block.Phase0 = resp.Data
	case spec.DataVersionAltair:
		var resp altairSignedBeaconBlockJSON
		if err := json.NewDecoder(res.reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse altair signed beacon block")
		}
		block.Altair = resp.Data
	case spec.DataVersionBellatrix:
		var resp bellatrixSignedBeaconBlockJSON
		if err := json.NewDecoder(res.reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse bellatrix signed beacon block")
		}
		block.Bellatrix = resp.Data
	case spec.DataVersionCapella:
		var resp capellaSignedBeaconBlockJSON
		if err := json.NewDecoder(res.reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse capella signed beacon block")
		}
		block.Capella = resp.Data
	case spec.DataVersionDeneb:
		var resp denebSignedBeaconBlockJSON
		if err := json.NewDecoder(res.reader).Decode(&resp); err != nil {
			return nil, errors.Wrap(err, "failed to parse deneb signed beacon block")
		}
		block.Deneb = resp.Data
	default:
		return nil, fmt.Errorf("unhandled block version %s", res.consensusVersion)
	}

	return block, nil
}
