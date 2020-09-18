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
	"regexp"
)

// Regexes to go from lighthouse JSON to spec JSON
// Values.
var lhToSpecRe1 = regexp.MustCompile(`:([0-9]+)`)

// Start of array.
var lhToSpecRe2 = regexp.MustCompile(`\[([0-9]+)`)

// Middle of array.
var lhToSpecRe3 = regexp.MustCompile(`,([0-9]+)`)

// End of array.
var lhToSpecRe4 = regexp.MustCompile(`([0-9]+)\]`)

// Regexes to go from spec JSON to lighthouse JSON
// Values.
var specToLHRe1 = regexp.MustCompile(`:"([0-9]+)"`)

// Start of array.
var specToLHRe2 = regexp.MustCompile(`\["([0-9]+)"`)

// Middle of array.
var specToLHRe3 = regexp.MustCompile(`,"([0-9]+)"`)

// End of array.
var specToLHRe4 = regexp.MustCompile(`"([0-9]+)"\]`)

// lhToSpec converts lighthouse JSON to spec JSON.
// This consumes (and closes) the input.
func lhToSpec(ctx context.Context, input io.ReadCloser) (io.Reader, error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if err := input.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close HTTP body")
	}

	// Lighthouse sends numbers unquoted.
	data = lhToSpecRe1.ReplaceAll(data, []byte(`:"$1"`))
	data = lhToSpecRe2.ReplaceAll(data, []byte(`["$1"`))
	data = lhToSpecRe3.ReplaceAll(data, []byte(`,"$1"`))
	data = lhToSpecRe4.ReplaceAll(data, []byte(`"$1"]`))
	return bytes.NewReader(data), nil
}

// specToLH converts spec JSON to lighthouse JSON.
func specToLH(ctx context.Context, input io.Reader) (io.Reader, error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Lighthouse wants numbers unquoted.
	data = specToLHRe1.ReplaceAll(data, []byte(`:$1`))
	data = specToLHRe2.ReplaceAll(data, []byte(`[$1`))
	data = specToLHRe3.ReplaceAll(data, []byte(`,$1`))
	data = specToLHRe4.ReplaceAll(data, []byte(`$1]`))
	return bytes.NewReader(data), nil
}
