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

package prysmgrpc

import (
	"context"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// CurrentEpoch is a helper that calculates the current epoch.
func (s *Service) CurrentEpoch(ctx context.Context) (uint64, error) {
	genesisTime, err := s.GenesisTime(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain genesis time for current epoch")
	}
	slotDuration, err := s.SlotDuration(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain slot duration for current epoch")
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain slots per epoch for current epoch")
	}
	var currentEpoch uint64
	if genesisTime.After(time.Now()) {
		currentEpoch = 0
	} else {
		currentEpoch = uint64(time.Since(genesisTime).Seconds()) / (uint64(slotDuration.Seconds()) * slotsPerEpoch)
	}

	return currentEpoch, nil
}

// CurrentSlot is a helper that calculates the current slot.
func (s *Service) CurrentSlot(ctx context.Context) (uint64, error) {
	genesisTime, err := s.GenesisTime(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain genesis time for current slot")
	}
	slotDuration, err := s.SlotDuration(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain slot duration for current slot")
	}
	var currentSlot uint64
	if genesisTime.After(time.Now()) {
		currentSlot = 0
	} else {
		currentSlot = uint64(time.Since(genesisTime).Seconds()) / uint64(slotDuration.Seconds())
	}

	return currentSlot, nil
}

// parseConfigByteArray parses a byte array returned by the prysm configuration call.
func parseConfigByteArray(val string) ([]byte, error) {
	vals := strings.Split(val[1:len(val)-1], " ")
	res := make([]byte, len(vals))
	for i, val := range vals {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert value %s for byte array", val)
		}
		res[i] = byte(intVal)
	}
	return res, nil
}

func (s *Service) indicesToPubKeys(ctx context.Context, indices []spec.ValidatorIndex) ([]spec.BLSPubKey, error) {
	if len(indices) == 0 {
		// Need to fetch all indices; bypass cache.
		return s.allIndicesToPubKeys(ctx)
	}

	pubKeys := make([]spec.BLSPubKey, len(indices))
	s.indexMapMu.Lock()
	defer s.indexMapMu.Unlock()

	// Need a mapping from validator index to position in res.
	unknownIndexMap := make(map[spec.ValidatorIndex]int)

	// Start by filling in all the keys we already know, and making a note of those we don't.
	unknownIndices := make([]spec.ValidatorIndex, 0, len(indices))
	for i, index := range indices {
		if pubKey, exists := s.indexMap[index]; exists {
			pubKeys[i] = pubKey
		} else {
			unknownIndexMap[index] = len(unknownIndices)
			unknownIndices = append(unknownIndices, index)
		}
	}

	if len(unknownIndices) == 0 {
		// We know all of them.
		return pubKeys, nil
	}

	// Fetch the ones we don't know.
	prysmValidators, err := s.ValidatorsByPubKey(ctx, "head", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validators")
	}

	// Cherry-pick the ones we need from the complete set.
	for k, v := range unknownIndexMap {
		pubKeys[v] = prysmValidators[k].Validator.PublicKey
	}

	return pubKeys, nil
}

// allIndicesToPubKeys fetches a list of the public keys of all known validators, ordered by index.
func (s *Service) allIndicesToPubKeys(ctx context.Context) ([]spec.BLSPubKey, error) {
	debug.PrintStack()
	prysmValidators, err := s.ValidatorsByPubKey(ctx, "head", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validators")
	}

	pubKeys := make([]spec.BLSPubKey, 0, len(prysmValidators))
	for _, validator := range prysmValidators {
		pubKeys = append(pubKeys, validator.Validator.PublicKey)
	}

	return pubKeys, nil
}
