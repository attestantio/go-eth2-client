// Copyright © 2021 Attestant Limited.
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

package v2

import (
	"bytes"
	"context"
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pkg/errors"
)

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Service) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *spec.SignedVoluntaryExit) error {
	specJSON, err := json.Marshal(voluntaryExit)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	_, err = s.post(ctx, "/eth/v1/beacon/pool/voluntary_exits", bytes.NewBuffer(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to submit voluntary exit")
	}

	return nil
}
