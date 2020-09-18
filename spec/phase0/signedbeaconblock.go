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

// SignedBeaconBlock is a signed beacon block.
type SignedBeaconBlock struct {
	Message   *BeaconBlock
	Signature []byte `ssz-size:"96"`
}

// signedBeaconBlockJSON is the spec representation of the struct.
type signedBeaconBlockJSON struct {
	Message   *BeaconBlock `json:"message"`
	Signature string       `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBeaconBlockJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBeaconBlock) UnmarshalJSON(input []byte) error {
	var err error

	var signedBeaconBlockJSON signedBeaconBlockJSON
	if err = json.Unmarshal(input, &signedBeaconBlockJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if signedBeaconBlockJSON.Message == nil {
		return errors.New("message missing")
	}
	s.Message = signedBeaconBlockJSON.Message
	if signedBeaconBlockJSON.Signature == "" {
		return errors.New("signature missing")
	}
	if s.Signature, err = hex.DecodeString(strings.TrimPrefix(signedBeaconBlockJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(s.Signature) != signatureLength {
		return fmt.Errorf("incorrect length %d for signature", len(s.Signature))
	}

	return nil
}

// String returns a string version of the structure.
func (s *SignedBeaconBlock) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
