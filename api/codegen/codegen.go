package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// api
	generator.BuildFile(
		filepath.Join(parentDir, "versionedblindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&api.VersionedBlindedBeaconBlock{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "versionedsignedblindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&api.VersionedSignedBlindedBeaconBlock{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "versionedsignedvalidatorregistration_ssz.go"),
		codegen.WithType(reflect.TypeOf(&api.VersionedSignedValidatorRegistration{})),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
