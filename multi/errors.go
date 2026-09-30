// Copyright © 2024 Attestant Limited.
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

package multi

import "errors"

// ErrIncorrectType is returned when the multi client obtain a response type it is not expecting.
var ErrIncorrectType = errors.New("incorrect response type")

// ErrCallNotSupported marks an error as saying that a client cannot serve this
// particular call -- it does not implement the provider interface, or the node
// behind it does not serve the endpoint.  That is a static property of the
// client rather than a health signal, so doCall moves on to the next client
// without deactivating this one: lacking an optional endpoint must not cost a
// healthy client its place in the rotation for every other call.
var ErrCallNotSupported = errors.New("client does not support this call")
