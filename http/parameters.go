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
	"errors"
	"net/http"
	"time"

	"github.com/attestantio/go-eth2-client/metrics"
	"github.com/rs/zerolog"
)

type parameters struct {
	logLevel           zerolog.Level
	monitor            metrics.Service
	address            string
	timeout            time.Duration
	indexChunkSize     int
	pubKeyChunkSize    int
	extraHeaders       map[string]string
	enforceJSON        bool
	allowDelayedStart  bool
	hooks              *Hooks
	reducedMemoryUsage bool
	customSpecSupport  bool
	client             *http.Client
	elConnectionCheck  bool
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

// WithMonitor sets the monitor for the service.
func WithMonitor(monitor metrics.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.monitor = monitor
	})
}

// WithAddress provides the address for the endpoint.
func WithAddress(address string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.address = address
	})
}

// WithTimeout sets the maximum duration for all requests to the endpoint.
func WithTimeout(timeout time.Duration) Parameter {
	return parameterFunc(func(p *parameters) {
		p.timeout = timeout
	})
}

// WithIndexChunkSize sets the maximum number of indices to send for individual validator requests.
func WithIndexChunkSize(indexChunkSize int) Parameter {
	return parameterFunc(func(p *parameters) {
		p.indexChunkSize = indexChunkSize
	})
}

// WithPubKeyChunkSize sets the maximum number of public keys to send for individual validator requests.
func WithPubKeyChunkSize(pubKeyChunkSize int) Parameter {
	return parameterFunc(func(p *parameters) {
		p.pubKeyChunkSize = pubKeyChunkSize
	})
}

// WithExtraHeaders sets additional headers to be sent with each HTTP request.
func WithExtraHeaders(headers map[string]string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.extraHeaders = headers
	})
}

// WithEnforceJSON forces all requests and responses to be in JSON, not sending or requesting SSZ.
func WithEnforceJSON(enforceJSON bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.enforceJSON = enforceJSON
	})
}

// WithAllowDelayedStart allows the service to start even if the client is unavailable.
func WithAllowDelayedStart(allowDelayedStart bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.allowDelayedStart = allowDelayedStart
	})
}

// WithHooks sets the hooks for client activation and sync events.
func WithHooks(hooks *Hooks) Parameter {
	return parameterFunc(func(p *parameters) {
		p.hooks = hooks
	})
}

// WithReducedMemoryUsage reduces memory usage by disabling certain actions that may take significant amount of memory.
// Enabling this may result in longer response times.
func WithReducedMemoryUsage(reducedMemoryUsage bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.reducedMemoryUsage = reducedMemoryUsage
	})
}

// WithCustomSpecSupport switches from the built in static SSZ library to a new dynamic SSZ library, which is able to handle
// non-mainnet presets.
// Dynamic SSZ en-/decoding is much slower than the static one, so this should only be used if required.
func WithCustomSpecSupport(customSpecSupport bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.customSpecSupport = customSpecSupport
	})
}

// WithHTTPClient provides a custom HTTP client for communication with the HTTP server.
// If not supplied then a standard HTTP client is used.
func WithHTTPClient(client *http.Client) Parameter {
	return parameterFunc(func(p *parameters) {
		p.client = client
	})
}

// WithELConnectionCheck enables making sure EL is not offline to consider the client synced.
func WithELConnectionCheck(elConnectionCheck bool) Parameter {
	return parameterFunc(func(p *parameters) {
		p.elConnectionCheck = elConnectionCheck
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel:          zerolog.GlobalLevel(),
		timeout:           2 * time.Second,
		indexChunkSize:    -1,
		pubKeyChunkSize:   -1,
		extraHeaders:      make(map[string]string),
		allowDelayedStart: false,
		hooks:             &Hooks{},
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if parameters.address == "" {
		return nil, errors.New("no address specified")
	}
	if parameters.timeout == 0 {
		return nil, errors.New("no timeout specified")
	}
	if parameters.indexChunkSize == 0 {
		return nil, errors.New("no index chunk size specified")
	}
	if parameters.pubKeyChunkSize == 0 {
		return nil, errors.New("no public key chunk size specified")
	}
	if parameters.hooks == nil {
		return nil, errors.New("no hooks specified")
	}

	return &parameters, nil
}
