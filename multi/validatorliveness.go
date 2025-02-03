package multi

import (
	"context"
	"errors"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// ValidatorLiveness provides the liveness information for validators across multiple clients.
func (s *Service) ValidatorLiveness(ctx context.Context,
	opts *api.ValidatorLivenessOpts,
) (
	*api.Response[[]*apiv1.ValidatorLiveness],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options provided for validator liveness")
	}
	if len(opts.Indices) == 0 {
		return nil, errors.New("no validator indices specified")
	}

	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		provider, ok := client.(consensusclient.ValidatorLivenessProvider)
		if !ok {
			return nil, errors.New("client does not support ValidatorLivenessProvider")
		}
		return provider.ValidatorLiveness(ctx, opts)
	}, nil)

	if err != nil {
		return nil, errors.Join(errors.New("failed to retrieve validator liveness"), err)
	}

	response, ok := res.(*api.Response[[]*apiv1.ValidatorLiveness])
	if !ok {
		return nil, ErrIncorrectType
	}

	return response, nil
}
