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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// singleAttestationJSON is the spec representation of the struct.
type singleAttestationJSON struct {
	CommitteeIndex string                  `json:"committee_index"`
	AttesterIndex  string                  `json:"attester_index"`
	Data           *phase0.AttestationData `json:"data"`
	Signature      string                  `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (a *SingleAttestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&singleAttestationJSON{
		CommitteeIndex: fmt.Sprintf("%d", a.CommitteeIndex),
		AttesterIndex:  fmt.Sprintf("%d", a.AttesterIndex),
		Data:           a.Data,
		Signature:      fmt.Sprintf("%#x", a.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *SingleAttestation) UnmarshalJSON(input []byte) error {
	var singleAttestationJSON singleAttestationJSON
	err := json.Unmarshal(input, &singleAttestationJSON)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&singleAttestationJSON)
}

func (a *SingleAttestation) unpack(singleAttestationJSON *singleAttestationJSON) error {
	var err error
	if singleAttestationJSON.CommitteeIndex == "" {
		return errors.New("committee index missing")
	}
	committeeIndex, err := strconv.ParseUint(singleAttestationJSON.CommitteeIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for committee index")
	}
	a.CommitteeIndex = phase0.CommitteeIndex(committeeIndex)
	if singleAttestationJSON.AttesterIndex == "" {
		return errors.New("attester index missing")
	}
	attesterIndex, err := strconv.ParseUint(singleAttestationJSON.AttesterIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for attester index")
	}
	a.AttesterIndex = phase0.ValidatorIndex(attesterIndex)
	a.Data = singleAttestationJSON.Data
	if a.Data == nil {
		return errors.New("data missing")
	}
	if singleAttestationJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(singleAttestationJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(a.Signature[:], signature)

	return nil
}
