package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// v1
	generator.BuildFile(
		filepath.Join(parentDir, "signedvalidatorregistration_ssz.go"),
		codegen.WithType(reflect.TypeOf(&v1.SignedValidatorRegistration{})),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "validatorregistration_ssz.go"),
		codegen.WithType(reflect.TypeOf(&v1.ValidatorRegistration{})),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
