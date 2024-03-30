package http

import (
	"bytes"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSONStruct(t *testing.T) {
	input := []byte(`{"execution_optimistic":false,"finalized":true,"data":{"previous_version":"0x00000001","current_version":"0x00000002","epoch":"3"}}`)
	resType := phase0.Fork{}
	expectedData := phase0.Fork{
		PreviousVersion: phase0.Version{0x00, 0x00, 0x00, 0x01},
		CurrentVersion:  phase0.Version{0x00, 0x00, 0x00, 0x02},
		Epoch:           3,
	}
	expectedMetadata := map[string]any{
		"execution_optimistic": false,
		"finalized":            true,
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(input), resType)
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
	require.Equal(t, expectedMetadata, metadata)
}

func TestDecodeJSONArray(t *testing.T) {
	input := []byte(`{"execution_optimistic":false,"finalized":true,"data":[{"previous_version":"0x00000001","current_version":"0x00000002","epoch":"3"},{"previous_version":"0x00000002","current_version":"0x00000003","epoch":"4"}]}`)
	resType := []phase0.Fork{}
	expectedData := []phase0.Fork{
		{
			PreviousVersion: phase0.Version{0x00, 0x00, 0x00, 0x01},
			CurrentVersion:  phase0.Version{0x00, 0x00, 0x00, 0x02},
			Epoch:           3,
		},
		{
			PreviousVersion: phase0.Version{0x00, 0x00, 0x00, 0x02},
			CurrentVersion:  phase0.Version{0x00, 0x00, 0x00, 0x03},
			Epoch:           4,
		},
	}
	expectedMetadata := map[string]any{
		"execution_optimistic": false,
		"finalized":            true,
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(input), resType)
	require.NoError(t, err)
	require.Equal(t, expectedData, data)
	require.Equal(t, expectedMetadata, metadata)
}
