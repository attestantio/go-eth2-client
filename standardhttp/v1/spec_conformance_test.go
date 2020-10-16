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

package v1_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v1"
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
		"DOMAIN_AGGREGATE_AND_PROOF":            []byte{},
		"DOMAIN_BEACON_ATTESTER":                []byte{},
		"DOMAIN_BEACON_PROPOSER":                []byte{},
		"DOMAIN_DEPOSIT":                        []byte{},
		"DOMAIN_RANDAO":                         []byte{},
		"DOMAIN_SELECTION_PROOF":                []byte{},
		"DOMAIN_VOLUNTARY_EXIT":                 []byte{},
		"EFFECTIVE_BALANCE_INCREMENT":           uint64(0),
		"EJECTION_BALANCE":                      uint64(0),
		"EPOCHS_PER_ETH1_VOTING_PERIOD":         uint64(0),
		"EPOCHS_PER_HISTORICAL_VECTOR":          uint64(0),
		"EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION": uint64(0),
		"EPOCHS_PER_SLASHINGS_VECTOR":           uint64(0),
		"ETH1_FOLLOW_DISTANCE":                  uint64(0),
		"GENESIS_DELAY":                         uint64(0),
		"GENESIS_FORK_VERSION":                  []byte{},
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
		"MIN_ATTESTATION_INCLUSION_DELAY":       uint64(0),
		"MIN_DEPOSIT_AMOUNT":                    uint64(0),
		"MIN_EPOCHS_TO_INACTIVITY_PENALTY":      uint64(0),
		"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT":    uint64(0),
		"MIN_GENESIS_TIME":                      uint64(0),
		"MIN_PER_EPOCH_CHURN_LIMIT":             uint64(0),
		"MIN_SEED_LOOKAHEAD":                    uint64(0),
		"MIN_SLASHING_PENALTY_QUOTIENT":         uint64(0),
		"MIN_VALIDATOR_WITHDRAWABILITY_DELAY":   uint64(0),
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

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithTimeout(timeout),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	fmt.Printf("Testing spec endpoint\n")
	spec, err := service.Spec(context.Background())
	require.NoError(t, err)
	require.NotNil(t, spec)

	// Ensure everything we expect is present.
	for k, v := range expected {
		val, exists := spec[k]
		if !exists {
			fmt.Printf("  Expected value %s not present\n", k)
			continue
		}

		if reflect.TypeOf(v) != reflect.TypeOf(val) {
			fmt.Printf("  Value %s has incorrect type\n", k)
		}
	}

	// Ensure nothing we don't expect is present.
	for k := range spec {
		_, exists := expected[k]
		if !exists {
			fmt.Printf("  Unexpected value %s is present\n", k)
		}
	}
}
