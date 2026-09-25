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

// Package eventtopic describes event topics: the name of each, and the type into which the data
// of its events decodes.  The topics themselves are listed in api/v1, which registers them here
// so that the rest of the module can reach them by type without their being public.
package eventtopic

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/attestantio/go-eth2-client/spec"
)

// Descriptor is a Topic of any type.
type Descriptor interface {
	// Name returns the name of the topic.
	Name() string
	// NewData returns a new value of the type into which the data of the topic's events decodes.
	NewData() any
	// dataType returns that type.
	dataType() reflect.Type
}

var _ Descriptor = (*Topic[struct{}])(nil)

// Topic describes an event topic: its name, and T, the type into which the data of its events
// decodes.
type Topic[T any] struct { //nolint:attgo // Checked above; attgo cannot see a check with a type argument.
	name string
	// version is the fork of the data of the topic's events, for a topic whose events the
	// beacon-API spec wraps as {"version": "...", "data": {...}}, and DataVersionUnknown
	// otherwise.
	version spec.DataVersion
}

// New creates a topic whose events carry their data bare.
func New[T any](name string) Topic[T] {
	return Topic[T]{name: name}
}

// NewVersioned creates a topic whose events the beacon-API spec wraps as
// {"version": "...", "data": {...}}, with data of the given fork.
func NewVersioned[T any](name string, version spec.DataVersion) Topic[T] {
	return Topic[T]{
		name:    name,
		version: version,
	}
}

// Name returns the name of the topic.
func (t Topic[T]) Name() string {
	return t.name
}

// NewData returns a new T.
func (Topic[T]) NewData() any {
	return new(T)
}

func (Topic[T]) dataType() reflect.Type {
	return reflect.TypeFor[T]()
}

// Decode decodes the data of an event of the topic as sent on the events stream.  For a topic
// whose events the beacon-API spec wraps as {"version": "...", "data": {...}}, the version must
// be the topic's fork, as data of another fork could decode into T without error while dropping
// fields T does not have; supporting a later fork means adding it here.  A bare, unwrapped
// object is accepted as well, for nodes that do not wrap it.
func (t Topic[T]) Decode(input []byte) (*T, error) {
	if t.version != spec.DataVersionUnknown {
		var wrapper struct {
			Version string          `json:"version"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(input, &wrapper); err == nil && len(wrapper.Data) > 0 && wrapper.Version != "" {
			version, err := spec.DataVersionFromString(wrapper.Version)
			if err != nil || version != t.version {
				return nil, fmt.Errorf("unsupported version %q for %s event", wrapper.Version, t.name)
			}

			input = wrapper.Data
		}
	}

	data := new(T)
	if err := json.Unmarshal(input, data); err != nil {
		return nil, err
	}

	return data, nil
}

// byType is each registered topic, by the type into which the data of its events decodes.
var byType = map[reflect.Type]Descriptor{}

// Register registers topics, so that they can be looked up with Lookup.  It panics if two
// topics decode into the same type, as Lookup could then not tell them apart.
func Register(topics ...Descriptor) {
	for _, topic := range topics {
		if existing, exists := byType[topic.dataType()]; exists {
			panic(fmt.Sprintf("event topics %s and %s both decode into %v", existing.Name(), topic.Name(), topic.dataType()))
		}

		byType[topic.dataType()] = topic
	}
}

// LookupType returns the registered topic whose data decodes into the given type, if any.
func LookupType(dataType reflect.Type) (Descriptor, bool) {
	topic, exists := byType[dataType]

	return topic, exists
}

// Lookup returns the registered topic whose data decodes into a T.  It panics if there is none.
func Lookup[T any]() Topic[T] {
	descriptor, _ := LookupType(reflect.TypeFor[T]())

	// A topic can be registered by value or, as *Topic also implements Descriptor, by pointer.
	switch topic := descriptor.(type) {
	case Topic[T]:
		return topic
	case *Topic[T]:
		return *topic
	default:
		panic(fmt.Sprintf("no event topic decodes into %v", reflect.TypeFor[T]()))
	}
}
