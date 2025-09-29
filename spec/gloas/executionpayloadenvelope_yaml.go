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

package gloas

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionPayloadEnvelope) MarshalYAML() ([]byte, error) {
	blobCommitments := make([]string, len(e.BlobKZGCommitments))
	for i := range e.BlobKZGCommitments {
		blobCommitments[i] = fmt.Sprintf("%#x", e.BlobKZGCommitments[i])
	}
	yamlBytes, err := yaml.MarshalWithOptions(&executionPayloadEnvelopeJSON{
		Payload:            e.Payload,
		ExecutionRequests:  e.ExecutionRequests,
		BuilderIndex:       fmt.Sprintf("%d", e.BuilderIndex),
		BeaconBlockRoot:    fmt.Sprintf("%#x", e.BeaconBlockRoot),
		Slot:               fmt.Sprintf("%d", e.Slot),
		BlobKZGCommitments: blobCommitments,
		StateRoot:          fmt.Sprintf("%#x", e.StateRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionPayloadEnvelope) UnmarshalYAML(input []byte) error {
	var data executionPayloadEnvelopeJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(marshaled)
}
