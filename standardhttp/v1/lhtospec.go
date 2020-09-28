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

package v1

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"regexp"
)

// lhEpochIntegerRE is a regular expression to fix the fact that Lighthouse returns epochs as integers.
var lhEpochIntegerRE = regexp.MustCompile(`"epoch":([0-9]+)`)

// lhGenesisEpochIntegerRE is another regular expression to fix the fact that Lighthouse returns epochs as integers.
var lhGenesisEpochIntegerRE = regexp.MustCompile(`"GENESIS_EPOCH":([0-9]+)`)

// lhToSpec patches Lighthouse data to match the spec.
func (s *Service) lhToSpec(ctx context.Context, input io.Reader) (io.Reader, error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	data = lhEpochIntegerRE.ReplaceAll(data, []byte(`"epoch":"$1"`))
	data = lhGenesisEpochIntegerRE.ReplaceAll(data, []byte(`"GENESIS_EPOCH":"$1"`))
	return bytes.NewReader(data), nil
}
