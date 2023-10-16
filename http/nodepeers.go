package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
	"net/http"
)

// Peers provides the peers a node has.
// peerStates allows you to filter by state; "disconnected", "connecting", "connected", "disconnecting"
// peerDirections allows you to filter by direction; "inbound", "outbound"
func (s *Service) Peers(ctx context.Context, peerStates []string, peerDirections []string) (*v1.Peers, error) {
	var filters []string

	for _, stateFilter := range peerStates {
		filters = append(filters, fmt.Sprintf("state=%s", stateFilter))
	}
	for _, directionFilter := range peerDirections {
		filters = append(filters, fmt.Sprintf("direction=%s", directionFilter))
	}

	return s.peers(ctx, filters)
}

func (s *Service) peersFromJSON(res *httpResponse) (*v1.Peers, error) {
	reader := bytes.NewBuffer(res.body)
	var resp v1.Peers
	if err := json.NewDecoder(reader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to decode peers")
	}
	return &resp, nil
}

func (s *Service) peers(ctx context.Context, filters []string) (*v1.Peers, error) {
	baseRequest := "/eth/v1/node/peers?"
	for ndx, filter := range filters {
		if ndx > 0 {
			baseRequest += "&"
		}
		baseRequest += filter
	}
	res, err := s.get2(ctx, baseRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request node peers")
	}
	if res.statusCode == http.StatusNotFound {
		return nil, nil
	}

	switch res.contentType {
	case ContentTypeJSON:
		return s.peersFromJSON(res)
	default:
		return nil, fmt.Errorf("unhandled content type for peers %v", res.contentType)
	}
}
