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
	Deposits       []*DepositRequest       `json:"deposits"`
	Withdrawals    []*WithdrawalRequest    `json:"withdrawals"`
	Consolidations []*ConsolidationRequest `json:"consolidations"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionRequests) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionRequestsJSON{
		Deposits:       e.Deposits,
		Withdrawals:    e.Withdrawals,
		Consolidations: e.Consolidations,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionRequests) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&executionRequestsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["deposits"], &e.Deposits); err != nil {
		return errors.Wrap(err, "deposits")
	}
	for i := range e.Deposits {
		if e.Deposits[i] == nil {
			return fmt.Errorf("deposits entry %d missing", i)
		}
	}

	if err := json.Unmarshal(raw["withdrawals"], &e.Withdrawals); err != nil {
		return errors.Wrap(err, "withdrawals")
	}
	for i := range e.Withdrawals {
		if e.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}

	if err := json.Unmarshal(raw["consolidations"], &e.Consolidations); err != nil {
		return errors.Wrap(err, "consolidations")
	}
	for i := range e.Consolidations {
		if e.Consolidations[i] == nil {
			return fmt.Errorf("consolidation requests entry %d missing", i)
		}
	}

	return nil
}
