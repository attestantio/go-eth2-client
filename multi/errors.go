// Copyright © 2024 - 2026 Attestant Limited.
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

import (
	"errors"

	consensusclient "github.com/attestantio/go-eth2-client"
)

// ErrIncorrectType is returned when the multi client obtain a response type it is not expecting.
var ErrIncorrectType = errors.New("incorrect response type")

// ErrCallNotSupported is client.ErrCallNotSupported. Wrappers outside this
// package return the root value, and doCall skips it instead of deactivating
// the client. A skip is not remembered, so the next call probes that client
// again. The error doCall returns names the active clients only.
var ErrCallNotSupported = consensusclient.ErrCallNotSupported
