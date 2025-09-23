// Copyright © 2025 Attestant Limited.
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

package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0) //nolint:dogsled // ignore unused vars
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// phase0
	generator.BuildFile(
		filepath.Join(parentDir, "aggregateandproof_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.AggregateAndProof{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "attestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.Attestation{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "attestationdata_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.AttestationData{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "attesterslashing_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.AttesterSlashing{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.BeaconBlockBody{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.BeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblockheader_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.BeaconBlockHeader{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconstate_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.BeaconState{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "checkpoint_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.Checkpoint{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "deposit_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.Deposit{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "depositdata_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.DepositData{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "depositmessage_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.DepositMessage{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "eth1data_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.ETH1Data{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "fork_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.Fork{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "forkdata_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.ForkData{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "indexedattestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.IndexedAttestation{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "pendingattestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.PendingAttestation{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "proposerslashing_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.ProposerSlashing{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedaggregateandproof_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.SignedAggregateAndProof{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.SignedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedbeaconblockheader_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.SignedBeaconBlockHeader{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedvoluntaryexit_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.SignedVoluntaryExit{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signingdata_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.SigningData{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "validator_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.Validator{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "voluntaryexit_ssz.go"),
		codegen.WithType(reflect.TypeOf(&phase0.VoluntaryExit{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
