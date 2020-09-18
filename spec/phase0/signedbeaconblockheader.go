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

// SignedBeaconBlockHeader provides information about a signed beacon block header.
type SignedBeaconBlockHeader struct {
	Message   *BeaconBlockHeader
	Signature []byte `ssz-size:"96"`
}

// signedBeaconBlockHeaderJSON is the spec representation of the struct.
type signedBeaconBlockHeaderJSON struct {
	Message   *BeaconBlockHeader `json:"message"`
	Signature string             `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBeaconBlockHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBeaconBlockHeaderJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBeaconBlockHeader) UnmarshalJSON(input []byte) error {
	var err error

	var signedBeaconBlockHeaderJSON signedBeaconBlockHeaderJSON
	if err = json.Unmarshal(input, &signedBeaconBlockHeaderJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	s.Message = signedBeaconBlockHeaderJSON.Message
	if s.Message == nil {
		return errors.New("message missing")
	}
	if s.Signature, err = hex.DecodeString(strings.TrimPrefix(signedBeaconBlockHeaderJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(s.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (s *SignedBeaconBlockHeader) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
