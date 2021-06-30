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

// +build conformance

package http_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/http"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecConformance(t *testing.T) {
	expected := map[string]interface{}{
		"BASE_REWARD_FACTOR":                    uint64(0),
		"BLS_WITHDRAWAL_PREFIX":                 []byte{},
		"CHURN_LIMIT_QUOTIENT":                  uint64(0),
		"CONFIG_NAME":                           "",
		"DEPOSIT_CHAIN_ID":                      uint64(0),
		"DEPOSIT_CONTRACT_ADDRESS":              []byte{},
		"DEPOSIT_NETWORK_ID":                    uint64(0),
		"DOMAIN_AGGREGATE_AND_PROOF":            spec.DomainType{},
		"DOMAIN_BEACON_ATTESTER":                spec.DomainType{},
		"DOMAIN_BEACON_PROPOSER":                spec.DomainType{},
		"DOMAIN_DEPOSIT":                        spec.DomainType{},
		"DOMAIN_RANDAO":                         spec.DomainType{},
		"DOMAIN_SELECTION_PROOF":                spec.DomainType{},
		"DOMAIN_VOLUNTARY_EXIT":                 spec.DomainType{},
		"EFFECTIVE_BALANCE_INCREMENT":           uint64(0),
		"EJECTION_BALANCE":                      uint64(0),
		"EPOCHS_PER_ETH1_VOTING_PERIOD":         uint64(0),
		"EPOCHS_PER_HISTORICAL_VECTOR":          uint64(0),
		"EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION": uint64(0),
		"EPOCHS_PER_SLASHINGS_VECTOR":           uint64(0),
		"ETH1_FOLLOW_DISTANCE":                  uint64(0),
		"GENESIS_DELAY":                         time.Duration(0),
		"GENESIS_FORK_VERSION":                  spec.Version{},
		"HISTORICAL_ROOTS_LIMIT":                uint64(0),
		"HYSTERESIS_DOWNWARD_MULTIPLIER":        uint64(0),
		"HYSTERESIS_QUOTIENT":                   uint64(0),
		"HYSTERESIS_UPWARD_MULTIPLIER":          uint64(0),
		"INACTIVITY_PENALTY_QUOTIENT":           uint64(0),
		"MAX_ATTESTATIONS":                      uint64(0),
		"MAX_ATTESTER_SLASHINGS":                uint64(0),
		"MAX_COMMITTEES_PER_SLOT":               uint64(0),
		"MAX_DEPOSITS":                          uint64(0),
		"MAX_EFFECTIVE_BALANCE":                 uint64(0),
		"MAX_PROPOSER_SLASHINGS":                uint64(0),
		"MAX_SEED_LOOKAHEAD":                    uint64(0),
		"MAX_VALIDATORS_PER_COMMITTEE":          uint64(0),
		"MAX_VOLUNTARY_EXITS":                   uint64(0),
		"MIN_ATTESTATION_INCLUSION_DELAY":       time.Duration(0),
		"MIN_DEPOSIT_AMOUNT":                    uint64(0),
		"MIN_EPOCHS_TO_INACTIVITY_PENALTY":      uint64(0),
		"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT":    uint64(0),
		"MIN_GENESIS_TIME":                      time.Time{},
		"MIN_PER_EPOCH_CHURN_LIMIT":             uint64(0),
		"MIN_SEED_LOOKAHEAD":                    uint64(0),
		"MIN_SLASHING_PENALTY_QUOTIENT":         uint64(0),
		"MIN_VALIDATOR_WITHDRAWABILITY_DELAY":   time.Duration(0),
		"PROPORTIONAL_SLASHING_MULTIPLIER":      uint64(0),
		"PROPOSER_REWARD_QUOTIENT":              uint64(0),
		"RANDOM_SUBNETS_PER_VALIDATOR":          uint64(0),
		"SAFE_SLOTS_TO_UPDATE_JUSTIFIED":        uint64(0),
		"SECONDS_PER_ETH1_BLOCK":                time.Duration(0),
		"SECONDS_PER_SLOT":                      time.Duration(0),
		"SHARD_COMMITTEE_PERIOD":                uint64(0),
		"SHUFFLE_ROUND_COUNT":                   uint64(0),
		"SLOTS_PER_EPOCH":                       uint64(0),
		"SLOTS_PER_HISTORICAL_ROOT":             uint64(0),
		"TARGET_AGGREGATORS_PER_COMMITTEE":      uint64(0),
		"TARGET_COMMITTEE_SIZE":                 uint64(0),
		"VALIDATOR_REGISTRY_LIMIT":              uint64(0),
		"WHISTLEBLOWER_REWARD_QUOTIENT":         uint64(0),
	}

	service, err := http.New(context.Background(),
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	spec, err := service.Spec(context.Background())
	require.NoError(t, err)
	require.NotNil(t, spec)

	// Ensure everything we expect is present.
	for k, v := range expected {
		val, exists := spec[k]
		assert.True(t, exists, fmt.Sprintf("Value %s not present", k))
		assert.Equal(t, reflect.TypeOf(v), reflect.TypeOf(val), fmt.Sprintf("Value %s has incorrect type", k))
	}

	// Ensure nothing we don't expect is present.
	for k := range spec {
		_, exists := expected[k]
		assert.True(t, exists, fmt.Sprintf("Value %s unexpected", k))
	}
}
