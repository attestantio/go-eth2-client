// Copyright © 2026 Attestant Limited.
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
	"encoding/json"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// builderExitRequestJSON is the spec representation of the struct.
type builderExitRequestJSON struct {
	SourceAddress bellatrix.ExecutionAddress `json:"source_address"`
	Pubkey        phase0.BLSPubKey           `json:"pubkey"`
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderExitRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderExitRequestJSON{
		SourceAddress: b.SourceAddress,
		Pubkey:        b.Pubkey,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderExitRequest) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&builderExitRequestJSON{}, input)
	if err != nil {
		return err
	}

	if err := b.SourceAddress.UnmarshalJSON(raw["source_address"]); err != nil {
		return errors.Wrap(err, "source_address")
	}

	if err := b.Pubkey.UnmarshalJSON(raw["pubkey"]); err != nil {
		return errors.Wrap(err, "pubkey")
	}

	return nil
}
