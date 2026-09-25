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

package eventtopic_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/internal/eventtopic"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	topic := eventtopic.NewVersioned[gloas.PayloadAttestationMessage]("payload_attestation_message", spec.DataVersionGloas)

	message := &gloas.PayloadAttestationMessage{
		ValidatorIndex: 123,
		Data:           &gloas.PayloadAttestationData{Slot: 10},
	}
	bare, err := json.Marshal(message)
	require.NoError(t, err)

	wrapped := func(version string) []byte {
		t.Helper()

		data, err := json.Marshal(map[string]any{"version": version, "data": json.RawMessage(bare)})
		require.NoError(t, err)

		return data
	}

	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{name: "Wrapped", input: wrapped("gloas")},
		{name: "Bare", input: bare},
		{name: "OtherFork", input: wrapped("fulu"), err: `unsupported version "fulu" for payload_attestation_message event`},
		{name: "UnknownFork", input: wrapped("unknown"), err: `unsupported version "unknown" for payload_attestation_message event`},
		// The data type rejects null itself, so no check of Decode's own is needed.
		{name: "WrappedNull", input: []byte(`{"version":"gloas","data":null}`), err: "validator index missing"},
		{name: "BareNull", input: []byte(`null`), err: "validator index missing"},
		{name: "BareNullWithWhitespace", input: []byte(` null `), err: "validator index missing"},
		{name: "Malformed", input: []byte(`invalid`), err: "invalid character 'i' looking for beginning of value"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := topic.Decode(test.input)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, phase0.ValidatorIndex(123), data.ValidatorIndex)
			require.Equal(t, phase0.Slot(10), data.Data.Slot)
		})
	}
}

// registerOnce registers a topic unless one with its type already is.  The registry is global, so
// without this a test registering a topic would panic when run again in the same process, as
// with -count.
func registerOnce(topic eventtopic.Descriptor) {
	if _, exists := eventtopic.LookupType(reflect.TypeOf(topic.NewData()).Elem()); !exists {
		eventtopic.Register(topic)
	}
}

func TestRegisterRejectsSharedType(t *testing.T) {
	type data struct{}
	registerOnce(eventtopic.New[data]("first"))

	require.PanicsWithValue(t, "event topics first and second both decode into eventtopic_test.data", func() {
		eventtopic.Register(eventtopic.New[data]("second"))
	})
	require.Equal(t, "first", eventtopic.Lookup[data]().Name())
}

func TestLookupPanicsWithoutTopic(t *testing.T) {
	type unregistered struct{}

	require.Panics(t, func() { eventtopic.Lookup[unregistered]() })
}

func TestLookupRegisteredByPointer(t *testing.T) {
	type data struct{}
	topic := eventtopic.New[data]("by_pointer")
	registerOnce(&topic)

	require.Equal(t, "by_pointer", eventtopic.Lookup[data]().Name())
}
