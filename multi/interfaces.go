// Copyright © 2026 Attestant Limited.
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

package multi

import consensusclient "github.com/attestantio/go-eth2-client"

// Compile-time assertions ensure that Service continues to satisfy its intended client interfaces.
var (
	_ consensusclient.Service                               = (*Service)(nil)
	_ consensusclient.SlotDurationProvider                  = (*Service)(nil)
	_ consensusclient.SlotsPerEpochProvider                 = (*Service)(nil)
	_ consensusclient.FarFutureEpochProvider                = (*Service)(nil)
	_ consensusclient.TargetAggregatorsPerCommitteeProvider = (*Service)(nil)
	_ consensusclient.SignedBeaconBlockProvider             = (*Service)(nil)
	_ consensusclient.BlobsProvider                         = (*Service)(nil)
	_ consensusclient.BlobSidecarsProvider                  = (*Service)(nil)
	_ consensusclient.BeaconCommitteesProvider              = (*Service)(nil)
	_ consensusclient.SyncCommitteesProvider                = (*Service)(nil)
	_ consensusclient.AggregateAttestationProvider          = (*Service)(nil)
	_ consensusclient.AggregateAttestationsSubmitter        = (*Service)(nil)
	_ consensusclient.AttestationDataProvider               = (*Service)(nil)
	_ consensusclient.AttestationPoolProvider               = (*Service)(nil)
	_ consensusclient.AttestationRewardsProvider            = (*Service)(nil)
	_ consensusclient.AttestationsSubmitter                 = (*Service)(nil)
	_ consensusclient.AttesterDutiesProvider                = (*Service)(nil)
	_ consensusclient.BlockRewardsProvider                  = (*Service)(nil)
	_ consensusclient.DepositContractProvider               = (*Service)(nil)
	_ consensusclient.SyncCommitteeDutiesProvider           = (*Service)(nil)
	_ consensusclient.SyncCommitteeMessagesSubmitter        = (*Service)(nil)
	_ consensusclient.SyncCommitteeSubscriptionsSubmitter   = (*Service)(nil)
	_ consensusclient.SyncCommitteeContributionProvider     = (*Service)(nil)
	_ consensusclient.SyncCommitteeContributionsSubmitter   = (*Service)(nil)
	_ consensusclient.SyncCommitteeRewardsProvider          = (*Service)(nil)
	_ consensusclient.BeaconBlockHeadersProvider            = (*Service)(nil)
	_ consensusclient.ProposalProvider                      = (*Service)(nil)
	_ consensusclient.BeaconBlockRootProvider               = (*Service)(nil)
	_ consensusclient.BeaconBlockSubmitter                  = (*Service)(nil)
	_ consensusclient.ProposalSubmitter                     = (*Service)(nil)
	_ consensusclient.BeaconCommitteeSubscriptionsSubmitter = (*Service)(nil)
	_ consensusclient.BeaconCommitteeSelectionsProvider     = (*Service)(nil)
	_ consensusclient.BeaconStateProvider                   = (*Service)(nil)
	_ consensusclient.BeaconStateRootProvider               = (*Service)(nil)
	_ consensusclient.BlindedBeaconBlockSubmitter           = (*Service)(nil)
	_ consensusclient.BlindedProposalSubmitter              = (*Service)(nil)
	_ consensusclient.ValidatorRegistrationsSubmitter       = (*Service)(nil)
	_ consensusclient.EventsProvider                        = (*Service)(nil)
	_ consensusclient.FinalityProvider                      = (*Service)(nil)
	_ consensusclient.ForkChoiceProvider                    = (*Service)(nil)
	_ consensusclient.ForkProvider                          = (*Service)(nil)
	_ consensusclient.ForkScheduleProvider                  = (*Service)(nil)
	_ consensusclient.GenesisProvider                       = (*Service)(nil)
	_ consensusclient.NodePeersProvider                     = (*Service)(nil)
	_ consensusclient.NodeSyncingProvider                   = (*Service)(nil)
	_ consensusclient.ValidatorLivenessProvider             = (*Service)(nil)
	_ consensusclient.NodeVersionProvider                   = (*Service)(nil)
	_ consensusclient.ProposalPreparationsSubmitter         = (*Service)(nil)
	_ consensusclient.ProposerDutiesProvider                = (*Service)(nil)
	_ consensusclient.SpecProvider                          = (*Service)(nil)
	_ consensusclient.ValidatorBalancesProvider             = (*Service)(nil)
	_ consensusclient.ValidatorsProvider                    = (*Service)(nil)
	_ consensusclient.VoluntaryExitSubmitter                = (*Service)(nil)
	_ consensusclient.VoluntaryExitPoolProvider             = (*Service)(nil)
	_ consensusclient.PendingDepositProvider                = (*Service)(nil)
	_ consensusclient.PendingConsolidationsProvider         = (*Service)(nil)
	_ consensusclient.PendingPartialWithdrawalsProvider     = (*Service)(nil)
	_ consensusclient.DomainProvider                        = (*Service)(nil)
	//nolint:staticcheck // GenesisTimeProvider is deprecated but still implemented.
	_ consensusclient.GenesisTimeProvider = (*Service)(nil)
)
