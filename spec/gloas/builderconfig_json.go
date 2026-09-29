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

type builderRequestAuthJSON struct {
	Data string `json:"data"`
	Slot string `json:"slot"`
}

type signedBuilderRequestAuthJSON struct {
	Message   *BuilderRequestAuth `json:"message"`
	Signature string              `json:"signature"`
}

type builderEntryJSON struct {
	URL                 string                    `json:"url"`
	Auth                *SignedBuilderRequestAuth `json:"auth"`
	BuilderPubkeys      []string                  `json:"builder_pubkeys"`
	MaxExecutionPayment string                    `json:"max_execution_payment"`
	MinBid              string                    `json:"min_bid"`
	BuilderBoostFactor  string                    `json:"builder_boost_factor"`
}

type builderConfigJSON struct {
	MinBid             string          `json:"min_bid"`
	BuilderBoostFactor string          `json:"builder_boost_factor"`
	Builders           []*BuilderEntry `json:"builders"`
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderRequestAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderRequestAuthJSON{
		Data: fmt.Sprintf("%#x", b.Data),
		Slot: fmt.Sprintf("%d", b.Slot),
	})
}

// MarshalJSON implements json.Marshaler.
func (b *SignedBuilderRequestAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBuilderRequestAuthJSON{
		Message:   b.Message,
		Signature: fmt.Sprintf("%#x", b.Signature),
	})
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderEntry) MarshalJSON() ([]byte, error) {
	builderPubkeys := make([]string, len(b.BuilderPubkeys))
	for i := range b.BuilderPubkeys {
		builderPubkeys[i] = fmt.Sprintf("%#x", b.BuilderPubkeys[i])
	}

	return json.Marshal(&builderEntryJSON{
		URL:                 string(b.URL),
		Auth:                b.Auth,
		BuilderPubkeys:      builderPubkeys,
		MaxExecutionPayment: fmt.Sprintf("%d", b.MaxExecutionPayment),
		MinBid:              fmt.Sprintf("%d", b.MinBid),
		BuilderBoostFactor:  fmt.Sprintf("%d", b.BuilderBoostFactor),
	})
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderConfigJSON{
		MinBid:             fmt.Sprintf("%d", b.MinBid),
		BuilderBoostFactor: fmt.Sprintf("%d", b.BuilderBoostFactor),
		Builders:           b.buildersForWire(),
	})
}

// buildersForWire normalises a nil builders slice to an empty one.  builders is
// a required array on the wire and an empty list is meaningful -- it asks for
// no builder bids, leaving only p2p ones -- so a nil slice has to encode as []
// rather than as null, which nothing will accept back.
func (b *BuilderConfig) buildersForWire() []*BuilderEntry {
	if b.Builders == nil {
		return []*BuilderEntry{}
	}

	return b.Builders
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderRequestAuth) UnmarshalJSON(input []byte) error {
	var data builderRequestAuthJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}
	if data.Data == "" {
		return errors.New("authorization data missing")
	}
	if !strings.HasPrefix(data.Data, "0x") {
		return errors.New("authorization data missing 0x prefix")
	}

	authData, err := hex.DecodeString(strings.TrimPrefix(data.Data, "0x"))
	if err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}
	if len(authData) == 0 {
		return errors.New("authorization data empty")
	}
	if len(authData) > 4096 {
		return errors.New("authorization data exceeds 4096 bytes")
	}

	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid slot: %w", err)
	}

	b.Data = authData
	b.Slot = phase0.Slot(slot)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *SignedBuilderRequestAuth) UnmarshalJSON(input []byte) error {
	var data signedBuilderRequestAuthJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}
	if !strings.HasPrefix(data.Signature, "0x") {
		return errors.New("authorization signature missing 0x prefix")
	}

	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}
	if len(signature) != phase0.SignatureLength {
		return fmt.Errorf("incorrect signature length %d", len(signature))
	}

	b.Message = data.Message
	copy(b.Signature[:], signature)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderEntry) UnmarshalJSON(input []byte) error {
	var data builderEntryJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return err
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
	minBid, err := strconv.ParseUint(data.MinBid, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid min bid: %w", err)
	}
	builderBoostFactor, err := strconv.ParseUint(data.BuilderBoostFactor, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid builder boost factor: %w", err)
	}

	if len(data.BuilderPubkeys) > 64 {
		return errors.New("too many builder public keys")
	}

	builderPubkeys := make([]phase0.BLSPubKey, len(data.BuilderPubkeys))
	for i := range data.BuilderPubkeys {
		if !strings.HasPrefix(data.BuilderPubkeys[i], "0x") {
			return fmt.Errorf("builder public key %d missing 0x prefix", i)
		}

		pubkey, err := hex.DecodeString(strings.TrimPrefix(data.BuilderPubkeys[i], "0x"))
		if err != nil {
			return fmt.Errorf("invalid builder public key %d: %w", i, err)
		}
		if len(pubkey) != phase0.PublicKeyLength {
			return fmt.Errorf("incorrect builder public key length %d", len(pubkey))
		}
		copy(builderPubkeys[i][:], pubkey)
	}

	b.URL = []byte(data.URL)
	b.Auth = data.Auth
	b.BuilderPubkeys = builderPubkeys
	b.MaxExecutionPayment = phase0.Gwei(maxExecutionPayment)
	b.MinBid = phase0.Gwei(minBid)
	b.BuilderBoostFactor = builderBoostFactor

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderConfig) UnmarshalJSON(input []byte) error {
	var data builderConfigJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}

	minBid, err := strconv.ParseUint(data.MinBid, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid min bid: %w", err)
	}
	builderBoostFactor, err := strconv.ParseUint(data.BuilderBoostFactor, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid builder boost factor: %w", err)
	}

	if data.Builders == nil {
		return errors.New("builders missing")
	}
	if len(data.Builders) > 64 {
		return errors.New("too many builders")
	}

	b.MinBid = phase0.Gwei(minBid)
	b.BuilderBoostFactor = builderBoostFactor
	b.Builders = data.Builders

	return nil
}
