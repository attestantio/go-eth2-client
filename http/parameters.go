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

package http

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type parameters struct {
	logLevel        zerolog.Level
	address         string
	timeout         time.Duration
	indexChunkSize  int
	pubKeyChunkSize int
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

// WithPubKeyChunkSize sets the maximum number of public kyes to send for individual validator requests.
func WithPubKeyChunkSize(pubKeyChunkSize int) Parameter {
	return parameterFunc(func(p *parameters) {
		p.pubKeyChunkSize = pubKeyChunkSize
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel:        zerolog.GlobalLevel(),
		timeout:         2 * time.Second,
		indexChunkSize:  -1,
		pubKeyChunkSize: -1,
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

	return &parameters, nil
}
