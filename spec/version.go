// Copyright © 2026 Attestant Limited.
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

package spec

import (
	"github.com/attestantio/go-eth2-client/spec/version"
)

// The version types and constants live in the leaf spec/version package so that
// fork packages can depend on them without importing spec (which would create an
// import cycle). These aliases keep the historical spec.DataVersion* surface
// intact — additive, non-breaking per ADR-0001.
type (
	BuilderVersion = version.BuilderVersion
	DataVersion    = version.DataVersion
)

// Constant aliases preserve the original const semantics of these identifiers.
const (
	DataVersionUnknown   = version.DataVersionUnknown
	DataVersionPhase0    = version.DataVersionPhase0
	DataVersionAltair    = version.DataVersionAltair
	DataVersionBellatrix = version.DataVersionBellatrix
	DataVersionCapella   = version.DataVersionCapella
	DataVersionDeneb     = version.DataVersionDeneb
	DataVersionElectra   = version.DataVersionElectra
	DataVersionFulu      = version.DataVersionFulu
	DataVersionGloas     = version.DataVersionGloas

	BuilderVersionV1 = version.BuilderVersionV1
)

// DataVersionFromString turns a fork string into a DataVersion, returning an
// error if the fork is not recognised. Preserved here (the fork drops it from
// package spec) to keep the port additive per ADR-0001.
var DataVersionFromString = version.DataVersionFromString
