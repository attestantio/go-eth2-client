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

package electra

//nolint:revive
// Need to `go install github.com/pk910/dynamic-ssz/dynssz-gen@latest` for this to work.
//go:generate rm -f aggregateandproof_ssz.go attestation_ssz.go attesterslashing_ssz.go beaconblockbody_ssz.go beaconblock_ssz.go beaconstate_ssz.go consolidation_ssz.go consolidationrequest_ssz.go depositrequest_ssz.go withdrawalrequest_ssz.go executionrequests_ssz.go indexedattestation_ssz.go pendingconsolidation_ssz.go pendingdeposit_ssz.go pendingpartialwithdrawal_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go singleattestation_ssz.go
//go:generate dynssz-gen -package . -legacy -without-dynamic-expressions -types AggregateAndProof:aggregateandproof_ssz.go,Attestation:attestation_ssz.go,AttesterSlashing:attesterslashing_ssz.go,BeaconBlockBody:beaconblockbody_ssz.go,BeaconBlock:beaconblock_ssz.go,BeaconState:beaconstate_ssz.go,Consolidation:consolidation_ssz.go,ConsolidationRequest:consolidationrequest_ssz.go,DepositRequest:depositrequest_ssz.go,WithdrawalRequest:withdrawalrequest_ssz.go,ExecutionRequests:executionrequests_ssz.go,IndexedAttestation:indexedattestation_ssz.go,PendingConsolidation:pendingconsolidation_ssz.go,PendingDeposit:pendingdeposit_ssz.go,PendingPartialWithdrawal:pendingpartialwithdrawal_ssz.go,SignedAggregateAndProof:signedaggregateandproof_ssz.go,SignedBeaconBlock:signedbeaconblock_ssz.go,SingleAttestation:singleattestation_ssz.go
