// Copyright © 2025, 2026 Attestant Limited.
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

package v1

import (
	"github.com/attestantio/go-eth2-client/spec/deneb"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// These SSZ methods are hand-crafted rather than generated. Blobs is a bare named
// slice, so it cannot carry an `ssz-max` tag directly; the List[Blob, 72] bound is
// instead expressed on the blobsSSZ wrapper below (mirroring api.BlobSidecars). This
// restores the length mix-in in HashTreeRoot and the max-72 cap in UnmarshalSSZ that
// generating against the untagged slice would otherwise drop.

// blobsSSZ is the SSZ wrapper for the Blobs object.
type blobsSSZ = dynssz.TypeWrapper[struct {
	Blobs []*deneb.Blob `ssz-max:"72"`
}, []*deneb.Blob]

// MarshalSSZ ssz marshals the Blobs object.
func (b *Blobs) MarshalSSZ() ([]byte, error) {
	return dynssz.GetGlobalDynSsz().MarshalSSZ(&blobsSSZ{Data: *b})
}

// MarshalSSZTo ssz marshals the Blobs object to the supplied buffer.
func (b *Blobs) MarshalSSZTo(buf []byte) ([]byte, error) {
	return dynssz.GetGlobalDynSsz().MarshalSSZTo(&blobsSSZ{Data: *b}, buf)
}

// UnmarshalSSZ ssz unmarshals the Blobs object.
func (b *Blobs) UnmarshalSSZ(buf []byte) error {
	return b.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// UnmarshalSSZDyn ssz unmarshals the Blobs object using the supplied dynamic SSZ
// instance, allowing the caller to decode against a custom (non-mainnet) spec.
func (b *Blobs) UnmarshalSSZDyn(dynSSZ *dynssz.DynSsz, buf []byte) error {
	wrapper := blobsSSZ{}
	if err := dynSSZ.UnmarshalSSZ(&wrapper, buf); err != nil {
		return err
	}

	*b = wrapper.Data

	return nil
}

// SizeSSZ returns the ssz encoded size in bytes for the Blobs object.
func (b *Blobs) SizeSSZ() int {
	// The error can only be non-nil for a structurally invalid type, which cannot
	// happen for this wrapper, so it is safe to discard here.
	size, _ := dynssz.GetGlobalDynSsz().SizeSSZ(&blobsSSZ{Data: *b})

	return size
}

// HashTreeRoot ssz hashes the Blobs object.
func (b *Blobs) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(&blobsSSZ{Data: *b})
}

// HashTreeRootWith ssz hashes the Blobs object with the supplied hasher.
func (b *Blobs) HashTreeRootWith(hh sszutils.HashWalker) error {
	return dynssz.GetGlobalDynSsz().HashTreeRootWith(&blobsSSZ{Data: *b}, hh)
}
