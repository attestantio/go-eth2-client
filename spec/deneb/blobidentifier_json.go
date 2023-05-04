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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blobIdentifierJSON is the spec representation of the struct.
type blobIdentifierJSON struct {
	BlockRoot string `json:"block_root"`
	Index     string `json:"index"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobIdentifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobIdentifierJSON{
		BlockRoot: b.BlockRoot.String(),
		Index:     fmt.Sprintf("%d", b.Index),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobIdentifier) UnmarshalJSON(input []byte) error {
	var data blobIdentifierJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return b.unpack(&data)
}

func (b *BlobIdentifier) unpack(data *blobIdentifierJSON) error {
	if data.BlockRoot == "" {
		return errors.New("block root missing")
	}
	blockRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block root")
	}
	if len(blockRoot) != phase0.RootLength {
		return errors.New("incorrect length for block root")
	}
	copy(b.BlockRoot[:], blockRoot)

	if data.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(data.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	b.Index = BlobIndex(index)

	return nil
}
