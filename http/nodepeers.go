package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"strings"
)

type peersResponseJSON struct {
	Data []apiv1.Peer           `json:"data"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

func (s *Service) NodePeers(ctx context.Context, opts *api.PeerOpts) (*api.Response[*apiv1.Peers], error) {
	var peerResonse peersResponseJSON
	// all options are considered optional
	request := "/eth/v1/node/peers"
	var additionalFields []string

	for _, stateFilter := range opts.State {
		additionalFields = append(additionalFields, fmt.Sprintf("state=%s", stateFilter))
	}

	for _, directionFilter := range opts.Direction {
		additionalFields = append(additionalFields, fmt.Sprintf("direction=%s", directionFilter))
	}

	if len(additionalFields) > 0 {
		request += "?"
		for _, field := range additionalFields {
			request += fmt.Sprintf("%s&", field)
		}
	}
	//remove the trailing '&'
	request = strings.TrimSuffix(request, "&")

	httpResponse, err := s.get2(ctx, request)
	if err != nil {
		return nil, err
	}

	if httpResponse.contentType != ContentTypeJSON {
		return nil, fmt.Errorf("unexpected content type %v (expected JSON)", httpResponse.contentType)
	}
	//data, meta, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), apiv1.Peers{})
	err = json.Unmarshal(httpResponse.body, &peerResonse)
	if err != nil {
		return nil, err
	}

	return &api.Response[*apiv1.Peers]{
		Data:     &apiv1.Peers{Peers: peerResonse.Data},
		Metadata: peerResonse.Meta,
	}, nil

}
