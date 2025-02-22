// Code generated by fastssz. DO NOT EDIT.
// Hash: 5c0795a737413b7dee222139ce353bfc25323debce66933bff7b3193d76324e8
// Version: 0.1.3
package electra

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the ExecutionRequests object
func (e *ExecutionRequests) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(e)
}

// MarshalSSZTo ssz marshals the ExecutionRequests object to a target array
func (e *ExecutionRequests) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Offset (0) 'Deposits'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(e.Deposits) * 192

	// Offset (1) 'Withdrawals'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(e.Withdrawals) * 76

	// Offset (2) 'Consolidations'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(e.Consolidations) * 116

	// Field (0) 'Deposits'
	if size := len(e.Deposits); size > 8192 {
		err = ssz.ErrListTooBigFn("ExecutionRequests.Deposits", size, 8192)
		return
	}
	for ii := 0; ii < len(e.Deposits); ii++ {
		if dst, err = e.Deposits[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	// Field (1) 'Withdrawals'
	if size := len(e.Withdrawals); size > 16 {
		err = ssz.ErrListTooBigFn("ExecutionRequests.Withdrawals", size, 16)
		return
	}
	for ii := 0; ii < len(e.Withdrawals); ii++ {
		if dst, err = e.Withdrawals[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	// Field (2) 'Consolidations'
	if size := len(e.Consolidations); size > 2 {
		err = ssz.ErrListTooBigFn("ExecutionRequests.Consolidations", size, 2)
		return
	}
	for ii := 0; ii < len(e.Consolidations); ii++ {
		if dst, err = e.Consolidations[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the ExecutionRequests object
func (e *ExecutionRequests) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o0, o1, o2 uint64

	// Offset (0) 'Deposits'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 12 {
		return ssz.ErrInvalidVariableOffset
	}

	// Offset (1) 'Withdrawals'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return ssz.ErrOffset
	}

	// Offset (2) 'Consolidations'
	if o2 = ssz.ReadOffset(buf[8:12]); o2 > size || o1 > o2 {
		return ssz.ErrOffset
	}

	// Field (0) 'Deposits'
	{
		buf = tail[o0:o1]
		num, err := ssz.DivideInt2(len(buf), 192, 8192)
		if err != nil {
			return err
		}
		e.Deposits = make([]*DepositRequest, num)
		for ii := 0; ii < num; ii++ {
			if e.Deposits[ii] == nil {
				e.Deposits[ii] = new(DepositRequest)
			}
			if err = e.Deposits[ii].UnmarshalSSZ(buf[ii*192 : (ii+1)*192]); err != nil {
				return err
			}
		}
	}

	// Field (1) 'Withdrawals'
	{
		buf = tail[o1:o2]
		num, err := ssz.DivideInt2(len(buf), 76, 16)
		if err != nil {
			return err
		}
		e.Withdrawals = make([]*WithdrawalRequest, num)
		for ii := 0; ii < num; ii++ {
			if e.Withdrawals[ii] == nil {
				e.Withdrawals[ii] = new(WithdrawalRequest)
			}
			if err = e.Withdrawals[ii].UnmarshalSSZ(buf[ii*76 : (ii+1)*76]); err != nil {
				return err
			}
		}
	}

	// Field (2) 'Consolidations'
	{
		buf = tail[o2:]
		num, err := ssz.DivideInt2(len(buf), 116, 2)
		if err != nil {
			return err
		}
		e.Consolidations = make([]*ConsolidationRequest, num)
		for ii := 0; ii < num; ii++ {
			if e.Consolidations[ii] == nil {
				e.Consolidations[ii] = new(ConsolidationRequest)
			}
			if err = e.Consolidations[ii].UnmarshalSSZ(buf[ii*116 : (ii+1)*116]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ExecutionRequests object
func (e *ExecutionRequests) SizeSSZ() (size int) {
	size = 12

	// Field (0) 'Deposits'
	size += len(e.Deposits) * 192

	// Field (1) 'Withdrawals'
	size += len(e.Withdrawals) * 76

	// Field (2) 'Consolidations'
	size += len(e.Consolidations) * 116

	return
}

// HashTreeRoot ssz hashes the ExecutionRequests object
func (e *ExecutionRequests) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith ssz hashes the ExecutionRequests object with a hasher
func (e *ExecutionRequests) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Deposits'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Deposits))
		if num > 8192 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range e.Deposits {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 8192)
	}

	// Field (1) 'Withdrawals'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Withdrawals))
		if num > 16 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range e.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
	}

	// Field (2) 'Consolidations'
	{
		subIndx := hh.Index()
		num := uint64(len(e.Consolidations))
		if num > 2 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range e.Consolidations {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 2)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the ExecutionRequests object
func (e *ExecutionRequests) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(e)
}
