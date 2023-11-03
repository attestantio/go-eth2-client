package v1

import (
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// BlobSidecarEvent is the data for the blob sidecar event.
type BlobSidecarEvent struct {
	BlockRoot     phase0.Root
	Slot          phase0.Slot
	Index         deneb.BlobIndex
	KZGCommitment deneb.KZGCommitment
	VersionedHash deneb.VersionedHash
}

// blobSidecarEventJSON is the spec representation of the struct.
type blobSidecarEventJSON struct {
	BlockRoot     string `json:"block_root"`
	Slot          string `json:"slot"`
	Index         string `json:"index"`
	KZGCommitment string `json:"kzg_commitment"`
	VersionedHash string `json:"versioned_hash"`
}

// MarshalJSON implements json.Marshaler.
func (e *BlobSidecarEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobSidecarEventJSON{
		BlockRoot:     fmt.Sprintf("%#x", e.BlockRoot),
		Slot:          fmt.Sprintf("%d", e.Slot),
		Index:         fmt.Sprintf("%d", e.Index),
		KZGCommitment: fmt.Sprintf("%#x", e.KZGCommitment),
		VersionedHash: fmt.Sprintf("%#x", e.VersionedHash),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *BlobSidecarEvent) UnmarshalJSON(input []byte) error {
	var err error

	var blobSidecarEventJSON blobSidecarEventJSON
	if err = json.Unmarshal(input, &blobSidecarEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if blobSidecarEventJSON.BlockRoot == "" {
		return errors.New("block_root missing")
	}
	err = e.BlockRoot.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, blobSidecarEventJSON.BlockRoot)))
	if err != nil {
		return errors.Wrap(err, "invalid value for block_root")
	}
	if blobSidecarEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	err = e.Slot.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, blobSidecarEventJSON.Slot)))
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if blobSidecarEventJSON.Index == "" {
		return errors.New("index missing")
	}
	err = e.Index.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, blobSidecarEventJSON.Index)))
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	if blobSidecarEventJSON.KZGCommitment == "" {
		return errors.New("kzg_commitment missing")
	}
	err = e.KZGCommitment.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, blobSidecarEventJSON.KZGCommitment)))
	if err != nil {
		return errors.Wrap(err, "invalid value for kzg_commitment")
	}
	if blobSidecarEventJSON.VersionedHash == "" {
		return errors.New("versioned_hash missing")
	}
	err = e.VersionedHash.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, blobSidecarEventJSON.VersionedHash)))
	if err != nil {
		return errors.Wrap(err, "invalid value for versioned_hash")
	}

	return nil
}

// String returns a string version of the structure.
func (e *BlobSidecarEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
