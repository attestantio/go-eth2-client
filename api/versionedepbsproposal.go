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

package api

import (
	"errors"
	"fmt"
	"math/big"

	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedEPBSProposal contains a versioned ePBS block proposal, as produced by
// the gloas-onwards block production endpoint.
type VersionedEPBSProposal struct {
	Version spec.DataVersion
	// ExecutionPayloadIncluded selects which arm below carries the proposal.
	ExecutionPayloadIncluded bool
	BuilderIndex             *gloas.BuilderIndex
	ConsensusValue           *big.Int
	ExecutionValue           *big.Int
	// Gloas is the proposal when the execution payload is not included.
	Gloas *gloas.BeaconBlock
	// GloasContents is the proposal when the execution payload is included.
	GloasContents *apiv1gloas.BlockContents
	// BeaconBlockBodyRoot is the hash tree root of the proposal's beacon block
	// body, calculated with the codec that decoded the proposal -- the one
	// aware of whatever preset the connected node actually runs.  The
	// generated HashTreeRoot methods on the gloas types inline mainnet preset
	// sizes (for example a 512-bit sync committee bitvector) as Go literals
	// and never consult that codec, so they compute the wrong root on any
	// other preset.  A beacon block's own root has no preset-dependent field
	// of its own: it merkleizes slot, proposer index, parent root, state root
	// and this body root, exactly like phase0.BeaconBlockHeader.  The body
	// root is therefore the only preset-sensitive quantity in a block, and
	// whatever obtains the proposal is responsible for setting this field
	// with the spec-aware root.  BodyRoot errors when it is unset rather than
	// falling back to the generated, potentially wrong, Body.HashTreeRoot().
	BeaconBlockBodyRoot *phase0.Root
}

// Slot returns the slot of the proposal.
func (v *VersionedEPBSProposal) Slot() (phase0.Slot, error) {
	block, err := v.block()
	if err != nil {
		return 0, err
	}

	return block.Slot, nil
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedEPBSProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	block, err := v.block()
	if err != nil {
		return 0, err
	}

	return block.ProposerIndex, nil
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedEPBSProposal) RandaoReveal() (phase0.BLSSignature, error) {
	block, err := v.block()
	if err != nil {
		return phase0.BLSSignature{}, err
	}

	if block.Body == nil {
		return phase0.BLSSignature{}, errors.New("no gloas beacon block body")
	}

	return block.Body.RANDAOReveal, nil
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedEPBSProposal) ParentRoot() (phase0.Root, error) {
	block, err := v.block()
	if err != nil {
		return phase0.Root{}, err
	}

	return block.ParentRoot, nil
}

// StateRoot returns the state root of the proposal.
func (v *VersionedEPBSProposal) StateRoot() (phase0.Root, error) {
	block, err := v.block()
	if err != nil {
		return phase0.Root{}, err
	}

	return block.StateRoot, nil
}

// BodyRoot returns the hash tree root of the proposal's beacon block body,
// the value a proposer signs over.
//
// It returns the retained BeaconBlockBodyRoot rather than recomputing it: the
// generated Body.HashTreeRoot() inlines mainnet preset sizes, so recomputing
// here would silently sign the wrong root on any other preset.  There is
// deliberately no fallback to that computation.
func (v *VersionedEPBSProposal) BodyRoot() (phase0.Root, error) {
	block, err := v.block()
	if err != nil {
		return phase0.Root{}, err
	}

	if block.Body == nil {
		return phase0.Root{}, errors.New("no gloas beacon block body")
	}

	if v.BeaconBlockBodyRoot == nil {
		return phase0.Root{}, errors.New("no beacon block body root")
	}

	return *v.BeaconBlockBodyRoot, nil
}

// Root returns the hash tree root of the proposal's beacon block.
//
// It composes a phase0.BeaconBlockHeader from the block's own fields plus the
// retained BodyRoot, rather than calling block.HashTreeRoot(): the generated
// block hasher recurses into the body and inherits the same mainnet-baked
// sizes that make Body.HashTreeRoot() wrong on other presets.  The header's
// fields are all preset-independent, so composing it this way can never
// disagree with BodyRoot() -- do not "simplify" this back to
// block.HashTreeRoot(), that reintroduces the bug BeaconBlockBodyRoot exists
// to avoid.
func (v *VersionedEPBSProposal) Root() (phase0.Root, error) {
	block, err := v.block()
	if err != nil {
		return phase0.Root{}, err
	}

	bodyRoot, err := v.BodyRoot()
	if err != nil {
		return phase0.Root{}, err
	}

	root, err := (&phase0.BeaconBlockHeader{
		Slot:          block.Slot,
		ProposerIndex: block.ProposerIndex,
		ParentRoot:    block.ParentRoot,
		StateRoot:     block.StateRoot,
		BodyRoot:      bodyRoot,
	}).HashTreeRoot()
	if err != nil {
		return phase0.Root{}, err
	}

	return phase0.Root(root), nil
}

// ExecutionPayloadEnvelope returns the execution payload envelope that travelled
// with the block.  It is an error to ask for it when the execution payload was
// not included: the caller must fetch the envelope from the producing node.
func (v *VersionedEPBSProposal) ExecutionPayloadEnvelope() (*gloas.ExecutionPayloadEnvelope, error) {
	contents, err := v.contents()
	if err != nil {
		return nil, err
	}

	if contents.ExecutionPayloadEnvelope == nil {
		return nil, errors.New("no gloas execution payload envelope")
	}

	return contents.ExecutionPayloadEnvelope, nil
}

// KZGProofs returns the KZG proofs that travelled with the block.
func (v *VersionedEPBSProposal) KZGProofs() ([]deneb.KZGProof, error) {
	contents, err := v.contents()
	if err != nil {
		return nil, err
	}

	return contents.KZGProofs, nil
}

// Blobs returns the blobs that travelled with the block.
func (v *VersionedEPBSProposal) Blobs() ([]deneb.Blob, error) {
	contents, err := v.contents()
	if err != nil {
		return nil, err
	}

	return contents.Blobs, nil
}

// Value returns the total value of the proposal, or nil when its execution
// value was not supplied by the beacon node.
func (v *VersionedEPBSProposal) Value() *big.Int {
	if v.ExecutionValue == nil {
		return nil
	}

	value := big.NewInt(0)
	if v.ConsensusValue != nil {
		value.Add(value, v.ConsensusValue)
	}

	return value.Add(value, v.ExecutionValue)
}

// IsEmpty returns true if no proposal is populated.
func (v *VersionedEPBSProposal) IsEmpty() bool {
	return v.Gloas == nil && v.GloasContents == nil
}

// String returns a string version of the structure.
func (v *VersionedEPBSProposal) String() string {
	switch v.Version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix,
		spec.DataVersionCapella, spec.DataVersionDeneb, spec.DataVersionElectra,
		spec.DataVersionFulu:
		return ""
	case spec.DataVersionGloas:
		if v.ExecutionPayloadIncluded {
			if v.GloasContents == nil {
				return ""
			}

			return v.GloasContents.String()
		}

		if v.Gloas == nil {
			return ""
		}

		return v.Gloas.String()
	default:
		return "unknown version"
	}
}

// block returns the beacon block the proposal carries, or an error explaining
// why there is none.  Every accessor that reads the block funnels through here,
// so the version, arm-selection and nil checks cannot drift between them.
func (v *VersionedEPBSProposal) block() (*gloas.BeaconBlock, error) {
	switch v.Version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix,
		spec.DataVersionCapella, spec.DataVersionDeneb, spec.DataVersionElectra,
		spec.DataVersionFulu:
		return nil, fmt.Errorf("no epbs proposal in %s", v.Version)
	case spec.DataVersionGloas:
		if v.ExecutionPayloadIncluded {
			contents, err := v.contents()
			if err != nil {
				return nil, err
			}

			if contents.Block == nil {
				return nil, errors.New("no gloas beacon block")
			}

			return contents.Block, nil
		}

		if v.Gloas == nil {
			return nil, errors.New("no gloas beacon block")
		}

		return v.Gloas, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// contents returns the block contents the proposal carries when the execution
// payload is included, or an error explaining why there are none.
func (v *VersionedEPBSProposal) contents() (*apiv1gloas.BlockContents, error) {
	switch v.Version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix,
		spec.DataVersionCapella, spec.DataVersionDeneb, spec.DataVersionElectra,
		spec.DataVersionFulu:
		return nil, fmt.Errorf("no epbs proposal in %s", v.Version)
	case spec.DataVersionGloas:
		if !v.ExecutionPayloadIncluded {
			return nil, errors.New("no block contents; the execution payload was not included")
		}

		if v.GloasContents == nil {
			return nil, errors.New("no gloas block contents")
		}

		return v.GloasContents, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}
