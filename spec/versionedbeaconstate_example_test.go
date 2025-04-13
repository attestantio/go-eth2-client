// Copyright © 2021 - 2024 Attestant Limited.
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

package spec_test

import (
	"context"
	"fmt"
	"os"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	eth2api "github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/rs/zerolog"
)

func ExampleVersionedBeaconState() {
	// Connect to a beacon node
	ctx := context.Background()
	client, err := http.New(ctx,
		http.WithAddress(os.Getenv("HOODI_CONSENSUS_ADDRESS")),
		http.WithLogLevel(zerolog.Disabled),
		http.WithTimeout(30*time.Second),
	)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Fetch beacon state
	beaconStateProvider, isProvider := client.(eth2client.BeaconStateProvider)
	if !isProvider {
		fmt.Println("Client does not support fetching beacon states")
		return
	}

	// Get state at slot 100
	fetchedState, err := beaconStateProvider.BeaconState(ctx, &eth2api.BeaconStateOpts{
		State: "100",
	})
	if err != nil {
		fmt.Printf("Failed to obtain beacon state: %v\n", err)
		return
	}

	state := fetchedState.Data

	// Generate proof for balances field
	proof, err := state.ProveField("Balances")
	if err != nil {
		fmt.Printf("Failed to generate proof: %v\n", err)
		return
	}

	// Get the field root (what we're proving)
	fieldRoot, err := state.FieldRoot("Balances")
	if err != nil {
		fmt.Printf("Failed to get field root: %v\n", err)
		return
	}
	fmt.Printf("Field root: %#x\n", fieldRoot)

	// Verify the proof
	valid, err := state.VerifyFieldProof(proof, "Balances")
	if err != nil {
		fmt.Printf("Failed to verify proof: %v\n", err)
		return
	}
	fmt.Printf("Proof verified: %v\n", valid)

	// Output:
	// Field root: 0xff5d10e166c63339da4ecbb013a77dd0d5d9dbf4a4e4115afd2fd11c5b57598d
	// Proof verified: true
}
