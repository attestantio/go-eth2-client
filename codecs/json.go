// Copyright © 2023 Attestant Limited.
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

package codecs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// RawJSON generates raw JSON for a struct,
// ensuring that all values are present.
func RawJSON(b any, input []byte) (map[string]json.RawMessage, error) {
	// Make generic map from input.
	base := make(map[string]json.RawMessage)
	if err := json.Unmarshal(input, &base); err != nil {
		return nil, errors.Wrap(err, "invalid JSON")
	}

	// Ensure all values are present.
	elem := reflect.TypeOf(b).Elem()
	fields := elem.NumField()
	for i := 0; i < fields; i++ {
		jsonTags, present := elem.Field(i).Tag.Lookup("json")
		if !present {
			return nil, fmt.Errorf("no json tags for field %d", i)
		}
		tags := strings.Split(jsonTags, ",")
		var emptyAllowed bool
		for i := range tags {
			if tags[i] == "allowempty" {
				// This can be omitted.
				emptyAllowed = true

				break
			}
		}
		if emptyAllowed {
			continue
		}
		if _, exists := base[tags[0]]; !exists {
			// This should be present but is not.
			return nil, fmt.Errorf("%s: missing", tags[0])
		}
	}

	return base, nil
}
