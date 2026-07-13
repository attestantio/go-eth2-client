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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// MarshalParticipationFlags emits the spec form: a JSON array of decimal
// strings, one per validator.
func MarshalParticipationFlags(flags []ParticipationFlags) (json.RawMessage, error) {
	strs := make([]string, len(flags))
	for i, f := range flags {
		strs[i] = strconv.FormatUint(uint64(f), 10)
	}

	return json.Marshal(strs)
}

// ParseParticipationFlags decodes a beacon-state `*_epoch_participation`
// field. It accepts the spec form (JSON array of decimal strings, e.g.
// `["0","3","2"]`) as well as the form Erigon's Caplin emits, which is a
// single 0x-prefixed hex string of the raw flag bytes. `field` is included
// in error messages.
func ParseParticipationFlags(raw json.RawMessage, field string) ([]ParticipationFlags, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	if trimmed[0] == '[' {
		var strs []string
		if err := json.Unmarshal(trimmed, &strs); err != nil {
			return nil, errors.Wrap(err, field)
		}

		out := make([]ParticipationFlags, len(strs))
		for i, s := range strs {
			if s == "" {
				return nil, errors.Errorf("%s: empty value at index %d", field, i)
			}
			v, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return nil, errors.Wrapf(err, "%s: invalid value %q at index %d", field, s, i)
			}
			out[i] = ParticipationFlags(v)
		}

		return out, nil
	}

	var hexStr string
	if err := json.Unmarshal(trimmed, &hexStr); err != nil {
		return nil, errors.Wrap(err, field)
	}

	decoded, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s: invalid hex", field)
	}

	out := make([]ParticipationFlags, len(decoded))
	for i, b := range decoded {
		out[i] = ParticipationFlags(b)
	}

	return out, nil
}
