// Copyright © 2021 Attestant Limited.
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

package spec

import (
	"fmt"
	"strings"
)

// DataVersion defines the spec version of the data in a response.
type DataVersion int

const (
	// DataVersionPhase0 is data applicable for the initial release of the beacon chain.
	DataVersionPhase0 DataVersion = iota
	// DataVersionAltair is data applicable for the Altair release of the beacon chain.
	DataVersionAltair
)

var responseVersionStrings = [...]string{
	"PHASE0",
	"ALTAIR",
}

// MarshalJSON implements json.Marshaler.
func (d *DataVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", responseVersionStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DataVersion) UnmarshalJSON(input []byte) error {
	var err error
	switch strings.ToUpper(string(input)) {
	case `"PHASE0"`:
		*d = DataVersionPhase0
	case `"ALTAIR"`:
		*d = DataVersionAltair
	default:
		err = fmt.Errorf("unrecognised response version %s", string(input))
	}
	return err
}

// String returns a string representation of the
func (d DataVersion) String() string {
	return responseVersionStrings[d]
}
