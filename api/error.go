// Copyright © 2020, 2023 Attestant Limited.
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

package api

import (
	"fmt"
)

// Error represents an API error.
type Error struct {
	Method     string
	Endpoint   string
	StatusCode int
	Data       []byte
}

func (e Error) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("%s failed with status %d: %s", e.Method, e.StatusCode, string(e.Data))
	}

	return fmt.Sprintf("%s failed with status %d", e.Method, e.StatusCode)
}
