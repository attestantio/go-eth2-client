// Copyright © 2020 Attestant Limited.
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

package prysmgrpc

import (
	"context"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SubmitVoluntaryExit submits a voluntary exit.
func (s *Service) SubmitVoluntaryExit(ctx context.Context, voluntaryExit *spec.SignedVoluntaryExit) error {
	exit := &ethpb.SignedVoluntaryExit{
		Exit: &ethpb.VoluntaryExit{
			ValidatorIndex: uint64(voluntaryExit.Message.ValidatorIndex),
			Epoch:          uint64(voluntaryExit.Message.Epoch),
		},
		Signature: voluntaryExit.Signature[:],
	}

	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling ProposeExit()")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	_, err := conn.ProposeExit(opCtx, exit)
	cancel()

	if err != nil {
		return errors.Wrap(err, "failed to submit voluntary exit")
	}
	return nil
}
