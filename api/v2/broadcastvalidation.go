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

package v2

import (
	"fmt"
	"strings"
)

// BroadcastValidation defines the validation to carry out prior to broadcasting proposals.
type BroadcastValidation int

const (
	// BroadcastValidationGossip means carry out lightweight gossip checks.
	BroadcastValidationGossip BroadcastValidation = iota
	// BroadcastValidationConsensus means carry out full consensus checks.
	BroadcastValidationConsensus
	// BroadcastValidationConsensusAndEquivocation means carry out consensus and equivocation checks.
	BroadcastValidationConsensusAndEquivocation
)

var broadcastValidationStrings = [...]string{
	"gossip",
	"consensus",
	"consensus_and_equivocation",
}

// MarshalJSON implements json.Marshaler.
func (b *BroadcastValidation) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", broadcastValidationStrings[*b]), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BroadcastValidation) UnmarshalJSON(input []byte) error {
	var err error

	switch strings.ToLower(string(input)) {
	case `"gossip"`:
		*b = BroadcastValidationGossip
	case `"consensus"`:
		*b = BroadcastValidationConsensus
	case `"consensus_and_equivocation"`:
		*b = BroadcastValidationConsensusAndEquivocation
	default:
		err = fmt.Errorf("unrecognised broadcast validation %s", string(input))
	}

	return err
}

func (b BroadcastValidation) String() string {
	if b < 0 || int(b) >= len(broadcastValidationStrings) {
		return broadcastValidationStrings[0] // unknown
	}

	return broadcastValidationStrings[b]
}
