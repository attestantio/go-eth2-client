// Copyright © 2020 Attestant Limited.
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

package phase0

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// ProposerSlashing provides information about a proposer slashing
type ProposerSlashing struct {
	Header1 *SignedBeaconBlockHeader
	Header2 *SignedBeaconBlockHeader
}

// proposerSlashingJSON is the spec representation of the struct.
type proposerSlashingJSON struct {
	SignedHeader1 *SignedBeaconBlockHeader `json:"signed_header_1"`
	SignedHeader2 *SignedBeaconBlockHeader `json:"signed_header_2"`
}

// MarshalJSON implements json.Marshaler.
func (p *ProposerSlashing) MarshalJSON() ([]byte, error) {
	return json.Marshal(&proposerSlashingJSON{
		SignedHeader1: p.Header1,
		SignedHeader2: p.Header2,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *ProposerSlashing) UnmarshalJSON(input []byte) error {
	var err error

	var proposerSlashingJSON proposerSlashingJSON
	if err = json.Unmarshal(input, &proposerSlashingJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if proposerSlashingJSON.SignedHeader1 == nil {
		return errors.New("signed header 1 missing")
	}
	p.Header1 = proposerSlashingJSON.SignedHeader1
	if proposerSlashingJSON.SignedHeader2 == nil {
		return errors.New("signed header 2 missing")
	}
	p.Header2 = proposerSlashingJSON.SignedHeader2

	return nil
}

// String returns a string version of the structure.
func (p *ProposerSlashing) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
