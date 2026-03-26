// Copyright © 2021 Attestant Limited.
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

package v1

//nolint:revive
// Need to `go install github.com/pk910/dynamic-ssz/dynssz-gen@latest` for this to work.
//go:generate rm -f blobs_ssz.go signedvalidatorregistration_ssz.go validatorregistration_ssz.go
//go:generate dynssz-gen -package . -legacy -without-dynamic-expressions -types Blobs:blobs_ssz.go,SignedValidatorRegistration:signedvalidatorregistration_ssz.go,ValidatorRegistration:validatorregistration_ssz.go
