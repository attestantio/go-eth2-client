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
)

// ValidatorLiveness provides the liveness information for validators.
func (s *Service) ValidatorLiveness(ctx context.Context, opts *api.ValidatorLivenessOpts) (*api.Response[[]*apiv1.ValidatorLiveness], error) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if len(opts.Indices) == 0 {
		return nil, errors.Join(errors.New("no validator indices specified"), client.ErrInvalidOptions)
	}

	// Marshal the indices array into JSON
	reqBody, err := json.Marshal(opts.Indices)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validator indices: %w", err)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/liveness/%d", opts.Epoch)

	httpResponse, err := s.post(ctx,
		endpoint,
		"",
		&api.CommonOpts{},
		bytes.NewReader(reqBody),
		ContentTypeJSON,
		nil,
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request validator liveness"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.ValidatorLiveness{})
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode validator liveness response"), err)
	}

	return &api.Response[[]*apiv1.ValidatorLiveness]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
