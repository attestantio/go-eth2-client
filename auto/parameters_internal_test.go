// Copyright © 2025 Attestant Limited.
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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
)

func Test_parseAndCheckParameters(t *testing.T) {
	type args struct {
		params []Parameter
	}
	tests := []struct {
		name    string
		args    args
		want    *parameters
		wantErr bool
	}{
		{
			name: "MissingAddress",
			args: args{
				params: []Parameter{},
			},
			wantErr: true,
		},
		{
			name: "AddressOnly",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.GlobalLevel(),
				timeout:  2 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "AddressWithLogLevel",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					WithLogLevel(zerolog.DebugLevel),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.DebugLevel,
				timeout:  2 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "AddressWithTimeout",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					WithTimeout(5 * time.Minute),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.GlobalLevel(),
				timeout:  5 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "AllParameters",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					WithLogLevel(zerolog.InfoLevel),
					WithTimeout(10 * time.Minute),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.InfoLevel,
				timeout:  10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "OverrideLogLevel",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					WithLogLevel(zerolog.DebugLevel),
					WithLogLevel(zerolog.WarnLevel),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.WarnLevel, // Last one wins
				timeout:  2 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "OverrideTimeout",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					WithTimeout(5 * time.Minute),
					WithTimeout(15 * time.Minute),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.GlobalLevel(),
				timeout:  15 * time.Minute, // Last one wins
			},
			wantErr: false,
		},
		{
			name: "NilParameter",
			args: args{
				params: []Parameter{
					WithAddress("localhost:5051"),
					nil,
					WithTimeout(3 * time.Minute),
				},
			},
			want: &parameters{
				address:  "localhost:5051",
				logLevel: zerolog.GlobalLevel(),
				timeout:  3 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "OnlyNilParameters",
			args: args{
				params: []Parameter{nil, nil},
			},
			wantErr: true, // Should fail because address is missing
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAndCheckParameters(tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseAndCheckParameters() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !cmp.Equal(tt.want, got, cmp.AllowUnexported(parameters{})) {
				t.Errorf("parseAndCheckParameters() = %v, want %v\ndiff=%s", got, tt.want, cmp.Diff(tt.want, got, cmp.AllowUnexported(parameters{})))
			}
		})
	}
}
