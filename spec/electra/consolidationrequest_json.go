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

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// consolidationRequestJSON is the spec representation of the struct.
type consolidationRequestJSON struct {
	SourceAddress bellatrix.ExecutionAddress `json:"source_address"`
	SourcePubkey  phase0.BLSPubKey           `json:"source_pubkey"`
	TargetPubkey  phase0.BLSPubKey           `json:"target_pubkey"`
}

// MarshalJSON implements json.Marshaler.
func (e *ConsolidationRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&consolidationRequestJSON{
		SourceAddress: e.SourceAddress,
		SourcePubkey:  e.SourcePubkey,
		TargetPubkey:  e.TargetPubkey,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ConsolidationRequest) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&consolidationRequestJSON{}, input)
	if err != nil {
		return err
	}

	if err := e.SourceAddress.UnmarshalJSON(raw["source_address"]); err != nil {
		return errors.Wrap(err, "source_address")
	}
	if err := e.SourcePubkey.UnmarshalJSON(raw["source_pubkey"]); err != nil {
		return errors.Wrap(err, "source_pubkey")
	}
	if err := e.TargetPubkey.UnmarshalJSON(raw["target_pubkey"]); err != nil {
		return errors.Wrap(err, "target_pubkey")
	}

	return nil
}
