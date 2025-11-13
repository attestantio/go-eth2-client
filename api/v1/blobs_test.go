package v1_test

import (
	"os"
	"testing"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
)

func TestBlobs(t *testing.T) {
	/* blob.ssz obtained from:
	http://<hoodi rpc url>/eth/v1/beacon/blobs/1713184 -H "Accept: application/octet-stream"
	Fulu hardfork
	*/
	var blobSSZRaw []byte
	blobSSZRaw, err := os.ReadFile("testdata/blobs/blob.ssz")
	if err != nil {
		t.Fatalf("Failed to read blob.ssz: %v", err)
	}

	var blobs v1.Blobs
	err = blobs.UnmarshalSSZ(blobSSZRaw)
	if err != nil {
		t.Fatalf("Failed to unmarshal blobs: %v", err)
	}

	// Test Blobs methods
	t.Run("String", func(t *testing.T) {
		t.Logf("Blobs: %s", blobs.String())
	})
	t.Run("MarshalSSZ", func(t *testing.T) {
		blobSSZ, err := blobs.MarshalSSZ()
		if err != nil {
			t.Fatalf("Failed to marshal blobs: %v", err)
		}
		t.Logf("Blobs: %s", blobSSZ)
	})
	t.Run("SizeSSZ", func(t *testing.T) {
		size := blobs.SizeSSZ()
		t.Logf("Size: %d", size)
	})
	t.Run("HashTreeRoot", func(t *testing.T) {
		hash, err := blobs.HashTreeRoot()
		if err != nil {
			t.Fatalf("Failed to get hash tree root: %v", err)
		}
		t.Logf("Hash: %s", hash)
	})
	t.Run("GetTree", func(t *testing.T) {
		_, err := blobs.GetTree()
		if err != nil {
			t.Fatalf("Failed to get tree: %v", err)
		}
	})
}
