// Copyright © 2020, 2023 Attestant Limited.
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

package phase0

//nolint:revive
// Need to `go install github.com/pk910/dynamic-ssz/dynssz-gen@latest` for this to work.
//go:generate rm -f aggregateandproof_ssz.go attestationdata_ssz.go attestation_ssz.go attesterslashing_ssz.go beaconblockbody_ssz.go beaconblock_ssz.go beaconblockheader_ssz.go beaconstate_ssz.go checkpoint_ssz.go depositdata_ssz.go deposit_ssz.go depositmessage_ssz.go eth1data_ssz.go forkdata_ssz.go fork_ssz.go indexedattestation_ssz.go pendingattestation_ssz.go proposerslashing_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go signedbeaconblockheader_ssz.go signedvoluntaryexit_ssz.go signingdata_ssz.go validator_ssz.go voluntaryexit_ssz.go
//go:generate dynssz-gen -package . -legacy -without-dynamic-expressions -types AggregateAndProof:aggregateandproof_ssz.go,AttestationData:attestationdata_ssz.go,Attestation:attestation_ssz.go,AttesterSlashing:attesterslashing_ssz.go,BeaconBlockBody:beaconblockbody_ssz.go,BeaconBlock:beaconblock_ssz.go,BeaconBlockHeader:beaconblockheader_ssz.go,BeaconState:beaconstate_ssz.go,Checkpoint:checkpoint_ssz.go,DepositData:depositdata_ssz.go,Deposit:deposit_ssz.go,DepositMessage:depositmessage_ssz.go,ETH1Data:eth1data_ssz.go,ForkData:forkdata_ssz.go,Fork:fork_ssz.go,IndexedAttestation:indexedattestation_ssz.go,PendingAttestation:pendingattestation_ssz.go,ProposerSlashing:proposerslashing_ssz.go,SignedAggregateAndProof:signedaggregateandproof_ssz.go,SignedBeaconBlock:signedbeaconblock_ssz.go,SignedBeaconBlockHeader:signedbeaconblockheader_ssz.go,SignedVoluntaryExit:signedvoluntaryexit_ssz.go,SigningData:signingdata_ssz.go,Validator:validator_ssz.go,VoluntaryExit:voluntaryexit_ssz.go
