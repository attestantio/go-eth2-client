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

// Need to `go install github.com/ferranbt/fastssz/sszgen@latest` for this to work.
//go:generate rm -f blindedbeaconblock_ssz.go blindedbeaconblockbody_ssz.go blindedblobsidecar_ssz.go blindedblockcontents_ssz.go blockcontents_ssz.go signedblindedbeaconblock_ssz.go signedblindedblockcontents_ssz.go signedblockcontents_ssz.go
//go:generate sszgen --include ../../../spec/phase0,../../../spec/altair,../../../spec/bellatrix,../../../spec/capella,../../../spec/deneb -path . --suffix ssz -objs BlindedBeaconBlock,BlindedBeaconBlockBody,BlindedBlobSidecar,BlindedBlockContents,BlockContents,SignedBlindedBeaconBlock,SignedBlindedBlockContents,SignedBlockContents
//go:generate goimports -w blindedbeaconblock_ssz.go blindedbeaconblockbody_ssz.go blindedblobsidecar_ssz.go blindedblockcontents_ssz.go blockcontents_ssz.go signedblindedbeaconblock_ssz.go signedblindedblockcontents_ssz.go,signedblockcontents_ssz.go
