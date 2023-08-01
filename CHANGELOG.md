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
