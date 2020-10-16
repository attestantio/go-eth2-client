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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// PendingAttestation is the Ethereum 2 pending attestation structure.
type PendingAttestation struct {
	AggregationBits bitfield.Bitlist `ssz-max:"2048"`
	Data            *AttestationData
	InclusionDelay  uint64
	ProposerIndex   uint64
}

// pendingAttestationJSON is the spec representation of the struct.
type pendingAttestationJSON struct {
	AggregationBits string           `json:"aggregation_bits"`
	Data            *AttestationData `json:"data"`
	InclusionDelay  string           `json:"inclusion_delay"`
	ProposerIndex   string           `json:"proposer_index"`
}

// MarshalJSON implements json.Marshaler.
func (p *PendingAttestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pendingAttestationJSON{
		AggregationBits: fmt.Sprintf("%#x", []byte(p.AggregationBits)),
		Data:            p.Data,
		InclusionDelay:  fmt.Sprintf("%d", p.InclusionDelay),
		ProposerIndex:   fmt.Sprintf("%d", p.ProposerIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PendingAttestation) UnmarshalJSON(input []byte) error {
	var err error

	var pendingAttestationJSON pendingAttestationJSON
	if err = json.Unmarshal(input, &pendingAttestationJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if pendingAttestationJSON.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}
	if p.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(pendingAttestationJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for aggregation bits")
	}
	p.Data = pendingAttestationJSON.Data
	if p.Data == nil {
		return errors.New("data missing")
	}
	if pendingAttestationJSON.InclusionDelay == "" {
		return errors.New("inclusion delay missing")
	}
	if p.InclusionDelay, err = strconv.ParseUint(pendingAttestationJSON.InclusionDelay, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for inclusion delay")
	}
	if pendingAttestationJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	if p.ProposerIndex, err = strconv.ParseUint(pendingAttestationJSON.ProposerIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}

	return nil
}

// String returns a string version of the structure.
func (p *PendingAttestation) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
