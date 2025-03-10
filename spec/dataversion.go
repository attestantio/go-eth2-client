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
)

var dataVersionStrings = [...]string{
	"unknown",
	"phase0",
	"altair",
	"bellatrix",
	"capella",
	"deneb",
	"electra",
}

var stringToDataVersion = map[string]DataVersion{
	dataVersionStrings[0]: DataVersionUnknown,
	dataVersionStrings[1]: DataVersionPhase0,
	dataVersionStrings[2]: DataVersionAltair,
	dataVersionStrings[3]: DataVersionBellatrix,
	dataVersionStrings[4]: DataVersionCapella,
	dataVersionStrings[5]: DataVersionDeneb,
	dataVersionStrings[6]: DataVersionElectra,
}

// MarshalJSON implements json.Marshaler.
func (d *DataVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", dataVersionStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DataVersion) UnmarshalJSON(input []byte) error {
	fork := strings.ReplaceAll(strings.ToLower(string(input)), `"`, ``)
	version, err := DataVersionFromString(fork)
	if err == nil {
		*d = version
	}
	return err
}

// String returns a string representation of the struct.
func (d DataVersion) String() string {
	if uint64(d) >= uint64(len(dataVersionStrings)) {
		return "unknown"
	}
	return dataVersionStrings[d]
}

func DataVersionFromString(fork string) (DataVersion, error) {
	if fork, ok := stringToDataVersion[fork]; ok {
		return fork, nil
	}
	return DataVersionUnknown, fmt.Errorf("unrecognised data version %s", string(input))
}
