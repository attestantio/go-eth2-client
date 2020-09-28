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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type specJSON struct {
	Data map[string]string `json:"data"`
}

// Spec provides the spec information of the chain.
func (s *Service) Spec(ctx context.Context) (map[string]interface{}, error) {
	if s.spec == nil {
		respBodyReader, err := s.get(ctx, "/eth/v1/config/spec")
		if err != nil {
			return nil, errors.Wrap(err, "failed to request spec")
		}

		specReader, err := s.lhToSpec(ctx, respBodyReader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert teku response to spec response")
		}

		var specJSON specJSON
		if err := json.NewDecoder(specReader).Decode(&specJSON); err != nil {
			return nil, errors.Wrap(err, "failed to parse spec")
		}

		spec := make(map[string]interface{})
		for k, v := range specJSON.Data {
			// Handle hex strings.
			if strings.HasPrefix(v, "0x") {
				byteVal, err := hex.DecodeString(strings.TrimPrefix(v, "0x"))
				if err == nil {
					spec[k] = byteVal
					continue
				}
			}

			// Handle durations.
			if strings.HasPrefix(k, "SECONDS_PER_") {
				intVal, err := strconv.ParseUint(v, 10, 64)
				if err == nil && intVal != 0 {
					spec[k] = time.Duration(intVal) * time.Second
					continue
				}
			}

			// Handle integers.
			if v == "0" {
				spec[k] = uint64(0)
				continue
			}
			intVal, err := strconv.ParseUint(v, 10, 64)
			if err == nil && intVal != 0 {
				spec[k] = intVal
				continue
			}

			// Unknown format.
			return nil, fmt.Errorf("invalid format of value %s", k)
		}
		s.spec = spec
	}
	return s.spec, nil
}
