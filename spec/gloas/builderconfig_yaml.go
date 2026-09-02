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
	"fmt"

	"github.com/goccy/go-yaml"
)

// MarshalYAML implements yaml.Marshaler.
func (b *BuilderConfig) MarshalYAML() ([]byte, error) {
	data, err := yaml.MarshalWithOptions(&builderConfigJSON{
		MinBid:             fmt.Sprintf("%d", b.MinBid),
		BuilderBoostFactor: fmt.Sprintf("%d", b.BuilderBoostFactor),
		Builders:           b.Builders,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(data, []byte(`"`), []byte(`'`)), nil
}
