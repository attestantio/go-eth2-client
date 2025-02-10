// Copyright © 2021 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package electra

//nolint:revive
// Need to `go install github.com/ferranbt/fastssz/sszgen@latest` for this to work.
//go:generate rm -f depositreceipts_ssz.go withdrawalrequests_ssz.go
//go:generate sszgen -suffix ssz -include ../../spec/phase0,../../spec/bellatrix,../../spec/electra -path . -objs DepositReceipts,ExecutionPayloadWithdrawalRequests
//go:generate goimports -w depositreceipts_ssz.go withdrawalrequests_ssz.go
