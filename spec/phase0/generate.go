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

// Need to `go install github.com/ferranbt/fastssz/sszgen@latest` for this to work.
//go:generate rm -f aggregateandproof_ssz.go attestationdata_ssz.go attestation_ssz.go attesterslashing_ssz.go beaconblockbody_ssz.go beaconblock_ssz.go beaconblockheader_ssz.go beaconstate_ssz.go checkpoint_ssz.go depositdata_ssz.go deposit_ssz.go depositmessage_ssz.go eth1data_ssz.go forkdata_ssz.go fork_ssz.go indexedattestation_ssz.go pendingattestation_ssz.go proposerslashing_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go signedbeaconblockheader_ssz.go signedvoluntaryexit_ssz.go signingdata_ssz.go validator_ssz.go voluntaryexit_ssz.go
//go:generate sszgen -suffix ssz -path . --objs AggregateAndProof,AttestationData,Attestation,AttesterSlashing,BeaconBlockBody,BeaconBlock,BeaconBlockHeader,BeaconState,Checkpoint,Deposit,DepositData,DepositMessage,ETH1Data,Fork,ForkData,IndexedAttestation,PendingAttestation,ProposerSlashing,SignedAggregateAndProof,SignedBeaconBlock,SignedBeaconBlockHeader,SignedVoluntaryExit,SigningData,Validator,VoluntaryExit
//go:generate goimports -w aggregateandproof_ssz.go attestationdata_ssz.go attestation_ssz.go attesterslashing_ssz.go beaconblockbody_ssz.go beaconblock_ssz.go beaconblockheader_ssz.go beaconstate_ssz.go checkpoint_ssz.go depositdata_ssz.go deposit_ssz.go depositmessage_ssz.go eth1data_ssz.go forkdata_ssz.go fork_ssz.go indexedattestation_ssz.go pendingattestation_ssz.go proposerslashing_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go signedbeaconblockheader_ssz.go signedvoluntaryexit_ssz.go signingdata_ssz.go validator_ssz.go voluntaryexit_ssz.go
