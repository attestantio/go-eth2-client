// Copyright © 2026 Attestant Limited.
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

package v1

import (
	"fmt"
	"strings"
)

// BuilderStatus defines the status of a Gloas builder.
type BuilderStatus int

const (
	// BuilderStatusUnknown means no information can be found about the builder.
	BuilderStatusUnknown BuilderStatus = iota
	// BuilderStatusPending means the builder has not exited and is not yet active.
	BuilderStatusPending
	// BuilderStatusActive means the builder is active.
	BuilderStatusActive
	// BuilderStatusExited means the builder has exited.
	BuilderStatusExited
)

var builderStatusStrings = [...]string{
	"unknown",
	"pending",
	"active",
	"exited",
}

// MarshalJSON implements json.Marshaler.
func (v BuilderStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", v.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *BuilderStatus) UnmarshalJSON(input []byte) error {
	switch strings.ToLower(string(input)) {
	case `"unknown"`:
		*v = BuilderStatusUnknown
	case `"pending"`:
		*v = BuilderStatusPending
	case `"active"`:
		*v = BuilderStatusActive
	case `"exited"`:
		*v = BuilderStatusExited
	default:
		return fmt.Errorf("unrecognised builder status %s", string(input))
	}

	return nil
}

func (v BuilderStatus) String() string {
	if v < 0 || int(v) >= len(builderStatusStrings) {
		return builderStatusStrings[0]
	}

	return builderStatusStrings[v]
}
