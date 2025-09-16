// Copyright © 2021 - 2024 Attestant Limited.
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

package spec_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/fulu"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func readVersionedBeaconState(t *testing.T, fileName string, version spec.DataVersion) (*spec.VersionedBeaconState, error) {
	// Read the state from the file
	// To download state files for testing, run the following command:
	// curl -X GET "https://<beacon-node-url>/eth/v2/debug/beacon/states/{slot_id}" -H "accept: application/octet-stream" > {netowrkname}_beaconstate_{slot_id}.ssz
	stateBytes, err := readFile(fileName)
	if err != nil {
		return nil, err
	}

	switch version {
	case spec.DataVersionDeneb:
		denebState := &deneb.BeaconState{}
		require.NoError(t, denebState.UnmarshalSSZ(stateBytes))

		state := &spec.VersionedBeaconState{
			Version: spec.DataVersionDeneb,
			Deneb:   denebState,
		}
		return state, nil
	case spec.DataVersionElectra:
		electraState := &electra.BeaconState{}
		require.NoError(t, electraState.UnmarshalSSZ(stateBytes))

		state := &spec.VersionedBeaconState{
			Version: spec.DataVersionElectra,
			Electra: electraState,
		}
		return state, nil
	case spec.DataVersionFulu:
		fuluState := &fulu.BeaconState{}
		require.NoError(t, fuluState.UnmarshalSSZ(stateBytes))

		state := &spec.VersionedBeaconState{
			Version: spec.DataVersionFulu,
			Fulu:    fuluState,
		}
		return state, nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}

func TestStateTree(t *testing.T) {
	state, err := readVersionedBeaconState(t, "holesky_beaconstate_2595934.ssz", spec.DataVersionDeneb)
	if err != nil {
		t.Skip("holesky_beaconstate_2595934.ssz not available")
	}

	tree, err := state.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	root := tree.Hash()
	require.Equal(t, "0x738b800105bb612abc36ea4040312ecf17bc5bd25404529701d0886074820572", fmt.Sprintf("%#x", root))

	hashTreeRoot, err := state.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, root[:], hashTreeRoot[:])
}

func TestFieldIndex(t *testing.T) {
	// Deneb state at Hoodi
	state, err := readVersionedBeaconState(t, "hoodi_beaconstate_100.ssz", spec.DataVersionDeneb)
	if err != nil {
		t.Skip("hoodi_beaconstate_100.ssz not available")
	}

	// Get the field index for the "Validators" field
	fieldIndex, err := state.FieldIndex("Validators")
	require.NoError(t, err)
	require.Equal(t, 11, fieldIndex)

	// Get the field index of non-existent field
	fieldIndex, err = state.FieldIndex("DepositRequestsStartIndex")
	require.ErrorContains(t, err, "field not found")
	require.Equal(t, 0, fieldIndex)

	// Electra state
	state, err = readVersionedBeaconState(t, "hoodi_beaconstate_66000.ssz", spec.DataVersionElectra)
	if err != nil {
		t.Skip("hoodi_beaconstate_66000.ssz not available")
	}

	// Get the field index for the "Validators" field
	fieldIndex, err = state.FieldIndex("Validators")
	require.NoError(t, err)
	require.Equal(t, 11, fieldIndex)

	// Get the field index of newly added field
	fieldIndex, err = state.FieldIndex("DepositRequestsStartIndex")
	require.NoError(t, err)
	require.Equal(t, 28, fieldIndex)
}

func TestFieldRoot(t *testing.T) {
	// Deneb state at Hoodi
	state, err := readVersionedBeaconState(t, "hoodi_beaconstate_100.ssz", spec.DataVersionDeneb)
	if err != nil {
		t.Skip("hoodi_beaconstate_100.ssz not available")
	}
	// Get the root for the "Validators" field
	fieldRoot, err := state.FieldRoot("Validators")
	require.NoError(t, err)

	stateTree, err := state.Deneb.GetTree()
	require.NoError(t, err)

	validatorGeneralizedIndex, err := state.FieldGeneralizedIndex("Validators")
	require.NoError(t, err)
	validatorTree, err := stateTree.Get(validatorGeneralizedIndex)
	require.NoError(t, err)

	validatorRoot := validatorTree.Hash()
	require.NoError(t, err)
	require.Equal(t, fieldRoot[:], validatorRoot[:])
}

func TestProveField(t *testing.T) {
	state, err := readVersionedBeaconState(t, "holesky_beaconstate_2649079.ssz", spec.DataVersionDeneb)
	if err != nil {
		t.Skip("holesky_beaconstate_2649079.ssz not available")
	}

	proof, err := state.ProveField("Balances")
	require.NoError(t, err)
	require.NotNil(t, proof)
	byteSlice := make([]byte, 0)
	for _, d := range proof {
		byteSlice = append(byteSlice, d[:]...)
	}
	require.Equal(t, "0x4210b8b20d920277a62b327759ef42ed5bcbe52513dd4228553995f37b15f893b169594a6c36fbc7ba4b74aa63f87d3442176cfca7d6befe8daf97d6510bd96b2ec272ee5bacd2f3a3fddc49aab127e31990e904ae5f3db808d43d383fa1d711fe9389daeb4a435cba0aaed1dfa536d1bbd711e4f50ae292106650251af358a4c02c8a559c6f3479b12efe98cb1ef318c711e5e6b82b721f2a9f4a825b5eab39", fmt.Sprintf("%#x", byteSlice))

	valid, err := state.VerifyFieldProof(proof, "Balances")
	require.NoError(t, err)
	require.True(t, valid)
}

func TestInvalidValidatorIndex(t *testing.T) {
	state, err := readVersionedBeaconState(t, "holesky_beaconstate_2649079.ssz", spec.DataVersionDeneb)
	if err != nil {
		t.Skip("holesky_beaconstate_2649079.ssz not available")
	}
	validatorIndex := phase0.ValidatorIndex(176565800)
	validator, err := state.ValidatorAtIndex(validatorIndex)
	require.Error(t, err, "validator index out of bounds")
	require.Nil(t, validator)

	balance, err := state.ValidatorBalance(validatorIndex)
	require.Error(t, err, "validator index out of bounds")
	require.Equal(t, phase0.Gwei(0), balance)
}
