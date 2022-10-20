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

package http

import (
	"context"
	"strings"
)

// NodeClient provides the client for the node.
func (s *Service) NodeClient(ctx context.Context) (string, error) {
	nodeVersion, err := s.NodeVersion(ctx)
	if err != nil {
		return "", err
	}

	nodeVersion = strings.ToLower(nodeVersion)

	switch {
	case strings.HasPrefix(nodeVersion, "lighthouse"):
		return "lighthouse", nil
	case strings.HasPrefix(nodeVersion, "nimbus"):
		return "nimbus", nil
	case strings.HasPrefix(nodeVersion, "prysm"):
		return "prysm", nil
	case strings.HasPrefix(nodeVersion, "teku"):
		return "teku", nil
	default:
		return nodeVersion, nil
	}
}
