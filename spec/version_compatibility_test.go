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

package spec_test

import (
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestVersionTypePackagePath(t *testing.T) {
	tests := []struct {
		name  string
		type_ any
	}{
		{
			name:  "BuilderVersion",
			type_: spec.BuilderVersion(0),
		},
		{
			name:  "DataVersion",
			type_: spec.DataVersion(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, "github.com/attestantio/go-eth2-client/spec", reflect.TypeOf(test.type_).PkgPath())
		})
	}
}
