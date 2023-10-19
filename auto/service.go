// Copyright © 2020, 2021 Attestant Limited.
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

package auto

import (
	"context"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new Ethereum 2 client service, trying different implementations at the given address.
// Deprecated.  Use the `http` module instead.
func New(ctx context.Context, params ...Parameter) (client.Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "auto").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	// Try HTTP.
	httpClient, err := tryHTTP(ctx, parameters)
	if err == nil {
		return httpClient, nil
	}
	log.Trace().Err(err).Msg("Attempt to connect via HTTP API failed")

	// No luck
	return nil, errors.New("failed to connect to Ethereum 2 client with any known method")
}

func tryHTTP(ctx context.Context, parameters *parameters) (client.Service, error) {
	httpParameters := make([]http.Parameter, 0)
	httpParameters = append(httpParameters, http.WithLogLevel(parameters.logLevel))
	httpParameters = append(httpParameters, http.WithAddress(parameters.address))
	httpParameters = append(httpParameters, http.WithTimeout(parameters.timeout))
	client, err := http.New(ctx, httpParameters...)
	if err != nil {
		return nil, errors.Wrap(err, "failed when trying to open connection with standard API")
	}

	return client, nil
}
