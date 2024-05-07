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
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// consolidationJSON is the spec representation of the struct.
type consolidationJSON struct {
	SourceIndex phase0.ValidatorIndex `json:"source_index"`
	TargetIndex phase0.ValidatorIndex `json:"target_index"`
	Epoch       phase0.Epoch          `json:"epoch"`
}

// MarshalJSON implements json.Marshaler.
func (c *Consolidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&consolidationJSON{
		SourceIndex: c.SourceIndex,
		TargetIndex: c.TargetIndex,
		Epoch:       c.Epoch,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Consolidation) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&consolidationJSON{}, input)
	if err != nil {
		return err
	}

	if err := c.SourceIndex.UnmarshalJSON(raw["source_index"]); err != nil {
		return errors.Wrap(err, "source_index")
	}
	if err := c.TargetIndex.UnmarshalJSON(raw["target_index"]); err != nil {
		return errors.Wrap(err, "target_index")
	}
	if err := c.Epoch.UnmarshalJSON(raw["epoch"]); err != nil {
		return errors.Wrap(err, "epoch")
	}

	return nil
}
