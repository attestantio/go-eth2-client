// Copyright © 2022, 2023 Attestant Limited.
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

package capella

//nolint:revive
// Need to `go install github.com/pk910/dynamic-ssz/dynssz-gen@latest` for this to work.
//go:generate rm -f beaconblockbody_ssz.go beaconblock_ssz.go beaconstate_ssz.go blstoexecutionchange_ssz.go executionpayloadheader_ssz.go executionpayload_ssz.go historicalsummary_ssz.go signedbeaconblock_ssz.go signedblstoexecutionchange_ssz.go withdrawal_ssz.go
//go:generate dynssz-gen -package . -legacy -without-dynamic-expressions -types BeaconBlockBody:beaconblockbody_ssz.go,BeaconBlock:beaconblock_ssz.go,BeaconState:beaconstate_ssz.go,BLSToExecutionChange:blstoexecutionchange_ssz.go,ExecutionPayloadHeader:executionpayloadheader_ssz.go,ExecutionPayload:executionpayload_ssz.go,HistoricalSummary:historicalsummary_ssz.go,SignedBeaconBlock:signedbeaconblock_ssz.go,SignedBLSToExecutionChange:signedblstoexecutionchange_ssz.go,Withdrawal:withdrawal_ssz.go
