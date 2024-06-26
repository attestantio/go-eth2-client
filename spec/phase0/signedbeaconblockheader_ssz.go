// Code generated by fastssz. DO NOT EDIT.
// Hash: fa97a18404df62375826ed5ee7ded00832a9164171537b8327e9357589e54c7a
// Version: 0.1.3
package phase0

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the SignedBeaconBlockHeader object
func (s *SignedBeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the SignedBeaconBlockHeader object to a target array
func (s *SignedBeaconBlockHeader) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'Message'
	if s.Message == nil {
		s.Message = new(BeaconBlockHeader)
	}
	if dst, err = s.Message.MarshalSSZTo(dst); err != nil {
		return
	}

	// Field (1) 'Signature'
	dst = append(dst, s.Signature[:]...)

	return
}

// UnmarshalSSZ ssz unmarshals the SignedBeaconBlockHeader object
func (s *SignedBeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 208 {
		return ssz.ErrSize
	}

	// Field (0) 'Message'
	if s.Message == nil {
		s.Message = new(BeaconBlockHeader)
	}
	if err = s.Message.UnmarshalSSZ(buf[0:112]); err != nil {
		return err
	}

	// Field (1) 'Signature'
	copy(s.Signature[:], buf[112:208])

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SignedBeaconBlockHeader object
func (s *SignedBeaconBlockHeader) SizeSSZ() (size int) {
	size = 208
	return
}

// HashTreeRoot ssz hashes the SignedBeaconBlockHeader object
func (s *SignedBeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the SignedBeaconBlockHeader object with a hasher
func (s *SignedBeaconBlockHeader) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Message'
	if s.Message == nil {
		s.Message = new(BeaconBlockHeader)
	}
	if err = s.Message.HashTreeRootWith(hh); err != nil {
		return
	}

	// Field (1) 'Signature'
	hh.PutBytes(s.Signature[:])

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the SignedBeaconBlockHeader object
func (s *SignedBeaconBlockHeader) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}
