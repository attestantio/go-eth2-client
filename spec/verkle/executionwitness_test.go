// Copyright © 2023 Guillaume Ballet.
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

package verkle

import (
	"encoding/json"
	"testing"
)

func TestDeserializationIssueKaustinenBlock71208(t *testing.T) {
	var input = `
    {
      "parent_hash": "0xaea97ced7a027b6cce3ba97725a866efee13dc6554f2e2fd6ec7f367f9380500",
      "fee_recipient": "0xf97e180c050e5ab072211ad2c213eb5aee4df134",
      "state_root": "0x4cb4af71950a8acb5fe69796b09900f173cab95f3d04f21d377815c1e01ea968",
      "receipts_root": "0xb20fdb5d387b3cc990d93c96a75ee0c0f8d8bf58b68eb13036a31f48e1c28993",
      "logs_bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "prev_randao": "0x61a4e713c554be051b4e5c4653a9b23f028e85b54feac76551e4eddcfc52cbf7",
      "block_number": "71208",
      "gas_limit": "30000000",
      "gas_used": "110210",
      "timestamp": "1701681108",
      "extra_data": "0xd983010c01846765746889676f312e32302e3130856c696e7578",
      "base_fee_per_gas": "0x7",
      "block_hash": "0xdd2a9bacff5d2bb44326e86b7ffab95cf482f9e04f39e3af3207e58456d8b223",
      "transactions": [
        "0x02f89283010f2c048459682f008459682f088301ae82948d76c584db91ad2f6d362f5f642b599f382aa1c580a46057361d0000000000000000000000000000000000000000000000000000000000000010c080a0c13f00da261a2468aff6e401e437ca9d24ecf680d5ad793835b86edd66c74541a05f684dfbea710dba19af914e751f52bcbba1ff622cc4bf5f7606089d76a17416"
      ],
      "withdrawals": [],
      "execution_witness": {
        "stateDiff": [
          {
            "stem": "0x3e5157bc3f9283e0d138b00063ea4792af3f8519681637321479b4cd9d7a81",
            "suffixDiffs": [
              {
                "suffix": "0",
                "current_value": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "1",
                "current_value": "0x3a6f542b7b4c8037000000000000000000000000000000000000000000000000",
                "new_value": "0xaccb6ec520b67f37000000000000000000000000000000000000000000000000"
              },
              {
                "suffix": "2",
                "current_value": "0x0400000000000000000000000000000000000000000000000000000000000000",
                "new_value": "0x0500000000000000000000000000000000000000000000000000000000000000"
              },
              {
                "suffix": "3",
                "current_value": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
                "new_value": null
              },
              {
                "suffix": "4",
                "current_value": null,
                "new_value": null
              }
            ]
          },
          {
            "stem": "0x52a7f96b7d8d198407cd32493972a0b322c2fd04aeeec97557a63a98df0372",
            "suffixDiffs": [
              {
                "suffix": "0",
                "current_value": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "1",
                "current_value": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "2",
                "current_value": "0x0100000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "3",
                "current_value": "0x0bbcf059067c2cc2d5be230210a565d9922cf2cb5cc61230e50c5d6cb17cb45b",
                "new_value": null
              },
              {
                "suffix": "4",
                "current_value": "0xce01000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "64",
                "current_value": null,
                "new_value": "0x0000000000000000000000000000000000000000000000000000000000000010"
              },
              {
                "suffix": "128",
                "current_value": "0x00608060405234801561001057600080fd5b50600436106100365760003560e0",
                "new_value": null
              },
              {
                "suffix": "129",
                "current_value": "0x001c80632e64cec11461003b5780636057361d14610059575b600080fd5b6100",
                "new_value": null
              },
              {
                "suffix": "130",
                "current_value": "0x0143610075565b60405161005091906100f0565b60405180910390f35b610073",
                "new_value": null
              },
              {
                "suffix": "131",
                "current_value": "0x00600480360381019061006e919061013c565b610094565b005b6000600160ff",
                "new_value": null
              },
              {
                "suffix": "132",
                "current_value": "0x00610200811061008d5761008c610169565b5b0154905090565b806000819055",
                "new_value": null
              },
              {
                "suffix": "133",
                "current_value": "0x005080600160ff61020081106100b2576100b1610169565b5b01819055506001",
                "new_value": null
              },
              {
                "suffix": "134",
                "current_value": "0x008061010161020081106100cf576100ce610169565b5b018190555050565b60",
                "new_value": null
              },
              {
                "suffix": "135",
                "current_value": "0x0100819050919050565b6100ea816100d7565b82525050565b60006020820190",
                "new_value": null
              },
              {
                "suffix": "136",
                "current_value": "0x005061010560008301846100e1565b92915050565b600080fd5b610119816100",
                "new_value": null
              },
              {
                "suffix": "137",
                "current_value": "0x01d7565b811461012457600080fd5b50565b6000813590506101368161011056",
                "new_value": null
              },
              {
                "suffix": "138",
                "current_value": "0x005b92915050565b6000602082840312156101525761015161010b565b5b6000",
                "new_value": null
              },
              {
                "suffix": "139",
                "current_value": "0x0061016084828501610127565b91505092915050565b7f4e487b710000000000",
                "new_value": null
              }
            ]
          },
          {
            "stem": "0x5a1765b2eb16864c894b048584a9d1d131742897fd2d6f0bc9d9e558ae7626",
            "suffixDiffs": [
              {
                "suffix": "0",
                "current_value": null,
                "new_value": "0x0000000000000000000000000000000000000000000000000000000000000010"
              },
              {
                "suffix": "2",
                "current_value": null,
                "new_value": "0x0000000000000000000000000000000000000000000000000000000000000001"
              }
            ]
          },
          {
            "stem": "0x8dc286880de0cc507d96583b7c4c2b2b25239e58f8e67509b32edb5bbf293c",
            "suffixDiffs": [
              {
                "suffix": "0",
                "current_value": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "1",
                "current_value": "0xf8613f9ae1822a1d000000000000000000000000000000000000000000000000",
                "new_value": "0xf83f19003c192b1d000000000000000000000000000000000000000000000000"
              },
              {
                "suffix": "2",
                "current_value": "0x0000000000000000000000000000000000000000000000000000000000000000",
                "new_value": null
              },
              {
                "suffix": "3",
                "current_value": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
                "new_value": null
              },
              {
                "suffix": "4",
                "current_value": null,
                "new_value": null
              }
            ]
          }
        ],
        "verkle_proof": {
          "otherStems": [],
          "depthExtensionPresent": "0x12121012",
          "commitmentsByPath": [
            "0x4894cc6695ee89f10c2da09cca27bfae66c9b183eb336f331d9ec62262fc8165",
            "0x36e89e52f57e1a548f66ba5889068ab565dc5ed9aac0f5933e635c759712f796",
            "0x03203502dff858a4a07d7c93f09ad4ec8671aa88a77ae364d570451ad6253b67",
            "0x0a31c4c5ddc542277d88ddd8007687976d26bd0bd5f57bdf366f5aba90488c1b",
            "0x40a8a0d0e6b72b5ea10e325d624057e19f69fee42f659e062b44f9200fc81b11",
            "0x6ec14ddf5c6a589c9ecff762e844f8d04fae2573576a42018bad865c8db5c2cf",
            "0x6154b84ffca2e07c1869cc831a650d520d0b16340b9f0de642d411200b3ee9e3",
            "0x2519804b999abf7113ff747d3be2acfe19f7039c15f6bd7f3915cf8f3937ba1d",
            "0x4abf483e780fb88ad1888d1e831925137980f13c004b64353959e8f744f62e98",
            "0x6ef1777ec06ada68e5218d7166ecb483813f7820086ff6b311b7f9899b640eee",
            "0x192b33b443671b0d35b76aa6f8774a1bec31eddf4ec409a206b94aca207c8bac"
          ],
          "d": "0x069a6435f89df41175ced2c1f6af7638dd3fb0fdce4e0923ae626baaf0d7ca6a",
          "ipaProof": {
            "cl": [
              "0x4468a46c4a53af39d036d4b969c9064e21b6ef6300d1e8b38bb7ac9687410ded",
              "0x28b168c654e233bcea68dccda1c3dd593f88d6081ee2e71a9eb277b0627fa1d3",
              "0x22e76996bc8d388ab8a70a4f1292eb014a402395b9f2ae1de6448e5a56e5300d",
              "0x6eafcb57cb9be5682c46ebc6800e904ac7888a65ea538fe76fa35a7181b41932",
              "0x0596c5e3bdac6525588cb91fd326d45bc79908f574991771c0c279a82d7d64e0",
              "0x3b0d64d2f7589f2656af9b23608ab4cbca09d6d0c3479e8599197950ceec1d45",
              "0x6fbc04881ec27a917ae2aca6220dd098447bb8ffcd2411733de171340a485730",
              "0x2579876cf47c94953cc4a9e0e2c7369d715603ca792e384db36624e04e541bfe"
            ],
            "cr": [
              "0x4abb9f6a2f78749d87b3d09c048cb993820d6f1c9bc222b179279067e4b39bf6",
              "0x29c4c3fa1525a1ae0c8d031465d045d4e7b995a586e71f876dab13fbad0d5e8b",
              "0x1ca3a5864bc2a1061679853302ab3ffd2b7858486a98b03376bcc3635b63b45d",
              "0x5d541092555e68e59377e29a5869b0378e876a32de4aee513de5848eff231953",
              "0x641d6805e2b9b7687b6b53c085882f13922e083e0df1580d420ceb78a03599fd",
              "0x3022a96626a833f37300d6346ff0562bc6126bf82103009961dbf5834cda8fed",
              "0x1b24d8e7c7e4182df84ec2bddade28ecbfb5bc63b72e073db78c8123e6cff5aa",
              "0x6b088795de6eba8a2931984105e7efd45299a86c0f92c8f2d6040e477936b56a"
            ],
            "finalEvaluation": "0x1617ec35c4560a790a290ef8a1c26a764c6ae474f737069e8db51647976ab0c2"
          }
        }
      }
    }
		`

	var ep ExecutionPayload
	if err := json.Unmarshal([]byte(input), &ep); err != nil {
		t.Fatal(err)
	}

	if ep.ExecutionWitness == nil {
		t.Fatal("nil execution witness")
	}

	if len(ep.ExecutionWitness.StateDiff) == 0 {
		t.Fatal("zero-length state diff")
	}

	if len(ep.ExecutionWitness.StateDiff[1].SuffixDiffs) != 18 {
		t.Fatalf("invalid length for suffix diffs, want 18 got %d", len(ep.ExecutionWitness.StateDiff[1].SuffixDiffs))
	}
}
