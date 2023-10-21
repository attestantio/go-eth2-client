package http

import (
	"encoding/json"
	"io"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/huandu/go-clone"
	"github.com/pkg/errors"
)

func decodeJSONResponse[T any](body io.Reader, res T) (T, map[string]any, error) {
	if body == nil {
		return res, nil, errors.New("no body to read")
	}

	decoded := make(map[string]json.RawMessage)

	if err := json.NewDecoder(body).Decode(&decoded); err != nil {
		return res, nil, errors.Wrap(err, "failed to parse JSON")
	}

	data := clone.Clone(res).(T)
	metadata := make(map[string]any)
	for k, v := range decoded {
		switch k {
		case "data":
			err := json.Unmarshal(v, &data)
			if err != nil {
				return res, nil, errors.Wrap(err, "failed to unmarshal data")
			}
		case "dependent_root":
			var val phase0.Root
			err := json.Unmarshal(v, &val)
			if err != nil {
				return res, nil, errors.Wrap(err, "failed to unmarshal dependent root")
			}
			metadata[k] = val
		default:
			var val any
			err := json.Unmarshal(v, &val)
			if err != nil {
				return res, nil, errors.Wrapf(err, "failed to unmarshal metadata %s", k)
			}
			metadata[k] = val
		}
	}

	return data, metadata, nil
}
