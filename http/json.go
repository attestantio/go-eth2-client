package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/huandu/go-clone"
)

func decodeJSONResponse[T any](body io.Reader, res T) (T, map[string]any, error) {
	if body == nil {
		return res, nil, errors.New("no body to read")
	}

	decoded := make(map[string]json.RawMessage)

	if err := json.NewDecoder(body).Decode(&decoded); err != nil {
		return res, nil, errors.Join(errors.New("failed to parse JSON"), err)
	}

	data := clone.Clone(res).(T)
	metadata := make(map[string]any)
	for k, v := range decoded {
		switch k {
		case "data":
			err := json.Unmarshal(v, &data)
			if err != nil {
				return res, nil, errors.Join(errors.New("failed to unmarshal data"), err)
			}
		case "dependent_root":
			var val phase0.Root
			err := json.Unmarshal(v, &val)
			if err != nil {
				return res, nil, errors.Join(errors.New("failed to unmarshal dependent root"), err)
			}
			metadata[k] = val
		default:
			var val any
			err := json.Unmarshal(v, &val)
			if err != nil {
				return res, nil, errors.Join(fmt.Errorf("failed to unmarshal metadata %s", k), err)
			}
			metadata[k] = val
		}
	}

	return data, metadata, nil
}
