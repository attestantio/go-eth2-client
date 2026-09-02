// Copyright © 2026 Attestant Limited.
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

package gloas

import "github.com/attestantio/go-eth2-client/spec/phase0"

// BuilderRequestAuth authorizes a builder bid request for a proposal slot.
type BuilderRequestAuth struct {
	Data []byte      `ssz-index:"0" ssz-max:"4096"`
	Slot phase0.Slot `ssz-index:"1"`
}

// SignedBuilderRequestAuth is a signed builder bid request authorization.
type SignedBuilderRequestAuth struct {
	Message   *BuilderRequestAuth `ssz-index:"0"`
	Signature phase0.BLSSignature `ssz-index:"1" ssz-size:"96"`
}

// BuilderEntry configures one direct builder bid request.
type BuilderEntry struct {
	URL                 []byte                    `ssz-index:"0" ssz-max:"2048"`
	Auth                *SignedBuilderRequestAuth `ssz-index:"1"`
	BuilderPubkeys      []phase0.BLSPubKey        `ssz-index:"2" ssz-max:"64" ssz-size:"?,48"`
	MaxExecutionPayment phase0.Gwei               `ssz-index:"3"`
	MinBid              phase0.Gwei               `ssz-index:"4"`
	BuilderBoostFactor  uint64                    `ssz-index:"5"`
}

// BuilderConfig configures P2P and direct builder bid selection for one proposal.
type BuilderConfig struct {
	MinBid             phase0.Gwei     `ssz-index:"0"`
	BuilderBoostFactor uint64          `ssz-index:"1"`
	Builders           []*BuilderEntry `ssz-index:"2" ssz-max:"64"`
}
