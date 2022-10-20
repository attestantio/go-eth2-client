// Copyright © 2020 Attestant Limited.
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
//go:generate rm -f aggregateandproof_encoding.go attestationdata_encoding.go attestation_encoding.go attesterslashing_encoding.go beaconblockbody_encoding.go beaconblock_encoding.go beaconblockheader_encoding.go beaconstate_encoding.go checkpoint_encoding.go depositdata_encoding.go deposit_encoding.go depositmessage_encoding.go eth1data_encoding.go forkdata_encoding.go fork_encoding.go indexedattestation_encoding.go pendingattestation_encoding.go proposerslashing_encoding.go signedaggregateandproof_encoding.go signedbeaconblock_encoding.go signedbeaconblockheader_encoding.go signedvoluntaryexit_encoding.go signingdata_encoding.go validator_encoding.go voluntaryexit_encoding.go
//go:generate sszgen --path . --objs AggregateAndProof,AttestationData,Attestation,AttesterSlashing,BeaconBlockBody,BeaconBlock,BeaconBlockHeader,BeaconState,Checkpoint,Deposit,DepositData,DepositMessage,ETH1Data,Fork,ForkData,IndexedAttestation,PendingAttestation,ProposerSlashing,SignedAggregateAndProof,SignedBeaconBlock,SignedBeaconBlockHeader,SignedVoluntaryExit,SigningData,Validator,VoluntaryExit
