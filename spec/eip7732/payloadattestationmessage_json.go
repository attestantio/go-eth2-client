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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// payloadAttestationMessageJSON is the spec representation of the struct.
type payloadAttestationMessageJSON struct {
	ValidatorIndex string                  `json:"validator_index"`
	Data           *PayloadAttestationData `json:"data"`
	Signature      string                  `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (p *PayloadAttestationMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&payloadAttestationMessageJSON{
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
		Data:           p.Data,
		Signature:      fmt.Sprintf("%#x", p.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttestationMessage) UnmarshalJSON(input []byte) error {
	var data payloadAttestationMessageJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid validator index")
	}
	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if data.Data == nil {
		return errors.New("data missing")
	}
	p.Data = data.Data

	if data.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid signature")
	}
	copy(p.Signature[:], signature)

	return nil
}
