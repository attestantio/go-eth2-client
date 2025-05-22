// Copyright © 2023 Attestant Limited.
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

package eip7732

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionPayloadHeader) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&executionPayloadHeaderJSON{
		ParentBlockHash:        fmt.Sprintf("%#x", e.ParentBlockHash),
		ParentBlockRoot:        fmt.Sprintf("%#x", e.ParentBlockRoot),
		BlockHash:              fmt.Sprintf("%#x", e.BlockHash),
		GasLimit:               fmt.Sprintf("%d", e.GasLimit),
		BuilderIndex:           fmt.Sprintf("%d", e.BuilderIndex),
		Slot:                   fmt.Sprintf("%d", e.Slot),
		Value:                  fmt.Sprintf("%d", e.Value),
		BlobKZGCommitmentsRoot: fmt.Sprintf("%#x", e.BlobKZGCommitmentsRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionPayloadHeader) UnmarshalYAML(input []byte) error {
	var data executionPayloadHeaderJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(marshaled)
}
