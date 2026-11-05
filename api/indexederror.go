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

package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// IndexedFailure is the failure of one item of a submitted list.
type IndexedFailure struct {
	// Index is the index of the item in the submitted list, as the node sent it.
	// Decoding does not check it against the length of that list.
	Index int `json:"index"`
	// Message describes the failure.
	Message string `json:"message"`
}

// UnmarshalJSON implements json.Unmarshaler.  It accepts the index as a string
// as well as the spec's number, because Teku sends a string.  A message that is
// not a string is left empty rather than failing the decode.
func (f *IndexedFailure) UnmarshalJSON(input []byte) error {
	var data struct {
		Index   json.Number     `json:"index"`
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}
	index, err := strconv.Atoi(data.Index.String())
	if err != nil {
		return fmt.Errorf("invalid index: %w", err)
	}
	var message string
	_ = json.Unmarshal(data.Message, &message)

	f.Index = index
	f.Message = message

	return nil
}

// IndexedError is a node's report that some items of a submitted list failed,
// decoded from the beacon API IndexedErrorMessage.  The HTTP status code is on
// the accompanying *Error.
type IndexedError struct {
	Message  string           `json:"message"`
	Failures []IndexedFailure `json:"failures"`
}

// UnmarshalJSON implements json.Unmarshaler.  It keeps every failure that
// decodes and skips the rest, so one malformed failure does not discard the
// others.  A message that is not a string is left empty.  It fails if the input
// is not JSON, is a JSON value other than an object or null, or has failures
// that is neither an array nor null.
func (e *IndexedError) UnmarshalJSON(input []byte) error {
	var data struct {
		Message  json.RawMessage   `json:"message"`
		Failures []json.RawMessage `json:"failures"`
	}
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}

	var message string
	_ = json.Unmarshal(data.Message, &message)

	failures := make([]IndexedFailure, 0, len(data.Failures))
	for _, raw := range data.Failures {
		var failure IndexedFailure
		if json.Unmarshal(raw, &failure) == nil {
			failures = append(failures, failure)
		}
	}

	e.Message = message
	e.Failures = failures

	return nil
}

func (e *IndexedError) Error() string {
	failures := make([]string, len(e.Failures))
	for i, failure := range e.Failures {
		failures[i] = fmt.Sprintf("index %d: %s", failure.Index, failure.Message)
	}

	return fmt.Sprintf("%s (%s)", e.Message, strings.Join(failures, "; "))
}
