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
)

// AggregateAndProof is the Ethereum 2 attestation structure.
type AggregateAndProof struct {
	AggregatorIndex uint64
	Aggregate       *Attestation
	SelectionProof  []byte `ssz-size:"96"`
}

// aggregateAndProofJSON is the spec representation of the struct.
type aggregateAndProofJSON struct {
	AggregatorIndex string       `json:"aggregator_index"`
	Aggregate       *Attestation `json:"aggregate"`
	SelectionProof  string       `json:"selection_proof"`
}

// MarshalJSON implements json.Marshaler.
func (a *AggregateAndProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aggregateAndProofJSON{
		AggregatorIndex: fmt.Sprintf("%d", a.AggregatorIndex),
		Aggregate:       a.Aggregate,
		SelectionProof:  fmt.Sprintf("%#x", a.SelectionProof),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AggregateAndProof) UnmarshalJSON(input []byte) error {
	var err error

	var aggregateAndProofJSON aggregateAndProofJSON
	if err = json.Unmarshal(input, &aggregateAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if aggregateAndProofJSON.AggregatorIndex == "" {
		return errors.New("aggregator index missing")
	}
	if a.AggregatorIndex, err = strconv.ParseUint(aggregateAndProofJSON.AggregatorIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for aggregator index")
	}
	if aggregateAndProofJSON.Aggregate == nil {
		return errors.New("aggregate missing")
	}
	a.Aggregate = aggregateAndProofJSON.Aggregate
	if aggregateAndProofJSON.SelectionProof == "" {
		return errors.New("selection proof missing")
	}
	if a.SelectionProof, err = hex.DecodeString(strings.TrimPrefix(aggregateAndProofJSON.SelectionProof, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for selection proof")
	}
	if len(a.SelectionProof) != signatureLength {
		return errors.New("incorrect length for selection proof")
	}

	return nil
}

// String returns a string version of the structure.
func (a *AggregateAndProof) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
