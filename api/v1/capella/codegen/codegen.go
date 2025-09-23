package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/api/v1/capella"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// capella
	generator.BuildFile(
		filepath.Join(parentDir, "blindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&capella.BlindedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "blindedbeaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&capella.BlindedBeaconBlockBody{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedblindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&capella.SignedBlindedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
