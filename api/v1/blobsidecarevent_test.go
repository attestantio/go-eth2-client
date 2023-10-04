package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBlobSidecarEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.blobSidecarEventJSON",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "block_root missing",
		},
		{
			name:  "BlockRootWrongType",
			input: []byte(`{"block_root": true, "slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobSidecarEventJSON.block_root of type string",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"block_root":"0xinvalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for block_root: invalid value invalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockRootShort",
			input: []byte(`{"block_root":"0x","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for block_root: incorrect length",
		},
		{
			name:  "BlockRootLong",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2a","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for block_root: incorrect length",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":true,"index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobSidecarEventJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"-1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for slot: invalid value -1: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":true,"kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobSidecarEventJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"-1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for index: invalid value -1: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "KZGCommitmentMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "kzg_commitment missing",
		},
		{
			name:  "KZGCommitmentWrongType",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":true,"versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobSidecarEventJSON.kzg_commitment of type string",
		},
		{
			name:  "KZGCommitmentInvalid",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0xinvalidfb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for kzg_commitment: invalid value invalidfb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "KZGCommitmentShort",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for kzg_commitment: incorrect length",
		},
		{
			name:  "KZGCommitmentLong",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2a","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for kzg_commitment: incorrect length",
		},
		{
			name:  "VersionedHashMissing",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2"}`),
			err:   "versioned_hash missing",
		},
		{
			name:  "VersionedHashWrongType",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobSidecarEventJSON.versioned_hash of type string",
		},
		{
			name:  "VersionedHashInvalid",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xinvalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
			err:   "invalid value for versioned_hash: invalid value invalide9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "VersionedHashShort",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0x"}`),
			err:   "invalid value for versioned_hash: incorrect length",
		},
		{
			name:  "VersionedHashLong",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2a"}`),
			err:   "invalid value for versioned_hash: incorrect length",
		},
		{
			name:  "Good",
			input: []byte(`{"block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","index":"1","kzg_commitment":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2","versioned_hash":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BlobSidecarEvent
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
