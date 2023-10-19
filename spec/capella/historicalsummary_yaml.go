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

package capella

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
)

// historicalSummaryYAML is the spec representation of the struct.
type historicalSummaryYAML struct {
	BlockSummaryRoot string `yaml:"block_summary_root"`
	StateSummaryRoot string `yaml:"state_summary_root"`
}

// MarshalYAML implements yaml.Marshaler.
func (h *HistoricalSummary) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&historicalSummaryYAML{
		BlockSummaryRoot: fmt.Sprintf("%#x", h.BlockSummaryRoot),
		StateSummaryRoot: fmt.Sprintf("%#x", h.StateSummaryRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (h *HistoricalSummary) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data historicalSummaryJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return h.unpack(&data)
}
