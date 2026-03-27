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

package gloas

//nolint:revive
// Need to `go install github.com/ferranbt/fastssz/sszgen@latest` for this to work.
//go:generate rm -f beaconblockbody_ssz.go beaconblock_ssz.go beaconstate_ssz.go signedbeaconblock_ssz.go
//go:generate sszgen --suffix=ssz --path . --include ../phase0,../altair,../bellatrix,../capella,../deneb,../electra --objs BeaconBlockBody,BeaconBlock,BeaconState,SignedBeaconBlock
//go:generate goimports -w beaconblockbody_ssz.go beaconblock_ssz.go beaconstate_ssz.go signedbeaconblock_ssz.go
