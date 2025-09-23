package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// altair
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.BeaconBlockBody{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.BeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconstate_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.BeaconState{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "contributionandproof_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.ContributionAndProof{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SignedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedcontributionandproof_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SignedContributionAndProof{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "syncaggregate_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SyncAggregate{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "syncaggregatorselectiondata_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SyncAggregatorSelectionData{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "synccommittee_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SyncCommittee{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "synccommitteecontribution_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SyncCommitteeContribution{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "synccommitteemessage_ssz.go"),
		codegen.WithType(reflect.TypeOf(&altair.SyncCommitteeMessage{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
