package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/attestantio/go-eth2-client/util/electra"
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
		filepath.Join(parentDir, "consolidation_requests_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.ConsolidationRequests{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "depositrequests_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.DepositRequests{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "withdrawalrequests_ssz.go"),
		codegen.WithType(reflect.TypeOf(&electra.WithdrawalRequests{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
