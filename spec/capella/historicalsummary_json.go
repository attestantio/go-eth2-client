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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// historicalSummaryJSON is the spec representation of the struct.
type historicalSummaryJSON struct {
	BlockSummaryRoot string `json:"block_summary_root"`
	StateSummaryRoot string `json:"state_summary_root"`
}

// MarshalJSON implements json.Marshaler.
func (h *HistoricalSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(&historicalSummaryJSON{
		BlockSummaryRoot: fmt.Sprintf("%#x", h.BlockSummaryRoot),
		StateSummaryRoot: fmt.Sprintf("%#x", h.StateSummaryRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (h *HistoricalSummary) UnmarshalJSON(input []byte) error {
	var data historicalSummaryJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return h.unpack(&data)
}

func (h *HistoricalSummary) unpack(data *historicalSummaryJSON) error {
	if data.BlockSummaryRoot == "" {
		return errors.New("block summary root missing")
	}
	blockSummaryRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlockSummaryRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block summary root")
	}
	if len(blockSummaryRoot) != phase0.RootLength {
		return errors.New("incorrect length for block summary root")
	}
	copy(h.BlockSummaryRoot[:], blockSummaryRoot)

	if data.StateSummaryRoot == "" {
		return errors.New("state summary root missing")
	}
	stateSummaryRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateSummaryRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state summary root")
	}
	if len(stateSummaryRoot) != phase0.RootLength {
		return errors.New("incorrect length for state summary root")
	}
	copy(h.StateSummaryRoot[:], stateSummaryRoot)

	return nil
}
