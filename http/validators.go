// Copyright © 2020 - 2024 Attestant Limited.
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

package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// indexChunkSizes defines the per-beacon-node size of an index chunk.
// Sizes are derived empirically.
var indexChunkSizes = map[string]int{
	"default":    1024,
	"lighthouse": 8192,
	"lodestar":   1024,
	"nimbus":     8192,
	"prysm":      8192,
	"teku":       8192,
}

// pubKeyChunkSizes defines the per-beacon-node size of a public key chunk.
// Sizes are derived empirically.
var pubKeyChunkSizes = map[string]int{
	"default":    128,
	"lighthouse": 512,
	"lodestar":   128,
	"nimbus":     1024,
	"prysm":      1024,
	"teku":       512,
}

// indexChunkSize is the maximum number of validator indices to send in each request.
func (s *Service) indexChunkSize(ctx context.Context) int {
	if s.userIndexChunkSize > 0 {
		return s.userIndexChunkSize
	}

	var nodeClient string
	response, err := s.NodeClient(ctx)
	if err == nil {
		nodeClient = response.Data
	} else {
		// Use default.
		nodeClient = "default"
	}

	if _, exists := indexChunkSizes[nodeClient]; exists {
		return indexChunkSizes[nodeClient]
	}

	return indexChunkSizes["default"]
}

// pubKeyChunkSize is the maximum number of validator public keys to send in each request.
func (s *Service) pubKeyChunkSize(ctx context.Context) int {
	if s.userPubKeyChunkSize > 0 {
		return s.userPubKeyChunkSize
	}

	var nodeClient string
	response, err := s.NodeClient(ctx)
	if err == nil {
		nodeClient = response.Data
	} else {
		// Use default.
		nodeClient = "default"
	}

	if _, exists := pubKeyChunkSizes[nodeClient]; exists {
		return pubKeyChunkSizes[nodeClient]
	}

	return pubKeyChunkSizes["default"]
}

// Validators provides the validators, with their balance and status, for the given options.
func (s *Service) Validators(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "Validators")
	defer span.End()

	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	span.SetAttributes(attribute.Int("validators", len(opts.Indices)+len(opts.PubKeys)))

	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}
	if len(opts.Indices) > 0 && len(opts.PubKeys) > 0 {
		return nil, errors.Join(errors.New("cannot specify both indices and public keys"), client.ErrInvalidOptions)
	}

	if len(opts.Indices) == 0 && len(opts.PubKeys) == 0 {
		// Request is for all validators; fetch from state.
		return s.validatorsFromState(ctx, opts)
	}

	if !s.reducedMemoryUsage {
		if len(opts.Indices) > s.indexChunkSize(ctx)*16 || len(opts.PubKeys) > s.pubKeyChunkSize(ctx)*16 {
			// Request is for multiple pages of validators; fetch from state.
			return s.validatorsFromState(ctx, opts)
		}
	}

	if len(opts.Indices) > s.indexChunkSize(ctx) || len(opts.PubKeys) > s.pubKeyChunkSize(ctx) {
		return s.chunkedValidators(ctx, opts)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", opts.State)
	query := ""
	switch {
	case len(opts.Indices) > 0:
		ids := make([]string, len(opts.Indices))
		for i := range opts.Indices {
			ids[i] = fmt.Sprintf("%d", opts.Indices[i])
		}
		query = "id=" + strings.Join(ids, ",")
	case len(opts.PubKeys) > 0:
		ids := make([]string, len(opts.PubKeys))
		for i := range opts.PubKeys {
			ids[i] = opts.PubKeys[i].String()
		}
		query = "id=" + strings.Join(ids, ",")
	}
	if len(opts.ValidatorStates) > 0 {
		states := make([]string, len(opts.ValidatorStates))
		for i := range opts.ValidatorStates {
			states[i] = opts.ValidatorStates[i].String()
		}
		if query == "" {
			query = "states=" + strings.Join(states, ",")
		} else {
			query = fmt.Sprintf("%s&states=%s", query, strings.Join(states, ","))
		}
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request validators"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.Validator{})
	if err != nil {
		return nil, err
	}

	// Data is returned as an array but we want it as a map.
	mapData := make(map[phase0.ValidatorIndex]*apiv1.Validator)
	for _, validator := range data {
		mapData[validator.Index] = validator
	}

	return &api.Response[map[phase0.ValidatorIndex]*apiv1.Validator]{
		Data:     mapData,
		Metadata: metadata,
	}, nil
}

// validatorsFromState fetches all validators from state.
// This is more efficient than fetching the validators endpoint, as validators uses JSON only,
// whereas state can be provided using SSZ.
func (s *Service) validatorsFromState(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "validatorsFromState")
	defer span.End()

	stateResponse, err := s.BeaconState(ctx, &api.BeaconStateOpts{State: opts.State, Common: opts.Common})
	if err != nil {
		return nil, err
	}
	if stateResponse == nil {
		return nil, errors.New("no beacon state")
	}

	validators, err := stateResponse.Data.Validators()
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain validators from state"), err)
	}

	balances, err := stateResponse.Data.ValidatorBalances()
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain validator balances from state"), err)
	}

	slot, err := stateResponse.Data.Slot()
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain slot from state"), err)
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, err
	}
	epoch := phase0.Epoch(uint64(slot) / slotsPerEpoch)

	farFutureEpoch, err := s.FarFutureEpoch(ctx)
	if err != nil {
		return nil, err
	}

	// Provide map of required pubkeys or indices.
	indices := make(map[phase0.ValidatorIndex]struct{})
	for _, index := range opts.Indices {
		indices[index] = struct{}{}
	}
	pubkeys := make(map[phase0.BLSPubKey]struct{})
	for _, pubkey := range opts.PubKeys {
		pubkeys[pubkey] = struct{}{}
	}
	validatorStates := make(map[apiv1.ValidatorState]struct{})
	for _, validatorState := range opts.ValidatorStates {
		validatorStates[validatorState] = struct{}{}
	}

	res := make(map[phase0.ValidatorIndex]*apiv1.Validator, len(validators))
	for i, validator := range validators {
		if len(pubkeys) > 0 {
			if _, exists := pubkeys[validator.PublicKey]; !exists {
				// We want specific public keys, and this isn't one of them.  Ignore.
				continue
			}
		}

		index := phase0.ValidatorIndex(i)
		if len(indices) > 0 {
			if _, exists := indices[index]; !exists {
				// We want specific indices, and this isn't one of them.  Ignore.
				continue
			}
		}

		state := apiv1.ValidatorToState(validator, &balances[i], epoch, farFutureEpoch)
		if len(validatorStates) > 0 {
			if _, exists := validatorStates[state]; !exists {
				// We want specific states, and this isn't one of them.  Ignore.
				continue
			}
		}

		res[index] = &apiv1.Validator{
			Index:     index,
			Balance:   balances[i],
			Status:    state,
			Validator: validator,
		}
	}

	return &api.Response[map[phase0.ValidatorIndex]*apiv1.Validator]{
		Data:     res,
		Metadata: stateResponse.Metadata,
	}, nil
}

// chunkedValidators obtains the validators a chunk at a time.
func (s *Service) chunkedValidators(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	if len(opts.Indices) > 0 {
		return s.chunkedValidatorsByIndex(ctx, opts)
	}

	return s.chunkedValidatorsByPubkey(ctx, opts)
}

// chunkedValidatorsByIndex obtains validators with index a chunk at a time.
func (s *Service) chunkedValidatorsByIndex(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	data := make(map[phase0.ValidatorIndex]*apiv1.Validator)
	metadata := make(map[string]any)
	indexChunkSize := s.indexChunkSize(ctx)
	for i := 0; i < len(opts.Indices); i += indexChunkSize {
		chunkStart := i
		chunkEnd := i + indexChunkSize
		if len(opts.Indices) < chunkEnd {
			chunkEnd = len(opts.Indices)
		}
		chunk := opts.Indices[chunkStart:chunkEnd]
		chunkRes, err := s.Validators(ctx, &api.ValidatorsOpts{State: opts.State, Indices: chunk})
		if err != nil {
			return nil, errors.Join(errors.New("failed to obtain chunk"), err)
		}
		for k, v := range chunkRes.Data {
			data[k] = v
		}
		for k, v := range chunkRes.Metadata {
			metadata[k] = v
		}
	}

	return &api.Response[map[phase0.ValidatorIndex]*apiv1.Validator]{
		Data:     data,
		Metadata: metadata,
	}, nil
}

// chunkedValidatorsByIndex obtains validators with public key a chunk at a time.
func (s *Service) chunkedValidatorsByPubkey(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	data := make(map[phase0.ValidatorIndex]*apiv1.Validator)
	metadata := make(map[string]any)
	pubkeyChunkSize := s.pubKeyChunkSize(ctx)
	for i := 0; i < len(opts.PubKeys); i += pubkeyChunkSize {
		chunkStart := i
		chunkEnd := i + pubkeyChunkSize
		if len(opts.PubKeys) < chunkEnd {
			chunkEnd = len(opts.PubKeys)
		}
		chunk := opts.PubKeys[chunkStart:chunkEnd]
		chunkRes, err := s.Validators(ctx, &api.ValidatorsOpts{State: opts.State, PubKeys: chunk})
		if err != nil {
			return nil, errors.Join(errors.New("failed to obtain chunk"), err)
		}
		for k, v := range chunkRes.Data {
			data[k] = v
		}
		for k, v := range chunkRes.Metadata {
			metadata[k] = v
		}
	}

	return &api.Response[map[phase0.ValidatorIndex]*apiv1.Validator]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
