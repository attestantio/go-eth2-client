// Code generated by fastssz. DO NOT EDIT.
// Hash: 17d4c9180818d70e873edf284079b326d586a16686d17c7c974a8a2fd19ec3e9
// Version: 0.1.3
package electra

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the ExecutionLayerWithdrawalRequest object
func (e *ExecutionLayerWithdrawalRequest) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(e)
}

// MarshalSSZTo ssz marshals the ExecutionLayerWithdrawalRequest object to a target array
func (e *ExecutionLayerWithdrawalRequest) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'SourceAddress'
	dst = append(dst, e.SourceAddress[:]...)

	// Field (1) 'ValidatorPubkey'
	dst = append(dst, e.ValidatorPubkey[:]...)

	// Field (2) 'Amount'
	dst = ssz.MarshalUint64(dst, uint64(e.Amount))

	return
}

// UnmarshalSSZ ssz unmarshals the ExecutionLayerWithdrawalRequest object
func (e *ExecutionLayerWithdrawalRequest) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 76 {
		return ssz.ErrSize
	}

	// Field (0) 'SourceAddress'
	copy(e.SourceAddress[:], buf[0:20])

	// Field (1) 'ValidatorPubkey'
	copy(e.ValidatorPubkey[:], buf[20:68])

	// Field (2) 'Amount'
	e.Amount = phase0.Gwei(ssz.UnmarshallUint64(buf[68:76]))

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ExecutionLayerWithdrawalRequest object
func (e *ExecutionLayerWithdrawalRequest) SizeSSZ() (size int) {
	size = 76
	return
}

// HashTreeRoot ssz hashes the ExecutionLayerWithdrawalRequest object
func (e *ExecutionLayerWithdrawalRequest) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith ssz hashes the ExecutionLayerWithdrawalRequest object with a hasher
func (e *ExecutionLayerWithdrawalRequest) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'SourceAddress'
	hh.PutBytes(e.SourceAddress[:])

	// Field (1) 'ValidatorPubkey'
	hh.PutBytes(e.ValidatorPubkey[:])

	// Field (2) 'Amount'
	hh.PutUint64(uint64(e.Amount))

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the ExecutionLayerWithdrawalRequest object
func (e *ExecutionLayerWithdrawalRequest) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(e)
}