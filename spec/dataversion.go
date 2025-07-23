// Copyright © 2021 - 2023 Attestant Limited.
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
type DataVersion uint64

const (
	// DataVersionUnknown is an unknown data version.
	DataVersionUnknown DataVersion = iota
	// DataVersionPhase0 is data applicable for the initial release of the beacon chain.
	DataVersionPhase0
	// DataVersionAltair is data applicable for the Altair release of the beacon chain.
	DataVersionAltair
	// DataVersionBellatrix is data applicable for the Bellatrix release of the beacon chain.
	DataVersionBellatrix
	// DataVersionCapella is data applicable for the Capella release of the beacon chain.
	DataVersionCapella
	// DataVersionDeneb is data applicable for the Deneb release of the beacon chain.
	DataVersionDeneb
	// DataVersionElectra is data applicable for the Electra release of the beacon chain.
	DataVersionElectra
	// DataVersionFulu is data applicable for the Fulu release of the beacon chain.
	DataVersionFulu
	// DataVersionEIP7732 is data applicable for the EIP-7732 release of the beacon chain.
	DataVersionEIP7732
)

var dataVersionStrings = [...]string{
	"unknown",
	"phase0",
	"altair",
	"bellatrix",
	"capella",
	"deneb",
	"electra",
	"fulu",
	"eip7732",
}

var dataVersionMap = map[string]DataVersion{
	`"phase0"`:    DataVersionPhase0,
	`"altair"`:    DataVersionAltair,
	`"bellatrix"`: DataVersionBellatrix,
	`"capella"`:   DataVersionCapella,
	`"deneb"`:     DataVersionDeneb,
	`"electra"`:   DataVersionElectra,
	`"fulu"`:      DataVersionFulu,
	`"eip7732"`:   DataVersionEIP7732,
}

// MarshalJSON implements json.Marshaler.
func (d *DataVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", dataVersionStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DataVersion) UnmarshalJSON(input []byte) error {
	lower := strings.ToLower(string(input))
	version, ok := dataVersionMap[lower]
	if !ok {
		return fmt.Errorf("unrecognised data version %s", string(input))
	}
	*d = version

	return nil
}

// String returns a string representation of the struct.
func (d DataVersion) String() string {
	if int(d) >= len(dataVersionStrings) {
		return "unknown"
	}

	return dataVersionStrings[d]
}

// DataVersionFromString turns a fork string into a DataVersion
// returns an error if the fork is not recognized.
func DataVersionFromString(fork string) (DataVersion, error) {
	var version DataVersion

	return version, version.UnmarshalJSON([]byte(fmt.Sprintf("\"%v\"", fork)))
}
