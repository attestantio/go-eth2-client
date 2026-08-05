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

import (
	"fmt"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// PayloadAttestation represents a payload attestation.
type PayloadAttestation struct {
	// AggregationBits is Bitvector[PTC_SIZE], so its width follows the preset: 64
	// bytes at mainnet, 2 on the minimal preset.  Four of the bitvector's own methods
	// are hard-coded to the mainnet width and must not be used: at any other preset
	// BitAt returns false for every index of a correctly-sized value, SetBitAt
	// silently discards the write, Len reports 512 regardless, and Shift panics
	// outright on a value shorter than eight bytes.  Use AttestedAt, SetAttestedAt and
	// PTCSize instead.  BitIndices, Count and Bytes are length-tolerant up to the
	// mainnet width; above it they silently ignore the excess and so disagree with
	// PTCSize, which a node sending an over-wide value can provoke.
	AggregationBits bitfield.Bitvector512   `dynssz-size:"PTC_SIZE/8" ssz-index:"0" ssz-size:"64"`
	Data            *PayloadAttestationData `ssz-index:"1"`
	Signature       phase0.BLSSignature     `ssz-index:"2"`
}

// PTCSize reports how many payload timeliness committee positions the aggregation
// bits cover.  Use this rather than AggregationBits.Len(); see AggregationBits.
//
// This is the width of the value as decoded, not a width verified against the chain's
// own PTC_SIZE, so for an attestation decoded from a beacon node it is that node's
// claim rather than an established fact.  A caller that needs the distinction must
// check it against the spec itself.
func (p *PayloadAttestation) PTCSize() uint64 {
	return uint64(len(p.AggregationBits)) * 8
}

// checkPTCIndex errors if the given position falls outside the committee the
// aggregation bits cover.  Both accessors share it so the bound and its error can
// never disagree between the read and the write direction.
func (p *PayloadAttestation) checkPTCIndex(ptcIndex uint64) error {
	if size := p.PTCSize(); ptcIndex >= size {
		return fmt.Errorf("payload timeliness committee index %d out of range for committee of %d", ptcIndex, size)
	}

	return nil
}

// AttestedAt reports whether the payload timeliness committee member at the given
// position attested.  The position is an index into the committee returned by
// get_ptc() for the attestation's slot, not a validator index.
//
// Use this rather than AggregationBits.BitAt(), which is mainnet-only; see
// AggregationBits.  An index outside the committee is an error here rather than a
// false, so a failed read cannot be mistaken for an unset bit.
func (p *PayloadAttestation) AttestedAt(ptcIndex uint64) (bool, error) {
	if err := p.checkPTCIndex(ptcIndex); err != nil {
		return false, err
	}

	return p.AggregationBits[ptcIndex/8]&(byte(1)<<(ptcIndex%8)) != 0, nil
}

// SetAttestedAt records participation for the payload timeliness committee member at
// the given position, which is an index into the committee returned by get_ptc() for
// the attestation's slot, not a validator index.
//
// Use this rather than AggregationBits.SetBitAt(), which is mainnet-only; see
// AggregationBits.  bitfield has no length-tolerant setter to delegate to, so the bit
// is set here directly, and an index outside the committee is an error rather than a
// dropped write.
//
// The attested bool mirrors the SetBitAt signature this replaces, so a caller passing a
// computed boolean keeps a single call rather than branching over two method names.
//
//nolint:revive // the bool mirrors the signature this replaces
func (p *PayloadAttestation) SetAttestedAt(ptcIndex uint64, attested bool) error {
	if err := p.checkPTCIndex(ptcIndex); err != nil {
		return err
	}

	bit := byte(1) << (ptcIndex % 8)
	if attested {
		p.AggregationBits[ptcIndex/8] |= bit
	} else {
		p.AggregationBits[ptcIndex/8] &^= bit
	}

	return nil
}

// String returns a string version of the structure.
func (p *PayloadAttestation) String() string {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
