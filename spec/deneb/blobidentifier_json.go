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

package deneb

import (
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blobIdentifierJSON is the spec representation of the struct.
type blobIdentifierJSON struct {
	BlockRoot phase0.Root `json:"block_root"`
	Index     string      `json:"index"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobIdentifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobIdentifierJSON{
		BlockRoot: b.BlockRoot,
		Index:     fmt.Sprintf("%d", b.Index),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobIdentifier) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&blobIdentifierJSON{}, input)
	if err != nil {
		return err
	}

	if err := b.BlockRoot.UnmarshalJSON(raw["block_root"]); err != nil {
		return errors.Wrap(err, "block_root")
	}

	if err := b.Index.UnmarshalJSON(raw["index"]); err != nil {
		return errors.Wrap(err, "index")
	}

	return nil
}
