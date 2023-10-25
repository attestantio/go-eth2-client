// Copyright © 2023 Attestant Limited.
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

package api

import (
	"errors"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedProposal contains a versioned signed beacon node proposal.
type VersionedSignedProposal struct {
	Version   spec.DataVersion
	Phase0    *phase0.SignedBeaconBlock
	Altair    *altair.SignedBeaconBlock
	Bellatrix *bellatrix.SignedBeaconBlock
	Capella   *capella.SignedBeaconBlock
	Deneb     *apiv1deneb.SignedBlockContents
}

// Slot returns the slot of the signed proposal.
func Slot(p VersionedSignedProposal) (phase0.Slot, error) {
	switch p.Version {
	case spec.DataVersionAltair:
		if p.Altair == nil || p.Altair.Message == nil {
			return 0, errors.New("no altair block")
		}

		return p.Altair.Message.Slot, nil
	case spec.DataVersionBellatrix:
		if p.Bellatrix == nil || p.Bellatrix.Message == nil {
			return 0, errors.New("no bellatrix block")
		}

		return p.Bellatrix.Message.Slot, nil
	case spec.DataVersionCapella:
		if p.Capella == nil || p.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}

		return p.Capella.Message.Slot, nil
	case spec.DataVersionDeneb:
		if p.Deneb == nil || p.Deneb.SignedBlock == nil || p.Deneb.SignedBlock.Message == nil {
			return 0, errors.New("no deneb block")
		}

		return p.Deneb.SignedBlock.Message.Slot, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedProposal) String() string {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	default:
		return "unsupported version"
	}
}
