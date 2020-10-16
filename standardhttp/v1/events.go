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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/r3labs/sse/v2"
)

// Events provides a channel for the specified events.
//func (s *Service) Events(ctx context.Context, topics []string, handler func(*api.Event)) error {
func (s *Service) Events(ctx context.Context, topics []string) error {
	if len(topics) == 0 {
		return errors.New("no topics supplied")
	}

	reference, err := url.Parse(fmt.Sprintf("/eth/v1/events?topics=%s", strings.Join(topics, "&topics=")))
	if err != nil {
		return errors.Wrap(err, "invalid endpoint")
	}
	url := s.base.ResolveReference(reference).String()
	log.Trace().Str("url", url).Msg("GET request to events stream")
	fmt.Printf("URL is %s\n", url)

	client := sse.NewClient(url)
	go func() {
		client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
			fmt.Printf("Message with topic %s: %s\n", string(msg.Event), string(msg.Data))
		})
	}()

	return nil
}
