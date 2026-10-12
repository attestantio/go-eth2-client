dev:
  - add Gloas consensus-spec vector CI and weekly release support review issues

0.30.0:
   - support Gloas consensus types and versioned helpers
   - add execution payload envelope retrieval and bid/envelope submission APIs
   - add payload timeliness committee duties and payload attestation APIs
   - add MultiForkProposalProvider for proposal production across the Gloas fork
   - support Gloas attestation hashing, submission and payload-availability votes
   - fix attestation-pool decoding and validation across forks
   - honor custom SSZ presets when submitting proposals
   - add mock SubmitBLSToExecutionChanges
   - handle invalid DataVersion values without panicking
   - bound HTTP response sizes for POSTs and Gloas fetch endpoints
   - add SubmitProposerPreferences with JSON and SSZ support in http, multi and mock clients
   - update ePBS block production to POST /eth/v4/validator/blocks/{slot} with a required EPBSProposalOpts.BuilderConfig
     - add Gloas BuilderConfig, BuilderEntry, BuilderRequestAuth and SignedBuilderRequestAuth types with JSON, YAML and SSZ support
   - add VersionedEPBSProposal.BuilderIndex and SubmitProposalOpts.BuilderURL for direct-builder proposal routing
   - validate ePBS bid policy and execution-value metadata, including minimum bids and execution-payment caps
   - add head_v2 events and HeadV2Handler, including fork version, payload status and epoch-dependent roots
   - add light_client_finality_update and light_client_optimistic_update events through the generic event handler
   - add ProposerDutiesV2Provider and ProposerDutiesV2 in http, multi and mock clients
     - use /eth/v2/validator/duties/proposer/{epoch} and its dependent-root semantics with head_v2; existing v1 behavior is unchanged
   - fix multi event subscriptions to preserve topics and handlers, forward only the active client's events, and continue retrying failed subscriptions
   - encode empty Gloas execution-payload transactions as "0x" in JSON and YAML, and accept them on JSON decode
  - add Gloas events: execution_payload, execution_payload_gossip, execution_payload_available, execution_payload_bid, payload_attestation_message, proposer_preferences and fast_confirmation
  - add http.ValidateEventsOpts, the events options check shared by the http and multi clients
  - Event.UnmarshalJSON decodes attester_slashing data as electra.AttesterSlashing rather than phase0.AttesterSlashing, matching the events stream
  - multi Events() validates its options as http Events() does, and returns an error if no client is subscribed or awaiting a retry
  - Events() no longer checks topics against SupportedEventTopics: removing a topic from the map no longer refuses subscriptions to it, and adding one no longer lets it through
  - multi Events() waits one retry interval before retrying a client that failed to subscribe, rather than retrying at once; set the interval, 5s by default, with WithEventsRetryInterval

0.29.0:
  - use dynssz library for SSZ handling

0.28.1:
  - update to HTTP tests
  - add beacon committee selections endpoint for distributed validators

0.28.0:
  - update dependency for go-bitfield to get from offchain instead of prysmatic org

0.27.2:
  - add /eth/v1/beacon/blobs/ endpoint support
  - update go version to 1.25
  - updated spec conformance tests for latest forks
  - updated linters and corresponding fixes

0.27.1:
  - set max possible blob count to 72

0.27.0:
  - support fulu
    - introduce data column sidecar event api and corresponding event handler
    - add blockcontents and signedblockcontents for fulu api
    - add beaconstate container updates as per spec
    - add fulu cases for all versioned spec, versioned api and http functions

0.26.0:
  - refactor http.Spec to allow more complex types in the keys
  - support pending consolidations and deposits

0.25.2:
  - add multi/submitblindedproposal

0.25.1:
  - add Merkle tree and proof generation utils
  - add convenience methods on versioned beaconstate for field access and proofs:
    - ValidatorAtIndex
    - ValidatorBalance
    - FieldIndex
    - FieldGeneralizedIndex
    - FieldRoot
    - FieldTree
    - ProveField
    - VerifyFieldProof

0.25.0:
  - update attestation pool endpoint to receive versioned attestations

0.24.2:
  - support single_attestation event
  - support change to attestation event; this event now emits a spec.VersionedAttestation
  - support change to attester_slashing event; this event now emits an electra.AttesterSlashing
  - update Events endpoint to provide specific handlers for each event

0.24.0:
  - support electra
    - the most notable change is that a number of functions now use spec.VersionedAttestation in place of phase0.Attestation
    - this release uses a number of new beacon API endpoints, specifically:
      - /eth/v2/validator/aggregate_attestation
      - /eth/v2/validator/aggregate_and_proofs
      - /eth/v2/beacon/pool/attestations
      These endpoints are supported in all current releases of major beacon nodes at the time of release

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
