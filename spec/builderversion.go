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

// BuilderVersion defines the builder spec version.
type BuilderVersion uint64

const (
	// BuilderVersionV1 is applicable for the V1 release of the builder spec.
	BuilderVersionV1 BuilderVersion = iota
)

var responseBuilderVersionStrings = [...]string{
	"v1",
}

// MarshalJSON implements json.Marshaler.
func (d *BuilderVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", responseBuilderVersionStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *BuilderVersion) UnmarshalJSON(input []byte) error {
	var err error
	switch strings.ToLower(string(input)) {
	case `"v1"`:
		*d = BuilderVersionV1
	default:
		err = fmt.Errorf("unrecognised response version %s", string(input))
	}

	return err
}

// String returns a string representation of the struct.
func (d BuilderVersion) String() string {
	if int(d) >= len(responseBuilderVersionStrings) {
		return "unknown"
	}

	return responseBuilderVersionStrings[d]
}
