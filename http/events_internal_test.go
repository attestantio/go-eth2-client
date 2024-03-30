// Copyright © 2020, 2021 Attestant Limited.
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

package http

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// timeout for tests.
var timeout = 60 * time.Second

func TestEventHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handled := false
	handler := func(*api.Event) {
		handled = true
	}

	tests := []struct {
		name    string
		message *sse.Event
		handler client.EventHandlerFunc
		handled bool
	}{
		{
			name:    "MessageNil",
			handler: handler,
			handled: false,
		},
		{
			name:    "MessageEmpty",
			message: &sse.Event{},
			handler: handler,
			handled: false,
		},
		{
			name: "EventUnknown",
			message: &sse.Event{
				Event: []byte("unknown"),
			},
			handler: handler,
			handled: false,
		},
		{
			name: "HandlerNil",
			message: &sse.Event{
				Event: []byte("head"),
			},
			handled: false,
		},
		{
			name: "HeadGood",
			message: &sse.Event{
				Event: []byte("head"),
				Data:  []byte(`{"slot":"4095940","block":"0x73d83c5f925716c9bd2d1e9c339fb99b0ec4addef3e93f6f35d4c5f1de7ae092","state":"0xead0e6eb4004576546864f10cfa4aeac31afbf96abc405a86c00cbda8f3e8ed0","epoch_transition":false,"previous_duty_dependent_root":"0xeca94cc9180212a2cff2659289cc7e6f2df08a645120e35e25d09c2ddc7db5f1","current_duty_dependent_root":"0xdda286c4a096fc8ec0d6ba9e14e688cbb046bfb33462fdf94953e75d0cea0074","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "BlockGood",
			message: &sse.Event{
				Event: []byte("block"),
				Data:  []byte(`{"slot":"4095943","block":"0x1c3981b7439cd2dc53dca1a99122e1cacb36a13796d426d4c8a03ba745cb0c8b","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "AttestationGood",
			message: &sse.Event{
				Event: []byte("attestation"),
				Data:  []byte(`{"aggregation_bits":"0x00002840403040000000020008800040008042800000020220","data":{"slot":"4095945","index":"12","beacon_block_root":"0xff27c7551bf4cfe4dc4cce00920e7a5c5074860d1dbd8aa8b3b5f888523f51ff","source":{"epoch":"127997","root":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783"},"target":{"epoch":"127998","root":"0x46d4629861bd81cfc94007501b4edb1b3ca9444b41d7a98681b6c2f4bdb978bd"}},"signature":"0xacb9f562a28c4ef5b60b88678068ea51573a3237d3331dda3b2d845a0d03bc56ab2994d2deb90d9f074a8bdab59945150d0a7717e74b1bf2627f8971c81091f724c211dfce8fa16fb839c6a1bfd341ddec5e7eb88472682fd1a170e373660534"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "VoluntaryExitGood",
			message: &sse.Event{
				Event: []byte("voluntary_exit"),
				Data:  []byte(`{"message":{"epoch":"1", "validator_index":"1"}, "signature":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505cc411d61252fb6cb3fa0017b679f8bb2305b26a285fa2737f175668d0dff91cc1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "FinalizedCheckpointGood",
			message: &sse.Event{
				Event: []byte("finalized_checkpoint"),
				Data:  []byte(`{"block":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783","state":"0x9c237b2a66df8636f816e6b2c8860ba287fc5b817d882b1be8b7111486fb4ddc","epoch":"127997","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ChainReorgGood",
			message: &sse.Event{
				Event: []byte("chain_reorg"),
				Data:  []byte(`{"slot":"4100237","depth":"2","old_head_block":"0x5c988c12b7d8638c06e6b9511e09e5e28511a16a33c153413416d3fd5d95353a","new_head_block":"0x783ccb4310c0ddab3ea500f8d2b88c5ad8d6b2d601513f4ebf491066cda1d180","old_head_state":"0xdbad017808a1c5a77866fccc3e15f14a67585fff25b3a4d947ae9cc6d937b4ab","new_head_state":"0x67f0302bb939f64b1feaaa907b30c9631c0588480f2305578377bfd87df68b95","epoch":"128132","execution_optimistic":false}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ProposerSlashingGood",
			message: &sse.Event{
				Event: []byte("proposer_slashing"),
				Data:  []byte(`{"signed_header_1":{"message":{"slot":"6137846","proposer_index":"552061","parent_root":"0x0000000000000000000000000000000000000000000000000000000000000000","state_root":"0x0000000000000000000000000000000000000000000000000000000000000000","body_root":"0x82710f9c7025617ac3ef1f722451c05046ea01d525d60b7594669e63080e5db7"},"signature":"0xb544eaa34654372143d0442fbba6713d978f780da85503b2a070998970867b4936e884c9f948d51c6e08f6c927f0cd4d0a334e4f04404876e08863e73d8add1b3949fd911eb386fa24ea4c1dc062564b0af8a2f95e8940737f958c8baa584cd6"},"signed_header_2":{"message":{"slot":"6137846","proposer_index":"552061","parent_root":"0x873f73baa696664b8b73b160e7cfe352a924238935324007086e93a158b6e23d","state_root":"0xfad223f846ec0798c8a128e52ded70c02eeda3bdcf733db7ee77b0f02a46cef9","body_root":"0xe8b27e1cb99c74d0a39d19a573743426aba74cca872b3fc44115113c669d5a96"},"signature":"0x90f6920c78ed00a06affa06231780b9ad632b54ee2e1a6531e3881c9a96f47dadfba648566cdd684386df19150a9b6eb18a11e79f3d71e5fc85c7af81dd5fcd92f8ddbc870e526caed49c48ebc426489abc1efc64e3c7a14541b558877309141"}}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "AttesterSlashingGood",
			message: &sse.Event{
				Event: []byte("attester_slashing"),
				Data:  []byte(`{"attestation_1":{"attesting_indices":["1379","6498","7152","7611","10719","11656","15711","20612","26976","28038","30613","39996","72161","72684","75322","86121","89806","92747","94094","96198","96587","99048","100892","101838","117219","135916","139472","140755","145132","146549","148793","153618","154320","154414","156422","163002","165632","166740","167801","169764","171695","174283","176982","177143","177856","178016","180422","180593","186627","190359","192302","192973","201784","202000","213132","214455","216157","218189","218670","218872","221371","223604","225384","228340","232488","236360","237398","237769","239585","239870","243367","245718","247409","250257","251557","251850","255266","269326","269413","277447","278479","280748","282230","286188","286499","287581","293997","294106","297527","298179","299246","304274","305482","308474","312315","312856","315565","318733","319788","320418","322572","329735","331990","337653","337896","341979","344354","349106","350432","350598","350719","351231","353417","356015","356912","358208","358253","358775","360731","361675","361809","362528","372325","372331","372770","377376","386524","389323","389993","391013","394262","395176","395752","398129","399427","405839","407363","408517","409937","410538","411482","413713","416159","421753","426669","427714","428875","429693","430575","432645","437734","445366","448354","450266","452487","455531","455575","458233","468388","468877","471109","477185","480640","484346","484790","485906","487414","487705","493398","493527","497774","497887","498530","498791","501663","502258","503864","504418","506511","507566","508833","511154","515792","517955","522747","524759","533659","539011","539517","542154","542387","545084","547320","564444","565691","566062","566942","567905","569019","570238","574236","574836","583402","586842","589094","589373","590472","591564","592365","599309","601533","614547","615289","618674","620483","620874","621114","631082","631580","635052","635357","635632","640540","641046","641712","642095","644393","645990","650639","651381","655674","658553","658730","660381","663012","666071","671248","672021","673352","677063","680655","684323","685145","687842","688381","688464","692537","710555","713517","716137","716725","716893","720080","720715","721026","721397","721969","722568","726847","728006","728026","728968","729431","730915","731268","734600","734705","736076","742690","743783","744046","746246","746281","747514","755518","756113","759722","763739","764264","764403","774204","774657","775367","776295","776468","777628","777724","778906","778935","784104","790549","793434","793732","793991","798735","798740","799926","800022","805848","806025","808386","817029","818647","820641","824328","826419","827758","829139","830010","831317","837969","841062","843916","845313","845802","847624","852439","856121","862363","863997","868399","872373","872685","877312","877842","878805","878819","882746","884016","885135","888929","891416","896204","897101","898180","902030","902136","902843","903615","909161","920873","922759","923975","924185","924653","925291","927531","930388","931933","935314","936880","939910","940153","941510","942425","944415","945704","946209","949636","950399","953330","954070","954479","957162","958646","959098","959979","963938","963981","966927","968451","970800","971047","975171","975715","976123","978689","980754","980953","983241","983503","985284","988359","988747","989661","989900","990231","991724","994685","995191","995380","997021","998184","1002418","1005226","1006270","1009818","1014829","1014891","1015330","1015344","1015979","1016933","1019651","1019804","1022015","1022177","1022438","1022480","1024014","1026716","1026957","1028027","1028042","1031558","1032367","1032944","1033387","1034311","1034709","1037676","1039639","1043684","1044190"],"data":{"slot":"7858464","index":"6","beacon_block_root":"0x47e718a353454bf041b953cbad56f814546514cb0dd0987d79bad2636c8af36f","source":{"epoch":"245576","root":"0x209ea5723ad520c4b288ddefd00c18ef92cc802fad92a05479f2da13039d2bc9"},"target":{"epoch":"245577","root":"0x47e718a353454bf041b953cbad56f814546514cb0dd0987d79bad2636c8af36f"}},"signature":"0xb2073a97debb76ed591afb65d3a59eb32ac16f88d85e73ff0b5a19a935feb5bb8e4d3c9bdeca93bc070fc2ae29c7903308c1e6480232d03a2bb0d813c3c8d15dbf7a573df382f08051d364eda0d16d26bcbe3bff5c02892a58278d16d94d18ab"},"attestation_2":{"attesting_indices":["386524"],"data":{"slot":"7858464","index":"6","beacon_block_root":"0x030dc7282d3cca8416970028f33361a044445176119ad35bd35c2ed9e797f910","source":{"epoch":"245576","root":"0x209ea5723ad520c4b288ddefd00c18ef92cc802fad92a05479f2da13039d2bc9"},"target":{"epoch":"245577","root":"0x030dc7282d3cca8416970028f33361a044445176119ad35bd35c2ed9e797f910"}},"signature":"0x80a603ce9cd34572947273feb894e715f0497d843cea397f84586603888620cb9721f479d2f46f3eb6568c8e7bf1d19a075b993679bd60bdd04613693d49799309fb44216f8ad03eb1f327e543a1292a8573a338dca2d8cb6beb593b93cfff56"}}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "BLSToExecutionChangeGood",
			message: &sse.Event{
				Event: []byte("bls_to_execution_change"),
				Data:  []byte(`{"validator_index":"63401","from_bls_pubkey":"0xa46ed2574770ec1942d577ef89e0bf7b0d601349dab791740dead3fb5a6e2624cf62b9e58de1074c49f44b986eb39002","to_execution_address":"0xd641D2Cc74C7b6A641861260d07D67eB67bc7403"}`),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ContributionAndProofGood",
			message: &sse.Event{
				Event: []byte("contribution_and_proof"),
				Data:  []byte(`{"message":{"aggregator_index":"355177","contribution":{"slot":"4095970","beacon_block_root":"0xfe17cc29bb937740eead4b84716751b00e361891dae7c8f0e98f4deb5f753cbc","subcommittee_index":"0","aggregation_bits":"0xffffffffffffffffffffdfffffffffff","signature":"0x926e8fbf2b8599f76a42e2dd02b954853d3841577a0c68303fb9a5690f7973e95454bc1df03118d53c80e3cf13dc33490f2516aefb4cb7766a724a10dd0536811ec43b7e5f08442ef7dbc072b4b484ea1acde78e5aae8d636b06dd18677c535b"},"selection_proof":"0x84c805b21a40315dc19ea89e6d64d8f5e913d5e003a813d74a23f16add72250791a950b7508e2ea69adab8169c2800b70d26fea4cdb51f1fb96d939baaa44567eb9e96d2021799478fd3f557326a62060215be95465fcc9c49f67dd8685a20bc"},"signature":"0x8093efce898e36cab5ab2b198a48046d029b36909a29ec33ca7075f389133288c4d7e13cf3e20396612050d4aebe9212154fd5a2be4bf356e6191600d65906d5c404bd46c95ae20fe4bc5e18c6e2808c97a4572f995bf90db8aaf3fd84fb87ac"}`),
			},
			handler: handler,
			handled: true,
		},
	}

	s, err := New(ctx,
		WithTimeout(timeout),
		WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	h, isHTTPService := s.(*Service)
	require.True(t, isHTTPService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handled = false
			log := zerolog.New(&bytes.Buffer{})
			ctx = log.WithContext(ctx)
			h.handleEvent(ctx, test.message, test.handler)
			require.Equal(t, test.handled, handled)
		})
	}
}
