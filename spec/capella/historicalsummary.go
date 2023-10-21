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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// HistoricalSummary represents a summary of past consensus state.
type HistoricalSummary struct {
	BlockSummaryRoot phase0.Root `ssz-size:"32"`
	StateSummaryRoot phase0.Root `ssz-size:"32"`
}

// String returns a string version of the structure.
func (h *HistoricalSummary) String() string {
	data, err := yaml.Marshal(h)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
