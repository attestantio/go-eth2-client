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

package http

import client "github.com/attestantio/go-eth2-client"

// Compile-time assertions ensure that Service continues to satisfy its intended client interfaces.
var (
	_ client.Service                               = (*Service)(nil)
	_ client.EpochFromStateIDProvider              = (*Service)(nil)
	_ client.SlotFromStateIDProvider               = (*Service)(nil)
	_ client.SlotDurationProvider                  = (*Service)(nil)
	_ client.SlotsPerEpochProvider                 = (*Service)(nil)
	_ client.FarFutureEpochProvider                = (*Service)(nil)
	_ client.TargetAggregatorsPerCommitteeProvider = (*Service)(nil)
	_ client.SignedBeaconBlockProvider             = (*Service)(nil)
	_ client.BlobsProvider                         = (*Service)(nil)
	_ client.BlobSidecarsProvider                  = (*Service)(nil)
	_ client.BeaconCommitteesProvider              = (*Service)(nil)
	_ client.SyncCommitteesProvider                = (*Service)(nil)
	_ client.AggregateAttestationProvider          = (*Service)(nil)
	_ client.AggregateAttestationsSubmitter        = (*Service)(nil)
	_ client.AttestationDataProvider               = (*Service)(nil)
	_ client.AttestationPoolProvider               = (*Service)(nil)
	_ client.AttestationRewardsProvider            = (*Service)(nil)
	_ client.AttestationsSubmitter                 = (*Service)(nil)
	_ client.AttesterSlashingSubmitter             = (*Service)(nil)
	_ client.AttesterDutiesProvider                = (*Service)(nil)
	_ client.BlockRewardsProvider                  = (*Service)(nil)
	_ client.DepositContractProvider               = (*Service)(nil)
	_ client.SyncCommitteeDutiesProvider           = (*Service)(nil)
	_ client.SyncCommitteeMessagesSubmitter        = (*Service)(nil)
	_ client.SyncCommitteeSubscriptionsSubmitter   = (*Service)(nil)
	_ client.SyncCommitteeContributionProvider     = (*Service)(nil)
	_ client.SyncCommitteeContributionsSubmitter   = (*Service)(nil)
	_ client.SyncCommitteeRewardsProvider          = (*Service)(nil)
	_ client.BLSToExecutionChangesSubmitter        = (*Service)(nil)
	_ client.BeaconBlockHeadersProvider            = (*Service)(nil)
	_ client.ProposalProvider                      = (*Service)(nil)
	_ client.ProposalSlashingSubmitter             = (*Service)(nil)
	_ client.BeaconBlockRootProvider               = (*Service)(nil)
	_ client.BeaconBlockSubmitter                  = (*Service)(nil)
	_ client.ProposalSubmitter                     = (*Service)(nil)
	_ client.BeaconCommitteeSubscriptionsSubmitter = (*Service)(nil)
	_ client.BeaconCommitteeSelectionsProvider     = (*Service)(nil)
	_ client.BeaconStateProvider                   = (*Service)(nil)
	_ client.BeaconStateRandaoProvider             = (*Service)(nil)
	_ client.BeaconStateRootProvider               = (*Service)(nil)
	_ client.BlindedBeaconBlockSubmitter           = (*Service)(nil)
	_ client.BlindedProposalSubmitter              = (*Service)(nil)
	_ client.ValidatorRegistrationsSubmitter       = (*Service)(nil)
	_ client.EventsProvider                        = (*Service)(nil)
	_ client.FinalityProvider                      = (*Service)(nil)
	_ client.ForkChoiceProvider                    = (*Service)(nil)
	_ client.ForkProvider                          = (*Service)(nil)
	_ client.ForkScheduleProvider                  = (*Service)(nil)
	_ client.GenesisProvider                       = (*Service)(nil)
	_ client.NodePeersProvider                     = (*Service)(nil)
	_ client.NodeSyncingProvider                   = (*Service)(nil)
	_ client.ValidatorLivenessProvider             = (*Service)(nil)
	_ client.NodeVersionProvider                   = (*Service)(nil)
	_ client.ProposalPreparationsSubmitter         = (*Service)(nil)
	_ client.ProposerDutiesProvider                = (*Service)(nil)
	_ client.SpecProvider                          = (*Service)(nil)
	_ client.ValidatorBalancesProvider             = (*Service)(nil)
	_ client.ValidatorsProvider                    = (*Service)(nil)
	_ client.VoluntaryExitSubmitter                = (*Service)(nil)
	_ client.VoluntaryExitPoolProvider             = (*Service)(nil)
	_ client.PendingDepositProvider                = (*Service)(nil)
	_ client.PendingConsolidationsProvider         = (*Service)(nil)
	_ client.PendingPartialWithdrawalsProvider     = (*Service)(nil)
	_ client.DomainProvider                        = (*Service)(nil)
	//nolint:staticcheck // GenesisTimeProvider is deprecated but still implemented.
	_ client.GenesisTimeProvider = (*Service)(nil)
	_ client.NodeClientProvider  = (*Service)(nil)
)
