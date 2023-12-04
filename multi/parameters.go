// Copyright © 2021, 2023 Attestant Limited.
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

package multi

import (
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/metrics"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type parameters struct {
	logLevel     zerolog.Level
	monitor      metrics.Service
	clients      []consensusclient.Service
	addresses    []string
	timeout      time.Duration
	extraHeaders map[string]string
	enforceJSON  bool
}

// Parameter is the interface for service parameters.
type Parameter interface {
	apply(p *parameters)
}

type parameterFunc func(*parameters)

func (f parameterFunc) apply(p *parameters) {
	f(p)
}

// WithLogLevel sets the log level for the module.
func WithLogLevel(logLevel zerolog.Level) Parameter {
	return parameterFunc(func(p *parameters) {
		p.logLevel = logLevel
	})
}

// WithTimeout sets the timeout for client requests.
func WithTimeout(timeout time.Duration) Parameter {
	return parameterFunc(func(p *parameters) {
		p.timeout = timeout
	})
}

// WithMonitor sets the monitor for the service.
func WithMonitor(monitor metrics.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.monitor = monitor
	})
}

// WithClients sets the pre-existing clients to add to the multi list.
func WithClients(clients []consensusclient.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.clients = clients
	})
}

// WithAddresses sets the addresses of clients to add to the multi list.
func WithAddresses(addresses []string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.addresses = addresses
	})
}

// WithEnforceJSON forces all requests and responses to be in JSON, not sending or requesting SSZ.
func WithEnforceJSON(enforceJSON bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.enforceJSON = enforceJSON
	})
}

// WithExtraHeaders sets additional headers to be sent with each HTTP request.
func WithExtraHeaders(headers map[string]string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.extraHeaders = headers
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel:     zerolog.GlobalLevel(),
		timeout:      2 * time.Second,
		extraHeaders: make(map[string]string),
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if len(parameters.addresses) > 0 && parameters.timeout == 0 {
		return nil, errors.New("no timeout specified")
	}
	if len(parameters.clients)+len(parameters.addresses) == 0 {
		return nil, errors.New("no Ethereum 2 clients specified")
	}

	return &parameters, nil
}
