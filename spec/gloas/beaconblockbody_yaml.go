// Copyright © 2023 Attestant Limited.
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

package gloas

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconBlockBody) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconBlockBodyJSON{
		RANDAOReveal:              fmt.Sprintf("%#x", b.RANDAOReveal),
		ETH1Data:                  b.ETH1Data,
		Graffiti:                  fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings:         b.ProposerSlashings,
		AttesterSlashings:         b.AttesterSlashings,
		Attestations:              b.Attestations,
		Deposits:                  b.Deposits,
		VoluntaryExits:            b.VoluntaryExits,
		SyncAggregate:             b.SyncAggregate,
		BLSToExecutionChanges:     b.BLSToExecutionChanges,
		SignedExecutionPayloadBid: b.SignedExecutionPayloadBid,
		PayloadAttestations:       b.PayloadAttestations,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalYAML(input []byte) error {
	var data beaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(marshaled)
}
