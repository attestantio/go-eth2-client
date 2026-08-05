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

package http

import (
	"context"
	"encoding/json"
	"errors"
)

// marshalRequestBody marshals a request body to the negotiated content type:
// SSZ, using the dynamic codec for the service's spec, unless JSON is enforced.
// It returns the encoded body and the content type to declare for it, or
// ContentTypeUnknown alongside any error.
func (s *Service) marshalRequestBody(ctx context.Context, value any) ([]byte, ContentType, error) {
	if s.enforceJSON {
		body, err := json.Marshal(value)
		if err != nil {
			return nil, ContentTypeUnknown, errors.Join(errors.New("failed to marshal JSON"), err)
		}

		return body, ContentTypeJSON, nil
	}

	ds, err := s.dynSSZForRequest(ctx)
	if err != nil {
		return nil, ContentTypeUnknown, err
	}

	body, err := ds.MarshalSSZ(value)
	if err != nil {
		return nil, ContentTypeUnknown, errors.Join(errors.New("failed to marshal SSZ"), err)
	}

	return body, ContentTypeSSZ, nil
}
