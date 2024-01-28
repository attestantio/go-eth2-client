// Copyright © 2020 - 2023 Attestant Limited.
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
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// indexChunkSizes defines the per-beacon-node size of an index chunk.
// A request should be no more than 8,000 bytes to work with all currently-supported clients.
// An index has variable size, but assuming 7 characters, including the comma separator, is safe.
// We also need to reserve space for the state ID and the endpoint itself, to be safe we go
// with 500 bytes for this which results in us having comfortable space for 1,000 indices.
// That said, some nodes have their own built-in limits so use them where appropriate.
var indexChunkSizes = map[string]int{
	"default":    1000,
	"lighthouse": 1000,
	"nimbus":     1000,
	"prysm":      1000,
	"teku":       1000,
}

// pubKeyChunkSizes defines the per-beacon-node size of a public key chunk.
// A request should be no more than 8,000 bytes to work with all currently-supported clients.
// A public key, including 0x header and comma separator, takes up 99 bytes.
// We also need to reserve space for the state ID and the endpoint itself, to be safe we go
// with 500 bytes for this which results in us having space for 75 public keys.
// That said, some nodes have their own built-in limits so use them where appropriate.
var pubKeyChunkSizes = map[string]int{
	"default":    75,
	"lighthouse": 75,
	"nimbus":     75,
	"prysm":      75,
	"teku":       75,
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

	switch {
	case strings.Contains(nodeClient, "lighthouse"):
		return indexChunkSizes["lighthouse"]
	case strings.Contains(nodeClient, "nimbus"):
		return indexChunkSizes["nimbus"]
	case strings.Contains(nodeClient, "prysm"):
		return indexChunkSizes["prysm"]
	case strings.Contains(nodeClient, "teku"):
		return indexChunkSizes["teku"]
	default:
		return indexChunkSizes["default"]
	}
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

	switch {
	case strings.Contains(nodeClient, "lighthouse"):
		return pubKeyChunkSizes["lighthouse"]
	case strings.Contains(nodeClient, "nimbus"):
		return pubKeyChunkSizes["nimbus"]
	case strings.Contains(nodeClient, "prysm"):
		return pubKeyChunkSizes["prysm"]
	case strings.Contains(nodeClient, "teku"):
		return pubKeyChunkSizes["teku"]
	default:
		return pubKeyChunkSizes["default"]
	}
}

// Validators provides the validators, with their balance and status, for the given options.
func (s *Service) Validators(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]*apiv1.Validator],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.State == "" {
		return nil, errors.New("no state specified")
	}
	if len(opts.Indices) > 0 && len(opts.PubKeys) > 0 {
		return nil, errors.New("cannot specify both indices and public keys")
	}

	if len(opts.Indices) == 0 && len(opts.PubKeys) == 0 {
		// Request is for all validators; fetch from state.
		return s.validatorsFromState(ctx, opts)
	}

	if len(opts.Indices) > indexChunkSizes["default"]*2 || len(opts.PubKeys) > pubKeyChunkSizes["default"]*2 {
		// Request is for multiple pages of validators; fetch from state.
		return s.validatorsFromState(ctx, opts)
	}

	if len(opts.Indices) > s.indexChunkSize(ctx) || len(opts.PubKeys) > s.pubKeyChunkSize(ctx) {
		return s.chunkedValidators(ctx, opts)
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", opts.State)
	switch {
	case len(opts.Indices) > 0:
		ids := make([]string, len(opts.Indices))
		for i := range opts.Indices {
			ids[i] = fmt.Sprintf("%d", opts.Indices[i])
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, ","))
	case len(opts.PubKeys) > 0:
		ids := make([]string, len(opts.PubKeys))
		for i := range opts.PubKeys {
			ids[i] = opts.PubKeys[i].String()
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, ","))
	}

	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request validators")
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
	stateResponse, err := s.BeaconState(ctx, &api.BeaconStateOpts{State: opts.State, Common: opts.Common})
	if err != nil {
		return nil, err
	}
	if stateResponse == nil {
		return nil, errors.New("no beacon state")
	}

	validators, err := stateResponse.Data.Validators()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validators from state")
	}

	balances, err := stateResponse.Data.ValidatorBalances()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validator balances from state")
	}

	slot, err := stateResponse.Data.Slot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slot from state")
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
			return nil, errors.Wrap(err, "failed to obtain chunk")
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
			return nil, errors.Wrap(err, "failed to obtain chunk")
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
