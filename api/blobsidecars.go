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

package api

import (
	"github.com/attestantio/go-eth2-client/spec/deneb"
	ssz "github.com/ferranbt/fastssz"
)

// BlobSidecars is an API construct to allow decoding an array of blob sidecars.
type BlobSidecars struct {
	Sidecars []*deneb.BlobSidecar `ssz-max:"6"`
}

// UnmarshalSSZ ssz unmarshals the BlobSidecars object.
// This is a hand-crafted function, as automatic generation does not support immediate arrays.
func (b *BlobSidecars) UnmarshalSSZ(buf []byte) error {
	num, err := ssz.DivideInt2(len(buf), 131928, 6)
	if err != nil {
		return err
	}
	b.Sidecars = make([]*deneb.BlobSidecar, num)
	for ii := 0; ii < num; ii++ {
		if b.Sidecars[ii] == nil {
			b.Sidecars[ii] = new(deneb.BlobSidecar)
		}
		if err = b.Sidecars[ii].UnmarshalSSZ(buf[ii*131928 : (ii+1)*131928]); err != nil {
			return err
		}
	}

	return nil
}
