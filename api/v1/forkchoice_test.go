package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkChoiceJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		err      string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.forkChoiceJSON",
		},
		{
			name:     "Good",
			input:    []byte(`{"justified_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}]}`),
			expected: `{"justified_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}]}`,
			err:      "",
		},
		{
			name:  "JustifiedCheckpointMissing",
			input: []byte(`{"finalized_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`),
			err:   "justified checkpoint missing",
		},
		{
			name:  "JustifiedCheckpointInvalid",
			input: []byte(`{"finalized_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"justified_checkpoint":-1,"fork_choice_nodes":[]}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal number into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "FinalizedCheckpointMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`),
			err:   "finalized checkpoint missing",
		},
		{
			name:  "FinalizedCheckpointInvalid",
			input: []byte(`{"justified_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":-1,"fork_choice_nodes":[]}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal number into Go value of type phase0.checkpointJSON",
		},
		{
			name:     "ForkChoiceNodesInvalid",
			input:    []byte(`{"justified_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":-1}`),
			expected: `{"justified_checkpoint":{"epoch":"1","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x0100000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`,
			err:      "invalid JSON: json: cannot unmarshal number into Go struct field forkChoiceJSON.fork_choice_nodes of type []*v1.ForkChoiceNode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var fc api.ForkChoice
			err := json.Unmarshal(test.input, &fc)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&fc)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), fc.String())
			}
		})
	}
}

func TestForkChoiceNodeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		err      string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.forkChoiceNodeJSON",
		},
		{
			name:     "Good",
			input:    []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			expected: `{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`,
			err:      "",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":1,"block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field forkChoiceNodeJSON.slot of type string",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"slot":"1962336","block_root":"","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "incorrect length 0 for block root",
		},
		{
			name:  "JustifiedEpochInvalid",
			input: []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":-1,"finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field forkChoiceNodeJSON.justified_epoch of type string",
		},
		{
			name:  "FinalizedEpochInvalid",
			input: []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":-1,"weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field forkChoiceNodeJSON.finalized_epoch of type string",
		},
		{
			name:  "WeightInvalid",
			input: []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":400,"validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field forkChoiceNodeJSON.weight of type string",
		},
		{
			name:  "ValidityInvalid",
			input: []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"NOT_A_VALID_VALIDITY","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid value for validity: NOT_A_VALID_VALIDITY: unrecognised fork choice validity: NOT_A_VALID_VALIDITY",
		},
		{
			name:  "ExecutionBlockHashInvalid",
			input: []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"abc","extra_data":{"justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7","state_root":"0x4bcecf56081291ab95df1dba25b0f83343d38217e9cec198510b61f6f35afdb3","unrealised_finalized_epoch":"61321","unrealised_justified_epoch":"61322","unrealized_finalized_root":"0x57a41f26678190d3e319c19fe9f4ea3830c4b21710a1e1ae41adcc23d0f030a2","unrealized_justified_root":"0xdee6c83ee7dc6c0916a8d43c4e7cda93655857da0487f193a62852699e5c39f7"}}`),
			err:   "invalid value for execution block hash: abc: encoding/hex: odd length hex string",
		},
		{
			name:     "ExtraDataMissing",
			input:    []byte(`{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9"}`),
			expected: `{"slot":"1962336","block_root":"0x0f61e82f7b51f41fcd552cbdc64547bd9e1ba54b8404482732927645c2e13ec6","parent_root":"0xe399a2ee74cf0570b4f980772983bdc5cfbfde1f87f3ab395f2bee96103978c7","justified_epoch":"61322","finalized_epoch":"61321","weight":"57481550000000","validity":"valid","execution_block_hash":"0x06a0277e02eae44c332bcec82d6715c3113dddce427982014cf5f43432f479e9"}`,
			err:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var fc api.ForkChoiceNode
			err := json.Unmarshal(test.input, &fc)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&fc)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), fc.String())
			}
		})
	}
}
