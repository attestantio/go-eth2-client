// Copyright © 2020 Attestant Limited.
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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// SignedVoluntaryExit provides information about a signed voluntary exit.
type SignedVoluntaryExit struct {
	Message   *VoluntaryExit
	Signature []byte `ssz-size:"96"`
}

// signedVoluntaryExitJSON is the spec representation of the struct.
type signedVoluntaryExitJSON struct {
	Message   *VoluntaryExit `json:"message"`
	Signature string         `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedVoluntaryExit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedVoluntaryExitJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedVoluntaryExit) UnmarshalJSON(input []byte) error {
	var err error

	var signedVoluntaryExitJSON signedVoluntaryExitJSON
	if err = json.Unmarshal(input, &signedVoluntaryExitJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	s.Message = signedVoluntaryExitJSON.Message
	if s.Message == nil {
		return errors.New("message missing")
	}
	if s.Signature, err = hex.DecodeString(strings.TrimPrefix(signedVoluntaryExitJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(s.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (s *SignedVoluntaryExit) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
