// Copyright © 2020 - 2026 Attestant Limited.
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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		err        string
		normalizes bool
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.eventJSON",
		},
		{
			name:  "TopicMissing",
			input: []byte(`{"data":{}}`),
			err:   "topic missing",
		},
		{
			name:  "TopicWrongType",
			input: []byte(`{"topic":[],"data":{}}`),
			err:   "invalid JSON: json: cannot unmarshal array into Go struct field eventJSON.topic of type string",
		},
		{
			name:  "TopicUnsupported",
			input: []byte(`{"topic":"foo","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
			err:   "unsupported event topic foo",
		},
		{
			name:  "DataMissing",
			input: []byte(`{"topic":"head"}`),
			err:   "data missing",
		},
		{
			name:       "GoodPhase0Attestation",
			normalizes: true,
			input:      []byte(`{"topic":"attestation","data":{"aggregation_bits":"0x010203","data":{"beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","index":"1","slot":"100","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
		},
		{
			name:       "GoodElectraAttestation",
			normalizes: true,
			input:      []byte(`{"topic":"attestation","data":{"aggregation_bits":"0xf77ffffffdfbfffffffdbfffffe5fff71f","data":{"slot":"98106","index":"0","beacon_block_root":"0xf8df02ed08b9adcb88a22cb22cd2a6074b184128ae6a240e3172109fdfacaa7b","source":{"epoch":"3064","root":"0x19ffd95e92753046cf63b4298f859e1fb1271a160ef0139ea1eb9f06d45d3b93"},"target":{"epoch":"3065","root":"0xffac9506e2262991ed19b1804bec9f7a1c4c4e61eb37444e6c8826bb362716d6"}},"signature":"0xb4f12c02e0f1a5db07999ceb8c1a4ccd41a3cb46ca15abe1c145337f1287360c49d5780fb7b44dfebeb96f3898824605008c9d458bdd2413358da3edf1b181d4e98edfe90d5fd016ac8f6aebc6646b2da83ab98722a7b4ee5264506bf6ae08e9","committee_bits":"0x0040000000000000"}}`),
		},
		{
			name:  "GoodSingleAttestation",
			input: []byte(`{"topic":"single_attestation","data":{"committee_index":"11","attester_index":"23784","data":{"slot":"98122","index":"0","beacon_block_root":"0x497033a5af8e64b748c554524e1e269da3c3af71515cf31f2d7bf9bab256a03c","source":{"epoch":"3065","root":"0xffac9506e2262991ed19b1804bec9f7a1c4c4e61eb37444e6c8826bb362716d6"},"target":{"epoch":"3066","root":"0x5982836668f92d786ef82f8841011a0888be5583bd09ab73760563855789d49b"}},"signature":"0xa448d3dc9520d4cb8e70094108169893a94ef7d074151ba333169ea92e2586da6f6efa622722743725c8012707f78efa02d99dd8ee094fa5bf5ca2b24066096ab0bc7671d4037521cbe69411871dd614e60d2c0eed8d2a0b2e4b77602b39d50e"}}`),
		},
		{
			name:       "GoodBlock",
			normalizes: true,
			input:      []byte(`{"topic":"block","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","slot":"1"}}`),
		},
		{
			name:  "GoodChainReorg",
			input: []byte(`{"topic":"chain_reorg","data":{"depth":"2","epoch":"16405","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","slot":"524986"}}`),
		},
		{
			name:  "GoodFinalizedCheckpoint",
			input: []byte(`{"topic":"finalized_checkpoint","data":{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","epoch":"2","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28"}}`),
		},
		{
			name:  "GoodHead",
			input: []byte(`{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
		},
		{
			name:  "GoodVoluntaryExit",
			input: []byte(`{"topic":"voluntary_exit","data":{"message":{"epoch":"1","validator_index":"2"},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
		},
		{
			name:  "GoodContributionAndProof",
			input: []byte(`{"topic":"contribution_and_proof","data":{"message":{"aggregator_index":"6568","contribution":{"aggregation_bits":"0x3f7f7f9fbffd9fddaf77fff7fffffdff","beacon_block_root":"0x3471a569ed74fb13f6638d7b759cd17c8ed08045d4668ae635349cc5f4dd2a75","signature":"0xa7260b90db427b85806cdaaecef08146a02c8c450aae96245be862f67fe6f54fefc7fdf1d4adfafad0164f5bbc0ceb65197b29e25b1dd9efd44ba8c390e95bd966b5dd97bf877a0ce277c757b68643054238659932348185775dc36d036b38da","slot":"45566","subcommittee_index":"2"},"selection_proof":"0x8c28b4b2f304f957735986e89ed3e429e007592e854d2c9a794333d5dfa05505412d70d0ba91e9fe3453816b01cd846415f82c864e7337e4796101ac9d2e351f2f8172d1d9061fd212f353ecf0ffd9dd17da42598adeae2046e5a74cbcb43474"},"signature":"0xb992ac86e1bbd6e2d1b7d18e8467aa435fecf583f5d13739db99b8d343093177caf167010b480c4e10f858f84cd05a1704c7ae3a253b1c454e5ebeefb7c35f7b8b51ba7aba0019cbc92d5bd8e9bcb61608a2ef47ce0a024b7b497ac9e813620f"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Event
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				if !test.normalizes {
					assert.JSONEq(t, string(test.input), string(rt))
				}
				assert.JSONEq(t, string(rt), res.String())
			}
		})
	}
}

func TestEventMarshalJSONUsesTypedData(t *testing.T) {
	input := []byte(`{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","epoch_transition":false,"slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`)

	var event api.Event
	require.NoError(t, json.Unmarshal(input, &event))
	event.Data.(*api.HeadEvent).Slot = 42

	output, err := json.Marshal(&event)
	require.NoError(t, err)
	require.JSONEq(t, `{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","epoch_transition":false,"slot":"42","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`, string(output))
}

func TestEventUnmarshalJSONReturnsTypedData(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected any
	}{
		{
			name:     "Head",
			input:    []byte(`{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
			expected: &api.HeadEvent{},
		},
		{
			name:     "HeadV2",
			input:    []byte(`{"topic":"head_v2","data":{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","payload_status":"empty","epoch_transition":false,"current_epoch_dependent_root":"0x5e0043f107cb57913498fbf2f99ff55e730bf1e151f02f221e977c91a90a0e91","next_epoch_dependent_root":"0x7f1154a218dc77addb1fe144266e83c04f67c8dfd03aacf955df05308cd8b27c","execution_optimistic":true}}}`),
			expected: &api.HeadEventV2{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event api.Event
			require.NoError(t, json.Unmarshal(test.input, &event))
			require.IsType(t, test.expected, event.Data)
		})
	}
}

func TestEventUnmarshalJSONVersionedGloasData(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		data     any
		expected any
	}{
		{
			name:     "ExecutionPayloadBid",
			topic:    "execution_payload_bid",
			data:     &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			expected: &gloas.SignedExecutionPayloadBid{},
		},
		{
			name:  "PayloadAttestationMessage",
			topic: "payload_attestation_message",
			data: &gloas.PayloadAttestationMessage{
				Data: &gloas.PayloadAttestationData{
					BeaconBlockRoot: phase0.Root{0x01},
					Slot:            1,
				},
			},
			expected: &gloas.PayloadAttestationMessage{},
		},
		{
			name:     "ProposerPreferences",
			topic:    "proposer_preferences",
			data:     &gloas.SignedProposerPreferences{Message: &gloas.ProposerPreferences{}},
			expected: &gloas.SignedProposerPreferences{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inner, err := json.Marshal(test.data)
			require.NoError(t, err)
			input, err := json.Marshal(map[string]any{
				"topic": test.topic,
				"data": map[string]any{
					"version": "gloas",
					"data":    json.RawMessage(inner),
				},
			})
			require.NoError(t, err)

			var event api.Event
			require.NoError(t, json.Unmarshal(input, &event))
			require.IsType(t, test.expected, event.Data)

			output, err := json.Marshal(&event)
			require.NoError(t, err)
			require.JSONEq(t, string(input), string(output))
		})
	}
}

func TestEverySupportedEventTopicHasTypedDispatch(t *testing.T) {
	tests := map[string]any{
		"attestation":                    &spec.VersionedAttestation{},
		"attester_slashing":              &phase0.AttesterSlashing{},
		"block":                          &api.BlockEvent{},
		"block_gossip":                   &api.BlockGossipEvent{},
		"bls_to_execution_change":        &capella.SignedBLSToExecutionChange{},
		"chain_reorg":                    &api.ChainReorgEvent{},
		"contribution_and_proof":         &altair.SignedContributionAndProof{},
		"data_column_sidecar":            &api.DataColumnSidecarEvent{},
		"execution_payload":              &api.ExecutionPayloadEvent{},
		"execution_payload_available":    &api.ExecutionPayloadAvailableEvent{},
		"execution_payload_bid":          &gloas.SignedExecutionPayloadBid{},
		"execution_payload_gossip":       &api.ExecutionPayloadEvent{},
		"fast_confirmation":              &api.FastConfirmationEvent{},
		"finalized_checkpoint":           &api.FinalizedCheckpointEvent{},
		"head":                           &api.HeadEvent{},
		"head_v2":                        &api.HeadEventV2{},
		"light_client_finality_update":   &api.LightClientFinalityUpdateEvent{},
		"light_client_optimistic_update": &api.LightClientOptimisticUpdateEvent{},
		"payload_attestation_message":    &gloas.PayloadAttestationMessage{},
		"payload_attributes":             &api.PayloadAttributesEvent{},
		"proposer_preferences":           &gloas.SignedProposerPreferences{},
		"proposer_slashing":              &phase0.ProposerSlashing{},
		"single_attestation":             &electra.SingleAttestation{},
		"voluntary_exit":                 &phase0.SignedVoluntaryExit{},
	}

	for topic, expected := range tests {
		t.Run(topic, func(t *testing.T) {
			input := []byte(`{"topic":"` + topic + `","data":{}}`)
			var event api.Event
			err := json.Unmarshal(input, &event)
			if err != nil {
				require.NotContains(t, err.Error(), "unsupported event topic")
			}
			require.IsType(t, expected, event.Data)
		})
	}
}

// TestSupportedEventTopics pins the exact topic enum from beacon-APIs commit
// ef98d512c03c8ca6b9d7cbdc45b9293ec2b24722.
func TestSupportedEventTopics(t *testing.T) {
	expected := map[string]bool{
		"attestation":                    true,
		"attester_slashing":              true,
		"block":                          true,
		"block_gossip":                   true,
		"bls_to_execution_change":        true,
		"chain_reorg":                    true,
		"contribution_and_proof":         true,
		"data_column_sidecar":            true,
		"execution_payload":              true,
		"execution_payload_available":    true,
		"execution_payload_bid":          true,
		"execution_payload_gossip":       true,
		"fast_confirmation":              true,
		"finalized_checkpoint":           true,
		"head":                           true,
		"head_v2":                        true,
		"light_client_finality_update":   true,
		"light_client_optimistic_update": true,
		"payload_attestation_message":    true,
		"payload_attributes":             true,
		"proposer_preferences":           true,
		"proposer_slashing":              true,
		"single_attestation":             true,
		"voluntary_exit":                 true,
	}

	require.Equal(t, expected, api.SupportedEventTopics)
}
