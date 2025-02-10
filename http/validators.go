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
	"encoding/json"
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type validatorsBody struct {
	IDs      []string `json:"ids,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
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
	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}
	span.SetAttributes(attribute.Int("validators", len(opts.Indices)+len(opts.PubKeys)))

	if len(opts.Indices) == 0 && len(opts.PubKeys) == 0 {
		// Request is for all validators; fetch from state.
		return s.validatorsFromState(ctx, opts)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", opts.State)
	query := ""

	body := &validatorsBody{
		IDs:      make([]string, 0),
		Statuses: make([]string, 0),
	}
	for i := range opts.Indices {
		body.IDs = append(body.IDs, fmt.Sprintf("%d", opts.Indices[i]))
	}
	for i := range opts.PubKeys {
		body.IDs = append(body.IDs, opts.PubKeys[i].String())
	}
	for i := range opts.ValidatorStates {
		body.Statuses = append(body.Statuses, opts.ValidatorStates[i].String())
	}

	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal request data"), err)
	}

	httpResponse, err := s.post(ctx, endpoint, query, &opts.Common, bytes.NewReader(reqData), ContentTypeJSON, map[string]string{})
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

	validatorStates := make(map[apiv1.ValidatorState]struct{})
	for _, validatorState := range opts.ValidatorStates {
		validatorStates[validatorState] = struct{}{}
	}

	res := make(map[phase0.ValidatorIndex]*apiv1.Validator, len(validators))
	for i, validator := range validators {
		index := phase0.ValidatorIndex(i)

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
