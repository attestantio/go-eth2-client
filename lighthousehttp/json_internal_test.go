// Copyright © 2020 Attestant Limited.
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

package lighthousehttp

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readCloser struct {
	io.Reader
}

func (s *readCloser) Close() error {
	return nil
}

func TestLHToSpec(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "IntValue",
			input:  `{"test":1}`,
			output: `{"test":"1"}`,
		},
		{
			name:   "IntArray1",
			input:  `{"test":[1]}`,
			output: `{"test":["1"]}`,
		},
		{
			name:   "IntArray2",
			input:  `{"test":[1,2]}`,
			output: `{"test":["1","2"]}`,
		},
		{
			name:   "IntArray3",
			input:  `{"test":[1,2,3]}`,
			output: `{"test":["1","2","3"]}`,
		},
		{
			name:   "IntArray4",
			input:  `{"test":[1,2,3,4]}`,
			output: `{"test":["1","2","3","4"]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputReader, err := lhToSpec(context.Background(), &readCloser{bytes.NewReader([]byte(test.input))})
			require.Nil(t, err)

			output, err := ioutil.ReadAll(outputReader)
			require.Nil(t, err)

			assert.Equal(t, test.output, string(output))
		})
	}
}

func TestSpecToLH(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "IntValue",
			input:  `{"test":"1"}`,
			output: `{"test":1}`,
		},
		{
			name:   "IntArray1",
			input:  `{"test":["1"]}`,
			output: `{"test":[1]}`,
		},
		{
			name:   "IntArray2",
			input:  `{"test":["1","2"]}`,
			output: `{"test":[1,2]}`,
		},
		{
			name:   "IntArray3",
			input:  `{"test":["1","2","3"]}`,
			output: `{"test":[1,2,3]}`,
		},
		{
			name:   "IntArray4",
			input:  `{"test":["1","2","3","4"]}`,
			output: `{"test":[1,2,3,4]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputReader, err := specToLH(context.Background(), bytes.NewReader([]byte(test.input)))
			require.Nil(t, err)

			output, err := ioutil.ReadAll(outputReader)
			require.Nil(t, err)

			assert.Equal(t, test.output, string(output))
		})
	}
}
