// Copyright © 2020 Attestant Limited.
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

package altair

// Slot is a slot number.
type Slot uint64

// Epoch is an epoch number.
type Epoch uint64

// CommitteeIndex is a committee index at a slot.
type CommitteeIndex uint64

// ValidatorIndex is a validator registry index.
type ValidatorIndex uint64

// Gwei is an amount in Gwei.
type Gwei uint64

// Root is a merkle root.
type Root [32]byte

// Version is a fork version.
type Version [4]byte

// DomainType is a domain type.
type DomainType [4]byte

// ForkDigest is a digest of fork data.
type ForkDigest [4]byte

// Domain is a signature domain.
type Domain [32]byte

// BLSPubKey is a BLS12-381 public key.
type BLSPubKey [48]byte

// BLSSignature is a BLS12-381 signature.
type BLSSignature [96]byte

// ParticipationFlags are validator participation flags in an epoch.
type ParticipationFlags uint8

// ParticipationFlag is an individual particiation flag for a validator.
type ParticipationFlag int

const (
	// TimelySourceFlagIndex is set when an attestation has a timely source value.
	TimelySourceFlagIndex ParticipationFlag = iota
	// TimelyTargetFlagIndex is set when an attestation has a timely target value.
	TimelyTargetFlagIndex
	// TimelyHeadFlagIndex is set when an attestation has a timely head value.
	TimelyHeadFlagIndex
)
