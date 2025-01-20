0.23.1:
  - add ability to override individual provider functions in mock client

0.23.0:
  - add attester_slashing, block_gossip, bls_to_execution_change and proposer_slashing events
  - add AttestationRewards, BlockRewards, and SyncCommitteeRewards functions

0.21.10:
  - better validator state when balance not supplied

0.21.9:
  - enable custom timeouts for POSTs

0.21.8:
  - remove Lodestar proposals workaround
  - add client headers for events stream

0.21.7:
  - use POST for specific validator and validator balance information

0.21.6:
  - use SSZ on a per-call basis

0.21.5:
  - ensure POST bodies are logged as JSON

0.21.4:
  - additional nil checks
  - allow non-mainnet configurations

0.21.3:
  - relax requirement for proposals to use the graffiti we request

0.21.2:
  - fuzz testing fixes

0.21.1:
  - fix potential crash when unmarshaling Gwei values
  - add `WithReducedMemoryUsage()` option for http service
  - more consistent tracing attributes and codes

0.21.0:
  - use v3 of the endpoint to obtain proposals
  - add bounds checking for ValidatorState

0.20.0:
  - allow delayed start of client, enabling the service even if the underlying beacon node is not ready
  - add IsActive() and IsSynced() methods to understand the status of the service
  - update multi clients to be aware of delayed start, only using clients that are synced
  - use standard errors for common function issues
  - add ProposerIndex() to VersionedSignedProposal
  - add name to multi clients to differentiate multiple instances
  - fully parse provided client URLs, allowing pass through of username, password, etc.

0.19.10:
  - add proposer_slashing and attester_slashing events
  - add bls_to_execution_change event

0.19.8
  - more efficient fetching for large numbers of validators

0.19.7:
  - add endpoint metrics for prometheus

0.19.5:
  - standardise names of options
  - add common options (currently just timeout) to options structs

0.19.4:
  - revert SubmitProposal() to use v1 of the API

0.19.0:
  - major rework of API; see docs/0.19-changes.md for details

0.18.3:
  - do not crash if beacon state is unavailable

0.18.2:
  - add 'withdrawable done' state to validators
  - use JSON metadata if not present in HTTP header

0.18.1:
  - add blinded block contents
  - add helpers to versioned signed blinded beacon block
  - add debug forkchoice endpoint support
  - add ProposerIndex() to BlindedBlocks
  - add helpers to versioned signed blinded beacon block
  - add BlockHash() to versioned signed beacon block
  - add ExecutionBlockHash() to versioned signed beacon block
  - rename data gas fields to blob gas for 1.4.0-beta1
 
0.18.0:
  - support Graffiti, ProposerIndex and RandaoReveal on VersionedBeaconBlock
  - use SSZ instead of JSON where available

0.17.0:
  - reworked JSON parsing for custom types to make easier to transition to another parser in future
  - added Deneb spec types
