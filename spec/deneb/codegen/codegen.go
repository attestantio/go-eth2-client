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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/pk910/dynamic-ssz/codegen"
)

func main() {
	// Create a code generator instance
	generator := codegen.NewCodeGenerator(nil)

	// Get the parent directory (where types are defined)
	_, currentFile, _, _ := runtime.Caller(0) //nolint:dogsled // ignore unused vars
	currentDir := filepath.Dir(currentFile)
	parentDir := filepath.Dir(currentDir)

	// deneb
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblockbody_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.BeaconBlockBody{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.BeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "beaconstate_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.BeaconState{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "blobidentifier_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.BlobIdentifier{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "blobsidecar_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.BlobSidecar{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "executionpayload_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.ExecutionPayload{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "executionpayloadheader_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.ExecutionPayloadHeader{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)
	generator.BuildFile(
		filepath.Join(parentDir, "signedbeaconblock_ssz.go"),
		codegen.WithType(reflect.TypeOf(&deneb.SignedBeaconBlock{})),
		codegen.WithoutDynamicExpressions(),
		codegen.WithCreateLegacyFn(),
	)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatal("Code generation failed:", err)
	}

	log.Println("Code generation completed successfully!")
}
