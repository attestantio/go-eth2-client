package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/api/v1/electra"
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
		filepath.Join(parentDir, "blindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BlindedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "blindedbeaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BlindedBeaconBlockBody{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "blockcontents_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.BlockContents{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedblindedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.SignedBlindedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedblockcontents_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.SignedBlockContents{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
