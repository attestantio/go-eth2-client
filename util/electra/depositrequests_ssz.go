// Code generated by fastssz. DO NOT EDIT.
// Hash: bd021322a0585214e3e2b5d17c3251ad136909024b45304057c72b6c80c84729
// Version: 0.1.3
package electra

import (
	"github.com/attestantio/go-eth2-client/spec/electra"
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the DepositRequests object
func (d *DepositRequests) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(d)
}

// MarshalSSZTo ssz marshals the DepositRequests object to a target array
func (d *DepositRequests) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(4)

	// Offset (0) 'DepositRequests'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(d.DepositRequests) * 192

	// Field (0) 'DepositRequests'
	if size := len(d.DepositRequests); size > 8192 {
		err = ssz.ErrListTooBigFn("DepositRequests.DepositRequests", size, 8192)
		return
	}
	for ii := 0; ii < len(d.DepositRequests); ii++ {
		if dst, err = d.DepositRequests[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the DepositRequests object
func (d *DepositRequests) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 4 {
		return ssz.ErrSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'DepositRequests'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 4 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (0) 'DepositRequests'
	{
		buf = tail[o0:]
		num, err := ssz.DivideInt2(len(buf), 192, 8192)
		if err != nil {
			return err
		}
		d.DepositRequests = make([]*electra.DepositRequest, num)
		for ii := 0; ii < num; ii++ {
			if d.DepositRequests[ii] == nil {
				d.DepositRequests[ii] = new(electra.DepositRequest)
			}
			if err = d.DepositRequests[ii].UnmarshalSSZ(buf[ii*192 : (ii+1)*192]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the DepositRequests object
func (d *DepositRequests) SizeSSZ() (size int) {
	size = 4

	// Field (0) 'DepositRequests'
	size += len(d.DepositRequests) * 192

	return
}

// HashTreeRoot ssz hashes the DepositRequests object
func (d *DepositRequests) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(d)
}

// HashTreeRootWith ssz hashes the DepositRequests object with a hasher
func (d *DepositRequests) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'DepositRequests'
	{
		subIndx := hh.Index()
		num := uint64(len(d.DepositRequests))
		if num > 8192 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range d.DepositRequests {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 8192)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the DepositRequests object
func (d *DepositRequests) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(d)
}