// Copyright © 2020 - 2024 Attestant Limited.
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
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Spec provides the spec information of the chain.
func (s *Service) Spec(ctx context.Context,
	opts *api.SpecOpts,
) (
	*api.Response[map[string]any],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}

	s.specMutex.RLock()
	if s.spec != nil {
		defer s.specMutex.RUnlock()

		return &api.Response[map[string]any]{
			Data:     s.spec,
			Metadata: make(map[string]any),
		}, nil
	}
	s.specMutex.RUnlock()

	s.specMutex.Lock()
	defer s.specMutex.Unlock()
	if s.spec != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return &api.Response[map[string]any]{
			Data:     s.spec,
			Metadata: make(map[string]any),
		}, nil
	}

	// Up to us to fetch the information.
	endpoint := "/eth/v1/config/spec"
	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, false)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), map[string]any{})
	if err != nil {
		return nil, err
	}

	config := s.parseSpecsContainer(data)

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

	s.spec = config

	return &api.Response[map[string]any]{
		Data:     s.spec,
		Metadata: metadata,
	}, nil
}

func (s *Service) parseSpecsContainer(data map[string]any) map[string]any {
	config := make(map[string]any)

	for k, v := range data {
		// Handle complex structures
		if k == "BLOB_SCHEDULE" {
			if arrVal, isArr := v.([]map[string]any); isArr {
				values := make([]map[string]any, len(arrVal))
				for idx, val := range arrVal {
					values[idx] = s.parseSpecsContainer(val)
				}
			}
		}

		// Handle domains.
		if strings.HasPrefix(k, "DOMAIN_") {
			var byteVal []byte
			var err error

			if intVal, isInt := v.(int); isInt {
				byteVal = make([]byte, 4)
				binary.BigEndian.PutUint32(byteVal, uint32(intVal))
			} else if strVal, isStr := v.(string); isStr {
				byteVal, err = hex.DecodeString(strings.TrimPrefix(strVal, "0x"))
			}

			if err == nil {
				var domainType phase0.DomainType
				copy(domainType[:], byteVal)
				config[k] = domainType

				continue
			}
		}

		// Handle fork versions.
		if strings.HasSuffix(k, "_FORK_VERSION") {
			var byteVal []byte
			var err error

			if intVal, isInt := v.(int); isInt {
				byteVal = make([]byte, 4)
				binary.BigEndian.PutUint32(byteVal, uint32(intVal))
			} else if strVal, isStr := v.(string); isStr {
				byteVal, err = hex.DecodeString(strings.TrimPrefix(strVal, "0x"))
			}

			if err == nil {
				var version phase0.Version
				copy(version[:], byteVal)
				config[k] = version

				continue
			}
		}

		// Handle hex strings.
		if strVal, isStr := v.(string); isStr && strings.HasPrefix(strVal, "0x") {
			byteVal, err := hex.DecodeString(strings.TrimPrefix(strVal, "0x"))
			if err == nil {
				config[k] = byteVal

				continue
			}
		}

		// Handle times.
		if strings.HasSuffix(k, "_TIME") {
			var int64Val int64
			if intVal, isInt := v.(int); isInt {
				int64Val = int64(intVal)
			} else if strVal, isStr := v.(string); isStr {
				intVal, err := strconv.ParseInt(strVal, 10, 64)
				if err == nil {
					int64Val = intVal
				}
			}

			config[k] = time.Unix(int64Val, 0)

			continue
		}

		// Handle durations.
		if strings.HasPrefix(k, "SECONDS_PER_") || k == "GENESIS_DELAY" {
			var int64Val int64
			if intVal, isInt := v.(int); isInt {
				int64Val = int64(intVal)
			} else if strVal, isStr := v.(string); isStr {
				intVal, err := strconv.ParseInt(strVal, 10, 64)
				if err == nil {
					int64Val = intVal
				}
			}

			config[k] = time.Duration(int64Val) * time.Second

			continue
		}

		// Handle integers.
		if intVal, isInt := v.(int); isInt {
			config[k] = uint64(intVal)

			continue
		} else if strVal, isStr := v.(string); isStr {
			intVal, err := strconv.ParseUint(strVal, 10, 64)
			if err == nil {
				config[k] = intVal

				continue
			}
		}

		// Assume string.
		if strVal, isStr := v.(string); isStr {
			config[k] = strVal

			continue
		}
	}

	return config
}
