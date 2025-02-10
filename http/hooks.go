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

package http

import "context"

// HookFunc is a function called when a hook is triggered.
type HookFunc func(ctx context.Context, s *Service)

// Hooks provides hooks that will be called when certain events occur.
type Hooks struct {
	OnActive   HookFunc
	OnInactive HookFunc
	OnSynced   HookFunc
	OnDesynced HookFunc
}
