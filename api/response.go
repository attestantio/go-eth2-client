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

// TODO remove.
// // ResultCode provides info about the result of a request.
// type ResultCode int
//
// // TODO do we need ResultCode?  Can we just return an error instead?
// // Errors are HTTP-specific, do we really care though?
// const (
// 	// ResultCodeUnknown is an unknown result code.
// 	ResultCodeUnknown ResultCode = iota
//
// 	// ResultCodeBadRequest states there was an issue with the options passed
// 	// in the request, meaning that the server could not enact the request.
// 	ResultCodeBadRequest
//
// 	// ResultCodeSyncing states the request was made but the server is currently
// 	// syncing with other beacon nodes, and so is unable to return the requested
// 	// information.
// 	ResultCodeSyncing
//
// 	// ResultCodeServerFailure states the request was made but the server failed to
// 	// handle the request properly due to some internal error.
// 	ResultCodeServerFailure
//
// 	// ResultCodeOK states a successful operation.  This does not
// 	// mean that Data is populated, for example a successful sending
// 	// of information may not return anything from the server.
// 	ResultCodeOK
//
// 	// ResultCodeNotFound states the request was made but no data was returned
// 	// by the server.  This could be due to the server not containing the data,
// 	// for example querying for an old state from a non-archive node.
// 	ResultCodeNotFound
// )

// Response is a response from the beacon API which may contain metadata.
// TODO should this be RequestResponse or similar?  Or does this work for submissions as well?
type Response[T any] struct {
	Data     T
	Metadata map[string]any
}
