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

package v1

import "context"

type providerKey struct{}

// WithProvider returns a context with the provider address set.
// This is the universal mechanism for obtaining the provider address
// across all event types.  Event structs owned by this package also
// carry the address directly in their Provider field.
func WithProvider(ctx context.Context, provider string) context.Context {
	return context.WithValue(ctx, providerKey{}, provider)
}

// Provider extracts the provider address from the context.
// It returns an empty string if no provider has been set.
func Provider(ctx context.Context) string {
	provider, ok := ctx.Value(providerKey{}).(string)
	if !ok {
		return ""
	}

	return provider
}
