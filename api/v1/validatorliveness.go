package v1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ValidatorLiveness represents the observed liveness state of a validator.
type ValidatorLiveness struct {
	// Index is the validator index.
	Index phase0.ValidatorIndex
	// IsLive indicates whether the validator is live in the given epoch.
	IsLive bool
}

// validatorLivenessJSON is the spec representation of the struct.
type validatorLivenessJSON struct {
	Index  string `json:"index"`
	IsLive bool   `json:"is_live"`
}

// MarshalJSON implements json.Marshaler.
func (v *ValidatorLiveness) MarshalJSON() ([]byte, error) {
	return json.Marshal(&validatorLivenessJSON{
		Index:  fmt.Sprintf("%d", v.Index),
		IsLive: v.IsLive,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *ValidatorLiveness) UnmarshalJSON(input []byte) error {
	var validatorLivenessJSON validatorLivenessJSON
	if err := json.Unmarshal(input, &validatorLivenessJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Convert Index from string to phase0.ValidatorIndex
	index, err := strconv.ParseUint(validatorLivenessJSON.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}

	v.Index = phase0.ValidatorIndex(index)
	v.IsLive = validatorLivenessJSON.IsLive

	return nil
}

// String returns a string version of the structure.
func (v *ValidatorLiveness) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
