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

package http

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ContentType defines the builder spec version.
type ContentType int

const (
	// ContentTypeUnknown implies an unknown content type.
	ContentTypeUnknown ContentType = iota
	// ContentTypeSSZ implies an SSZ content type.
	ContentTypeSSZ
	// ContentTypeJSON implies a JSON content type.
	ContentTypeJSON
)

var contentTypeStrings = [...]string{
	"Unknown",
	"SSZ",
	"JSON",
}

// MarshalJSON implements json.Marshaler.
func (c *ContentType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", contentTypeStrings[*c])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ContentType) UnmarshalJSON(input []byte) error {
	var err error
	switch strings.ToUpper(string(input)) {
	case `"SSZ"`:
		*c = ContentTypeSSZ
	case `"JSON"`:
		*c = ContentTypeJSON
	default:
		*c = ContentTypeUnknown
		err = fmt.Errorf("unrecognised content type %s", string(input))
	}

	return err
}

// String returns a string representation of the struct.
func (c ContentType) String() string {
	if int(c) >= len(contentTypeStrings) {
		return contentTypeStrings[0]
	}

	return contentTypeStrings[c]
}

// ParseFromMediaType parses a content type string as per
// http://www.iana.org/assignments/media-types/media-types.xhtml
func ParseFromMediaType(input string) (ContentType, error) {
	if input == "" {
		return ContentTypeUnknown, errors.New("no content type supplied")
	}

	contentTypeParts := strings.Split(input, ";")

	switch strings.ToLower(contentTypeParts[0]) {
	case "application/octet-stream":
		return ContentTypeSSZ, nil
	case "application/json":
		return ContentTypeJSON, nil
	default:
		return ContentTypeUnknown, fmt.Errorf("unrecognised content type %s", input)
	}
}

// MediaType returns the IANA name of the media type.
func (c ContentType) MediaType() string {
	switch c {
	case ContentTypeJSON:
		return "application/json"
	case ContentTypeSSZ:
		return "application/octet-stream"
	default:
		return "unknown"
	}
}
