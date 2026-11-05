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

package gloas

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type builderPreferencesEntryJSON struct {
	ProposerPubkey      string                    `json:"proposer_pubkey"`
	URL                 string                    `json:"url"`
	Auth                *SignedBuilderRequestAuth `json:"auth"`
	MaxExecutionPayment string                    `json:"max_execution_payment"`
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderPreferencesEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderPreferencesEntryJSON{
		ProposerPubkey:      fmt.Sprintf("%#x", b.ProposerPubkey),
		URL:                 string(b.URL),
		Auth:                b.Auth,
		MaxExecutionPayment: fmt.Sprintf("%d", b.MaxExecutionPayment),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderPreferencesEntry) UnmarshalJSON(input []byte) error {
	var data builderPreferencesEntryJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}

	if data.ProposerPubkey == "" {
		return errors.New("proposer public key missing")
	}
	if !strings.HasPrefix(data.ProposerPubkey, "0x") {
		return errors.New("proposer public key missing 0x prefix")
	}
	pubkey, err := hex.DecodeString(strings.TrimPrefix(data.ProposerPubkey, "0x"))
	if err != nil {
		return fmt.Errorf("invalid proposer public key: %w", err)
	}
	if len(pubkey) != phase0.PublicKeyLength {
		return fmt.Errorf("incorrect proposer public key length %d", len(pubkey))
	}

	if data.URL == "" {
		return errors.New("builder URL missing")
	}
	if len(data.URL) > 2048 {
		return errors.New("builder URL exceeds 2048 bytes")
	}
	if data.Auth == nil || data.Auth.Message == nil {
		return errors.New("builder authorization missing")
	}

	maxExecutionPayment, err := strconv.ParseUint(data.MaxExecutionPayment, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid max execution payment: %w", err)
	}

	copy(b.ProposerPubkey[:], pubkey)
	b.URL = []byte(data.URL)
	b.Auth = data.Auth
	b.MaxExecutionPayment = phase0.Gwei(maxExecutionPayment)

	return nil
}
