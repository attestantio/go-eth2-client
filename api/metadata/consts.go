// Copyright © 2023 Attestant Limited.
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

// Package metadata contains keys for well-known metadata fields provided in
// an API response.
package metadata

const (
	// Finalized describes if the response contains finalized data.
	Finalized = "finalized"
	// ExecutionOptimistic is a boolean value describing if the response contains execution data
	// that has not been fully verified at the time of response.
	ExecutionOptimistic = "execution_optimistic"
	// DependentRoot is the block root on which the returned data is based.
	DependentRoot = "dependent_root"
)
