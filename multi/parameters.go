// Copyright © 2021 Attestant Limited.
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

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/metrics"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type parameters struct {
	logLevel zerolog.Level
	monitor  metrics.Service
	clients  []eth2client.Service
	timeout  time.Duration
}

// Parameter is the interface for service parameters.
type Parameter interface {
	apply(*parameters)
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

// WithClients sets the clients.
func WithClients(clients []eth2client.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.clients = clients
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel: zerolog.GlobalLevel(),
		timeout:  2 * time.Second,
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if parameters.timeout == 0 {
		return nil, errors.New("no timeout specified")
	}
	if len(parameters.clients) == 0 {
		return nil, errors.New("no Ethereum 2 clients specified")
	}

	return &parameters, nil
}
