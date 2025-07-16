package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestDataColumnSidecarEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.dataColumnSidecarEventJSON",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "block_root missing",
		},
		{
			name:  "BlockRootWrongType",
			input: []byte(`{"block_root": true, "slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field dataColumnSidecarEventJSON.block_root of type string",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"block_root":"0xinvalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "invalid value for block_root: invalid value invalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "slot missing",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "index missing",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"-1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "invalid value for index: expected integer",
		},
		{
			name:  "KZGCommitmentsMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "kzg_commitments missing",
		},
		{
			name:  "KZGCommitmentsEmpty",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitments":[],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
			err:   "kzg_commitments missing",
		},
		{
			name:  "VersionedHashesMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"]}`),
			err:   "versioned_hashes missing",
		},
		{
			name:  "VersionedHashesEmpty",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":[]}`),
			err:   "versioned_hashes missing",
		},
		{
			name:  "Good",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f"]}`),
		},
		{
			name:  "RealExample",
			input: []byte(`{"block_root":"0xd565ec354fd256c54ff336b4d62a25312ffb0ef3004ea249cde2bf661a444ff9","slot":"127504","index":"60","kzg_commitments":["0xa590e760fdce951756d59c46b037bab8de815fe8ffc25e6e3a7b45e43289e1fdc942854cdfea1615385a0db63442f363","0x8479e2829495f00c47b1414332da60f961dba8767bbcfc8bd545b68df4966c3b4831ee1362eebe8d1d56e30b63bafb65"],"versioned_hashes":["0x01b27973b4b8de64263ccfb9068532eeff8d9e6280cd745ee754bbde6715ce2f","0x0122e2792a9b4ec6d42e16af0e028fb60dc09bffd0c48a7827db113c75eba0c2"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.DataColumnSidecarEvent
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
