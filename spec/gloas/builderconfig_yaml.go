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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// MarshalYAML implements yaml.Marshaler.
func (b *BuilderRequestAuth) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&builderRequestAuthJSON{
		Data: fmt.Sprintf("%#x", b.Data),
		Slot: fmt.Sprintf("%d", b.Slot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BuilderRequestAuth) UnmarshalYAML(input []byte) error {
	var data builderRequestAuthJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}

// MarshalYAML implements yaml.Marshaler.
func (b *SignedBuilderRequestAuth) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedBuilderRequestAuthJSON{
		Message:   b.Message,
		Signature: fmt.Sprintf("%#x", b.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *SignedBuilderRequestAuth) UnmarshalYAML(input []byte) error {
	var data signedBuilderRequestAuthJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}

// yamlQuotedString renders as an explicitly single-quoted YAML scalar.  A
// builder URL carries "://", which goccy emits as a plain scalar; in the flow
// style this package marshals with, the embedded colon then makes the enclosing
// document re-parse the value as a nested mapping and lose the URL.
type yamlQuotedString string

// MarshalYAML implements yaml.BytesMarshaler.
func (q yamlQuotedString) MarshalYAML() ([]byte, error) {
	return []byte("'" + strings.ReplaceAll(string(q), "'", "''") + "'"), nil
}

// builderEntryYAML mirrors builderEntryJSON with the URL forced to a quoted
// scalar.  Unmarshaling still goes through builderEntryJSON, which takes a
// plain string.
type builderEntryYAML struct {
	URL                 yamlQuotedString          `json:"url"`
	Auth                *SignedBuilderRequestAuth `json:"auth"`
	BuilderPubkeys      []string                  `json:"builder_pubkeys"`
	MaxExecutionPayment string                    `json:"max_execution_payment"`
	MinBid              string                    `json:"min_bid"`
	BuilderBoostFactor  string                    `json:"builder_boost_factor"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BuilderEntry) MarshalYAML() ([]byte, error) {
	builderPubkeys := make([]string, len(b.BuilderPubkeys))
	for i := range b.BuilderPubkeys {
		builderPubkeys[i] = fmt.Sprintf("%#x", b.BuilderPubkeys[i])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&builderEntryYAML{
		URL:                 yamlQuotedString(b.URL),
		Auth:                b.Auth,
		BuilderPubkeys:      builderPubkeys,
		MaxExecutionPayment: fmt.Sprintf("%d", b.MaxExecutionPayment),
		MinBid:              fmt.Sprintf("%d", b.MinBid),
		BuilderBoostFactor:  fmt.Sprintf("%d", b.BuilderBoostFactor),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BuilderEntry) UnmarshalYAML(input []byte) error {
	var data builderEntryJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}

// MarshalYAML implements yaml.Marshaler.
func (b *BuilderConfig) MarshalYAML() ([]byte, error) {
	data, err := yaml.MarshalWithOptions(&builderConfigJSON{
		MinBid:             fmt.Sprintf("%d", b.MinBid),
		BuilderBoostFactor: fmt.Sprintf("%d", b.BuilderBoostFactor),
		Builders:           b.buildersForWire(),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(data, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BuilderConfig) UnmarshalYAML(input []byte) error {
	var data builderConfigJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}
