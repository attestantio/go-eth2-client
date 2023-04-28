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
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type specJSON struct {
	Data map[string]string `json:"data"`
}

// Spec provides the spec information of the chain.
func (s *Service) Spec(ctx context.Context) (map[string]interface{}, error) {
	s.specMutex.RLock()
	if s.spec != nil {
		defer s.specMutex.RUnlock()
		return s.spec, nil
	}
	s.specMutex.RUnlock()

	s.specMutex.Lock()
	defer s.specMutex.Unlock()
	if s.spec != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return s.spec, nil
	}

	// Up to us to fetch the information.
	respBodyReader, err := s.get(ctx, "/eth/v1/config/spec")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request spec")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain spec")
	}

	var specJSON specJSON
	if err := json.NewDecoder(respBodyReader).Decode(&specJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse spec")
	}

	config := make(map[string]interface{})
	for k, v := range specJSON.Data {
		// Handle domains.
		if strings.HasPrefix(k, "DOMAIN_") {
			byteVal, err := hex.DecodeString(strings.TrimPrefix(v, "0x"))
			if err == nil {
				var domainType phase0.DomainType
				copy(domainType[:], byteVal)
				config[k] = domainType
				continue
			}
		}

		// Handle fork versions.
		if strings.HasSuffix(k, "_FORK_VERSION") {
			byteVal, err := hex.DecodeString(strings.TrimPrefix(v, "0x"))
			if err == nil {
				var version phase0.Version
				copy(version[:], byteVal)
				config[k] = version
				continue
			}
		}

		// Handle hex strings.
		if strings.HasPrefix(v, "0x") {
			byteVal, err := hex.DecodeString(strings.TrimPrefix(v, "0x"))
			if err == nil {
				config[k] = byteVal
				continue
			}
		}

		// Handle times.
		if strings.HasSuffix(k, "_TIME") {
			intVal, err := strconv.ParseInt(v, 10, 64)
			if err == nil && intVal != 0 {
				config[k] = time.Unix(intVal, 0)
				continue
			}
		}

		// Handle durations.
		if strings.HasPrefix(k, "SECONDS_PER_") || k == "GENESIS_DELAY" {
			intVal, err := strconv.ParseUint(v, 10, 64)
			if err == nil && intVal != 0 {
				config[k] = time.Duration(intVal) * time.Second
				continue
			}
		}

		// Handle integers.
		if v == "0" {
			config[k] = uint64(0)
			continue
		}
		intVal, err := strconv.ParseUint(v, 10, 64)
		if err == nil && intVal != 0 {
			config[k] = intVal
			continue
		}

		// Assume string.
		config[k] = v
	}

	// The application mask domain type is not provided by all nodes, so add it here if not present.
	if _, exists := config["DOMAIN_APPLICATION_MASK"]; !exists {
		config["DOMAIN_APPLICATION_MASK"] = phase0.DomainType{0x00, 0x00, 0x00, 0x01}
	}
	// The BLS to execution change domain type is not provided by all nodes, so add it here if not present.
	if _, exists := config["DOMAIN_BLS_TO_EXECUTION_CHANGE"]; !exists {
		config["DOMAIN_BLS_TO_EXECUTION_CHANGE"] = phase0.DomainType{0x0a, 0x00, 0x00, 0x00}
	}
	// The builder application domain type is not officially part of the spec, so add it here if not present.
	if _, exists := config["DOMAIN_APPLICATION_BUILDER"]; !exists {
		config["DOMAIN_APPLICATION_BUILDER"] = phase0.DomainType{0x00, 0x00, 0x00, 0x01}
	}
	// The blob sidecar domain type is not provided by all nodes, so add it here if not present.
	if _, exists := config["DOMAIN_BLOB_SIDECAR"]; !exists {
		config["DOMAIN_BLOB_SIDECAR"] = phase0.DomainType{0x0b, 0x00, 0x00, 0x00}
	}

	s.spec = config
	return s.spec, nil
}
