package verkle

const IPA_PROOF_DEPTH = 8

type IPAProof struct {
	CL              [IPA_PROOF_DEPTH][32]byte `ssz-size:"256"`
	CR              [IPA_PROOF_DEPTH][32]byte `ssz-size:"256"`
	FinalEvaluation [32]byte                  `ssz-size:"32"`
}

type VerkleProof struct {
	OtherStems            [][31]byte `json:"otherStems"`
	DepthExtensionPresent []byte     `json:"depthExtensionPresent"`
	CommitmentsByPath     [][32]byte `json:"commitmentsByPath"`
	D                     [32]byte   `json:"d"`
	IPAProof              *IPAProof  `json:"ipa_proof"`
}

type SuffixStateDiff struct {
	Suffix       byte      `json:"suffix"`
	CurrentValue *[32]byte `json:"currentValue"`
	NewValue     *[32]byte `json:"newValue"`
}
type SuffixStateDiffs []SuffixStateDiff

type StemStateDiff struct {
	Stem        [31]byte         `json:"stem"`
	SuffixDiffs SuffixStateDiffs `json:"suffixDiffs"`
}

type StateDiff []StemStateDiff
