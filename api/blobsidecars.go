// Copyright © 2023 - 2026 Attestant Limited.
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
	dynssz "github.com/pk910/dynamic-ssz"
)

// BlobSidecars is an API construct to allow decoding an array of blob sidecars.
type BlobSidecars struct {
	Sidecars []*deneb.BlobSidecar
}

// blobSidecarsSSZ is the SSZ wrapper for the BlobSidecars object.
type blobSidecarsSSZ = dynssz.TypeWrapper[struct {
	Sidecars []*deneb.BlobSidecar `ssz-max:"72"`
}, []*deneb.BlobSidecar]

// UnmarshalSSZ ssz unmarshals the BlobSidecars object.
func (b *BlobSidecars) UnmarshalSSZ(buf []byte) error {
	return b.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// UnmarshalSSZDyn ssz unmarshals the BlobSidecars object using the supplied dynamic SSZ instance,
// allowing the caller to decode against a custom (non-mainnet) spec rather than the mainnet global.
func (b *BlobSidecars) UnmarshalSSZDyn(dynSSZ *dynssz.DynSsz, buf []byte) error {
	blobs := blobSidecarsSSZ{}
	if err := dynSSZ.UnmarshalSSZ(&blobs, buf); err != nil {
		return err
	}

	b.Sidecars = blobs.Data

	return nil
}

// MarshalSSZ ssz marshals the BlobSidecars object.
func (b *BlobSidecars) MarshalSSZ() ([]byte, error) {
	return dynssz.GetGlobalDynSsz().MarshalSSZ(&blobSidecarsSSZ{
		Data: b.Sidecars,
	})
}

// SizeSSZ returns the size of the BlobSidecars object.
func (b *BlobSidecars) SizeSSZ() int {
	// The error can only be non-nil for a structurally invalid type, which cannot
	// happen for this wrapper, so it is safe to discard here.
	size, _ := dynssz.GetGlobalDynSsz().SizeSSZ(&blobSidecarsSSZ{
		Data: b.Sidecars,
	})

	return size
}

// HashTreeRoot ssz hashes the BlobSidecars object.
func (b *BlobSidecars) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(&blobSidecarsSSZ{
		Data: b.Sidecars,
	})
}
