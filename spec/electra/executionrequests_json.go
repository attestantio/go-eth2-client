// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/pkg/errors"
)

// executionRequestsJSON is the spec representation of the struct.
type executionRequestsJSON struct {
	DepositRequests       []*DepositRequest       `json:"deposit_requests"`
	WithdrawalRequests    []*WithdrawalRequest    `json:"withdrawal_requests"`
	ConsolidationRequests []*ConsolidationRequest `json:"consolidation_requests"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionRequests) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionRequestsJSON{
		DepositRequests:       e.DepositRequests,
		WithdrawalRequests:    e.WithdrawalRequests,
		ConsolidationRequests: e.ConsolidationRequests,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionRequests) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&executionRequestsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["deposit_requests"], &e.DepositRequests); err != nil {
		return errors.Wrap(err, "deposit_requests")
	}
	for i := range e.DepositRequests {
		if e.DepositRequests[i] == nil {
			return fmt.Errorf("deposit receipts entry %d missing", i)
		}
	}

	if err := json.Unmarshal(raw["withdrawal_requests"], &e.WithdrawalRequests); err != nil {
		return errors.Wrap(err, "withdrawal_requests")
	}
	for i := range e.WithdrawalRequests {
		if e.WithdrawalRequests[i] == nil {
			return fmt.Errorf("withdraw requests entry %d missing", i)
		}
	}

	if err := json.Unmarshal(raw["consolidation_requests"], &e.ConsolidationRequests); err != nil {
		return errors.Wrap(err, "consolidation_requests")
	}
	for i := range e.ConsolidationRequests {
		if e.ConsolidationRequests[i] == nil {
			return fmt.Errorf("consolidation requests entry %d missing", i)
		}
	}

	return nil
}
