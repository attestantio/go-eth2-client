package capella

import "github.com/attestantio/go-eth2-client/spec/capella"

// ExecutionPayloadWithdrawals provides information about withdrawals.
type ExecutionPayloadWithdrawals struct {
	Withdrawals []*capella.Withdrawal `ssz-max:"8"`
}
