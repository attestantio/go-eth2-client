// Copyright © 2026 Attestant Limited.
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
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func validBuilderConfigForValidation() *gloas.BuilderConfig {
	return &gloas.BuilderConfig{
		Builders: []*gloas.BuilderEntry{{
			URL: []byte("https://builder.example"),
			Auth: &gloas.SignedBuilderRequestAuth{
				Message:   &gloas.BuilderRequestAuth{Data: []byte{0x01}},
				Signature: phase0.BLSSignature{0x02},
			},
		}},
	}
}

func TestValidateBuilderConfig(t *testing.T) {
	tooManyBuilders := make([]*gloas.BuilderEntry, 65)
	tooManyPubkeys := validBuilderConfigForValidation()
	tooManyPubkeys.Builders[0].BuilderPubkeys = make([]phase0.BLSPubKey, 65)

	tests := []struct {
		name   string
		config *gloas.BuilderConfig
		err    string
	}{
		{
			name:   "Valid",
			config: validBuilderConfigForValidation(),
		},
		{
			name: "MissingConfig",
			err:  "no builder config supplied",
		},
		{
			name:   "BuildersOmitted",
			config: &gloas.BuilderConfig{},
			err:    "no builders supplied",
		},
		{
			name:   "TooManyBuilders",
			config: &gloas.BuilderConfig{Builders: tooManyBuilders},
			err:    "too many builders supplied",
		},
		{
			name:   "MissingBuilder",
			config: &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{nil}},
			err:    "builder 0 missing",
		},
		{
			name: "InvalidURL",
			config: func() *gloas.BuilderConfig {
				config := validBuilderConfigForValidation()
				config.Builders[0].URL = []byte{0xff}
				return config
			}(),
			err: "builder 0 has invalid URL",
		},
		{
			name: "MissingAuthorization",
			config: func() *gloas.BuilderConfig {
				config := validBuilderConfigForValidation()
				config.Builders[0].Auth = nil
				return config
			}(),
			err: "builder 0 has invalid authorization",
		},
		{
			name: "EmptyAuthorizationData",
			config: func() *gloas.BuilderConfig {
				config := validBuilderConfigForValidation()
				config.Builders[0].Auth.Message.Data = []byte{}
				return config
			}(),
			err: "builder 0 has invalid authorization",
		},
		{
			name:   "TooManyPubkeys",
			config: tooManyPubkeys,
			err:    "builder 0 has too many public keys",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateBuilderConfig(test.config)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.err)
			}
		})
	}
}
