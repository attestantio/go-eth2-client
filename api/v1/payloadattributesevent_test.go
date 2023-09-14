package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestPayloadAttributesEventJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.payloadAttributesEventJSON",
		},
		{
			name:  "VersionMissing",
			input: []byte(`{"data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "unsupported data version",
		},
		{
			name:  "VersionInvalid",
			input: []byte(`{"version":"invalid","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: unrecognised data version \"invalid\"",
		},
		{
			name:  "ProposerIndexMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "proposer index missing",
		},
		{
			name:  "ProposerIndexWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":123,"proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.proposer_index of type string",
		},
		{
			name:  "ProposerIndexInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"invalid","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for proposer index: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ProposerSlotMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "proposal slot missing",
		},
		{
			name:  "ProposerSlotWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":10,"parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.proposal_slot of type string",
		},
		{
			name:  "ProposerSlotInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"invalid","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for proposal slot: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ParentBlockNumberMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block number missing",
		},
		{
			name:  "ParentBlockNumberWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":9,"parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.parent_block_number of type string",
		},
		{
			name:  "ParentBlockNumberInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"invalid","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block number: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ParentBlockRootMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block root missing",
		},
		{
			name:  "ParentBlockRootWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":true,"parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesDataJSON.data.parent_block_root of type string",
		},
		{
			name:  "ParentBlockRootInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"invalid","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ParentBlockHashMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block hash missing",
		},
		{
			name:  "ParentBlockHashWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":true,"payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesDataJSON.data.parent_block_hash of type string",
		},
		{
			name:  "ParentBlockHashInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"invalid","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BadPayloadAttributesV1Data",
			input: []byte(`{"version":"bellatrix","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV2Data",
			input: []byte(`{"version":"capella","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV3Data",
			input: []byte(`{"version":"deneb","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV1JSON",
		},
		{
			name:  "BadPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV2JSON",
		},
		{
			name:  "BadPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV3JSON",
		},
		{
			name:  "BadPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "BadPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "BadPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "MissingPayloadAttributes",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}}`),
			err:   "payload attributes missing",
		},
		{
			name:  "MissingWithdrawals",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "payload attributes withdrawals missing",
		},
		{
			name:  "GoodPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
		},
		{
			name:  "GoodPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}]}}}`),
		},
		{
			name:  "GoodPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.PayloadAttributesEvent
			err := json.Unmarshal(test.input, &res)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
