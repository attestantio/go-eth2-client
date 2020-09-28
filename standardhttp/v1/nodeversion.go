// Copyright © 2020 Attestant Limited.
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

package v1

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type nodeVersionJSON struct {
	Data *nodeVersionDataJSON `json:"data"`
}

type nodeVersionDataJSON struct {
	Version string `json:"version"`
}

// NodeVersion provides the version information of the node.
func (s *Service) NodeVersion(ctx context.Context) (string, error) {
	if s.nodeVersion == "" {
		respBodyReader, err := s.get(ctx, "/eth/v1/node/version")
		if err != nil {
			return "", errors.Wrap(err, "failed to request node version")
		}

		var resp nodeVersionJSON
		if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
			return "", errors.Wrap(err, "failed to parse node version")
		}
		s.nodeVersion = resp.Data.Version
	}

	return s.nodeVersion, nil
}
