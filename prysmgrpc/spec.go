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

package prysmgrpc

import (
	"context"
	"strconv"
	"strings"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type specJSON struct {
	Data map[string]string `json:"data"`
}

var specNameMapping = map[string]string{
	"BaseRewardFactor":                 "BASE_REWARD_FACTOR",
	"ChurnLimitQuotient":               "CHURN_LIMIT_QUOTIENT",
	"DomainAggregateAndProof":          "DOMAIN_AGGREGATE_AND_PROOF",
	"DomainBeaconAttester":             "DOMAIN_BEACON_ATTESTER",
	"DomainBeaconProposer":             "DOMAIN_BEACON_PROPOSER",
	"DomainDeposit":                    "DOMAIN_DEPOSIT",
	"DomainRandao":                     "DOMAIN_RANDAO",
	"DomainSelectionProof":             "DOMAIN_SELECTION_PROOF",
	"DomainVoluntaryExit":              "DOMAIN_VOLUNTARY_EXIT",
	"EffectiveBalanceIncrement":        "EFFECTIVE_BALANCE_INCREMENT",
	"EjectionBalance":                  "EJECTION_BALANCE",
	"EpochsPerEth1VotingPeriod":        "EPOCHS_PER_ETH1_VOTING_PERIOD",
	"EpochsPerHistoricalVector":        "EPOCHS_PER_HISTORICAL_VECTOR",
	"EpochsPerSlashingsVector":         "EPOCHS_PER_SLASHINGS_VECTOR",
	"Eth1FollowDistance":               "ETH1_FOLLOW_DISTANCE",
	"GenesisDelay":                     "GENESIS_DELAY",
	"GenesisForkVersion":               "GENESIS_FORK_VERSION",
	"HistoricalRootsLimit":             "HISTORICAL_ROOTS_LIMIT",
	"HysteresisDownwardMultiplier":     "HYSTERESIS_DOWNWARD_MULTIPLIER",
	"HysteresisQuotient":               "HYSTERESIS_QUOTIENT",
	"HysteresisUpwardMultiplier":       "HYSTERESIS_UPWARD_MULTIPLIER",
	"InactivityPenaltyQuotient":        "INACTIVITY_PENALTY_QUOTIENT",
	"MaxAttestations":                  "MAX_ATTESTATIONS",
	"MaxAttesterSlashings":             "MAX_ATTESTER_SLASHINGS",
	"MaxCommitteesPerSlot":             "MAX_COMMITTEES_PER_SLOT",
	"MaxDeposits":                      "MAX_DEPOSITS",
	"MaxEffectiveBalance":              "MAX_EFFECTIVE_BALANCE",
	"MaxProposerSlashings":             "MAX_PROPOSER_SLASHINGS",
	"MaxSeedLookahead":                 "MAX_SEED_LOOKAHEAD",
	"MaxValidatorsPerCommittee":        "MAX_VALIDATORS_PER_COMMITTEE",
	"MaxVoluntaryExits":                "MAX_VOLUNTARY_EXITS",
	"MinAttestationInclusionDelay":     "MIN_ATTESTATION_INCLUSION_DELAY",
	"MinDepositAmount":                 "MIN_DEPOSIT_AMOUNT",
	"MinEpochsToInactivityPenalty":     "MIN_EPOCHS_TO_INACTIVITY_PENALTY",
	"MinGenesisActiveValidatorCount":   "MIN_GENESIS_ACTIVE_VALIDATOR_COUNT",
	"MinGenesisTime":                   "MIN_GENESIS_TIME",
	"MinPerEpochChurnLimit":            "MIN_PER_EPOCH_CHURN_LIMIT",
	"MinSeedLookahead":                 "MIN_SEED_LOOKEAHEAD",
	"MinSlashingPenaltyQuotient":       "MIN_SLASHING_PENALTY_QUOTIENT",
	"MinValidatorWithdrawabilityDelay": "MIN_VALIDATOR_WITHDRAWABILITY_DELAY",
	"ProportionalSlashingMultiplier":   "PROPORTIONAL_SLASHING_MULTIPLIER",
	"ProposerRewardQuotient":           "PROPOSER_REWARD_QUOTIENT",
	"SafeSlotsToUpdateJustified":       "SAFE_SLOTS_TO_UPDATE_JUSTIFIED",
	"SecondsPerETH1Block":              "SECONDS_PER_ETH1_BLOCK",
	"SecondsPerSlot":                   "SECONDS_PER_SLOT",
	"ShardCommitteePeriod":             "SHARD_COMMITTEE_PERIOD",
	"ShuffleRoundCount":                "SHUFFLE_ROUND_COUNT",
	"SlotsPerEpoch":                    "SLOTS_PER_EPOCH",
	"SlotsPerHistoricalRoot":           "SLOTS_PER_HISTORICAL_ROOT",
	"TargetAggregatorsPerCommittee":    "TARGET_AGGREGATORS_PER_COMMITTEE",
	"TargetCommitteeSize":              "TARGET_COMMITTEE_SIZE",
	"ValidatorRegistryLimit":           "VALIDATOR_REGISTRY_LIMIT",
	"WhistleBlowerRewardQuotient":      "WHISTLEBLOWER_REWARD_QUOTIENT",
}

// Spec provides the spec information of the chain.
func (s *Service) Spec(ctx context.Context) (map[string]interface{}, error) {
	if s.spec == nil {
		conn := ethpb.NewBeaconChainClient(s.conn)
		log.Trace().Msg("Fetching beacon chain spec")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		config, err := conn.GetBeaconConfig(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}

		res := make(map[string]interface{})
		for k, v := range config.Config {
			// Ensure we know about this value, and map it to the official name.
			if _, exists := specNameMapping[k]; !exists {
				continue
			}
			k = specNameMapping[k]
			// Handle domain types.
			if strings.HasPrefix(k, "DOMAIN_") {
				byteArrayVal, err := parseConfigByteArray(v)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse domain type")
				}
				var domain spec.DomainType
				copy(domain[:], byteArrayVal)
				res[k] = domain
				continue
			}

			// Handle durations.
			if strings.HasPrefix(k, "SECONDS_PER_") {
				intVal, err := strconv.ParseUint(v, 10, 64)
				if err == nil && intVal != 0 {
					res[k] = time.Duration(intVal) * time.Second
					continue
				}
			}

			// Handle integers.
			if v == "0" {
				res[k] = uint64(0)
				continue
			}
			intVal, err := strconv.ParseUint(v, 10, 64)
			if err == nil && intVal != 0 {
				res[k] = intVal
				continue
			}

			// Assume string.
			res[k] = v
		}
		s.spec = res
	}
	return s.spec, nil
}
