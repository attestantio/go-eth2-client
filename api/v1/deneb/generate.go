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

//nolint:revive
// Need to `go install github.com/pk910/dynamic-ssz/dynssz-gen@latest` for this to work.
//go:generate rm -f blindedbeaconblock_ssz.go blindedbeaconblockbody_ssz.go blockcontents_ssz.go signedblindedbeaconblock_ssz.go signedblockcontents_ssz.go
//go:generate dynssz-gen -package . -legacy -without-dynamic-expressions -types BlindedBeaconBlockBody:blindedbeaconblockbody_ssz.go,BlindedBeaconBlock:blindedbeaconblock_ssz.go,BlockContents:blockcontents_ssz.go,SignedBlindedBeaconBlock:signedblindedbeaconblock_ssz.go,SignedBlockContents:signedblockcontents_ssz.go
