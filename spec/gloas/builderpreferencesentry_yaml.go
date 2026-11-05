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
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// builderPreferencesEntryYAML mirrors builderPreferencesEntryJSON with every
// string forced to a single-quoted scalar; see yamlQuotedString.  Quoting here,
// rather than replacing double quotes in the output as sibling types do, keeps
// a double quote inside the URL intact.
type builderPreferencesEntryYAML struct {
	ProposerPubkey      yamlQuotedString          `json:"proposer_pubkey"`
	URL                 yamlQuotedString          `json:"url"`
	Auth                *SignedBuilderRequestAuth `json:"auth"`
	MaxExecutionPayment yamlQuotedString          `json:"max_execution_payment"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BuilderPreferencesEntry) MarshalYAML() ([]byte, error) {
	return yaml.MarshalWithOptions(&builderPreferencesEntryYAML{
		ProposerPubkey:      yamlQuotedString(fmt.Sprintf("%#x", b.ProposerPubkey)),
		URL:                 yamlQuotedString(b.URL),
		Auth:                b.Auth,
		MaxExecutionPayment: yamlQuotedString(fmt.Sprintf("%d", b.MaxExecutionPayment)),
	}, yaml.Flow(true))
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BuilderPreferencesEntry) UnmarshalYAML(input []byte) error {
	var data builderPreferencesEntryJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}
