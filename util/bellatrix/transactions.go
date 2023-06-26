package bellatrix

import "github.com/attestantio/go-eth2-client/spec/bellatrix"

// ExecutionPayloadTransactions provides information about transactions.
type ExecutionPayloadTransactions struct {
	Transactions []bellatrix.Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}
