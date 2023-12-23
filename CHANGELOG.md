dev:
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
