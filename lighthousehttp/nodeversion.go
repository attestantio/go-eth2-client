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

package lighthousehttp

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// NodeVersion returns a free-text string with the node version.
func (s *Service) NodeVersion(ctx context.Context) (string, error) {
	respBodyReader, cancel, err := s.get(ctx, "/node/version")
	if err != nil {
		return "", errors.Wrap(err, "failed to obtain node version")
	}
	defer cancel()

	version, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return "", errors.Wrap(err, "failed to read node version")
	}
	return strings.TrimSuffix(strings.TrimPrefix(string(version), `"`), `"`), nil
}
