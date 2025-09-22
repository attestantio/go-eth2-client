package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// electra
	generator.BuildFile(
		filepath.Join(parentDir, "aggregateandproof_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.AggregateAndProof{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "attestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.Attestation{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "attesterslashing_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.AttesterSlashing{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BeaconBlockBody{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BeaconBlock{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconstate_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BeaconState{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "consolidation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.Consolidation{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "consolidationrequest_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.ConsolidationRequest{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "depositrequest_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.DepositRequest{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "executionrequests_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.ExecutionRequests{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "indexedattestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.IndexedAttestation{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "pendingconsolidation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.PendingConsolidation{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "pendingdeposit_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.PendingDeposit{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "pendingpartialwithdrawal_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.PendingPartialWithdrawal{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.SignedBeaconBlock{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "singleattestation_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.SingleAttestation{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "withdrawalrequest_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.WithdrawalRequest{})),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
