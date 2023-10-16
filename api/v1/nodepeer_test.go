package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestNodePeersJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "GoodNoENR",
			input: []byte(`{"data":[{"peer_id":"16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/43402","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmDnjN82ckWtkRR2MGWC8vHryZ2BkYavB9ddnjrGDSvPyP","last_seen_p2p_address":"/ip4/10.0.20.20/tcp/45604","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmS7YMyhDwdBYGpiZXjYgqDD6EPtpArpUC4syqK5f8D6Dx","last_seen_p2p_address":"/ip4/10.0.20.17/tcp/60144","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm2q9YfmAJCn8GThcUmCCpra4a595oiWqbVEHZ6Qim1H2u","last_seen_p2p_address":"/ip4/10.0.20.13/tcp/37096","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmEqiEtEaourJ78Wkmbir1QbeeHy7FGTAWVMsDrMWdtYCt","last_seen_p2p_address":"/ip4/10.0.20.23/tcp/9000","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAm1fB6uumkVkWrgJckQm1vC1SNxMvus2FeCbzrD1HsfnCv","last_seen_p2p_address":"/ip4/10.0.20.18/tcp/38182","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmFj8Y4QZj3EUEFW9Hhf81nERESkdhU2s7tmodQnb2mDRS","last_seen_p2p_address":"/ip4/10.0.20.22/tcp/42234","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmKpMt5tXuao9nCe3NjhhPn3ny6VY9SD4neMrCyPgjD5Qt","last_seen_p2p_address":"/ip4/10.0.20.21/tcp/59882","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm1V3SzAg6CrC3vcyWkGFQpTBX6pywws12pNzXKmiqpEtv","last_seen_p2p_address":"/ip4/10.0.20.9/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm48vBTmKSaYtgSdHF3Y38b9Ajsa9xEJ4WFWAn65YymqLb","last_seen_p2p_address":"/ip4/10.0.20.4/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm5SdiyWCnR4pN5zUZYan5tfPBbLLCkDYsvJNLwW7fnrFj","last_seen_p2p_address":"/ip4/10.0.20.6/tcp/40380","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm3uj2dJ3WNQ8bVpPersrzpoetpKuXH3smTfW6EcDn3WFv","last_seen_p2p_address":"/ip4/10.0.20.14/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmHAZ9XGJiciW3ADQZRd1r88M7djpk16jWp8YJdEnEXZmZ","last_seen_p2p_address":"/ip4/10.0.20.19/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAkxwmpMqLLgJFgKi3RsHsPUNznggk5KQuWaTx3C6L4wuxe","last_seen_p2p_address":"/ip4/10.0.20.12/tcp/34260","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm6G56FAy33hzjZ2Qb8LnAaRpuW4Abt1ZvfSfp5aQK6ZGw","last_seen_p2p_address":"/ip4/10.0.20.11/tcp/53066","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmRzt5nPmxqU4YTNdWy5XEPgiFTyz1mTBPojEj9aefNS9F","last_seen_p2p_address":"/ip4/10.0.20.7/tcp/42976","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmQRZp5sxhEHjBYhf7GHKm67DnUyNcC1ApMNpkx8bxGxnY","last_seen_p2p_address":"/ip4/10.0.20.16/tcp/53284","state":"connected","direction":"inbound"}],"meta":{"count":17}}`),
		},
		{
			name:  "GoodWithENR",
			input: []byte(`{"data":[{"peer_id":"16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","enr":"enr:-MS4QExfvXqHhj-nqAqkg1Sn55uV7YgpRtlImGCvMJkrkbnLDo8sGhecAGid9B3NjXzN3UtGxpOOUqHZVcEDQxkniwoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBWEcXVpY4IjKYlzZWNwMjU2azGhA9mr4bIskWVeMt0dEn4IlQJhOFgOqgR9V3gkHTl1lTioiHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo","last_seen_p2p_address":"/ip4/10.0.20.21/udp/9001/quic-v1/p2p/16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmMtsF6vRyxhecq7tWFFkaQb9JTAXt3qADuYQybmEC1Azc","enr":"enr:-MK4QIpk8gn6FwjQWkDVvb1yxWJ_zI8zil7MVnoYoXbt_QOAanfBSeYBuK76T-EO_BjkM9ZQ5ap4hl_m4REytf-QKYqGAYs1FFUIh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBOJc2VjcDI1NmsxoQOJR0czLiBiqSYHGvyI9eYR1XRjA14trUpWxRmo3LluIYhzeW5jbmV0cwCDdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.19/tcp/9000/p2p/16Uiu2HAmMtsF6vRyxhecq7tWFFkaQb9JTAXt3qADuYQybmEC1Azc","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmKsbtgFhPUAxqgu3LBY9FCGcbYykWAtUheLmSj74ksEHM","last_seen_p2p_address":"/ip4/10.0.20.7/tcp/53854","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm8Qate2mJRqi34Znd5ZMD16KLMKYQtCGSoBTtPoy76Tfv","enr":"enr:-MK4QOMybXQ4D69oOCmSGk8w1OeG20w2gt3zs25Y5mvtObROI6cJvUTudnd5PiZ1-PgeGCmDJlVen_l020vjMj62-8mGAYs1FFRQh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFA6Jc2VjcDI1NmsxoQLA4JhhQgoWMvCzYl2puFenymWZ9dlcvbSksibMCK_QIYhzeW5jbmV0cwCDdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.14/tcp/9000/p2p/16Uiu2HAm8Qate2mJRqi34Znd5ZMD16KLMKYQtCGSoBTtPoy76Tfv","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAm2GC2bBhQkfzy1qoh513EjcWpkAKV3sfEj1YAyJGbPAj3","last_seen_p2p_address":"/ip4/10.0.20.22/tcp/45320","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm6ZLwRQee1qKp42zs7FugRSTQb18krF9NADCEBZ97yN6G","enr":"enr:-Ly4QK-p_H7Bvh8GwtZNB2Bj3v5YXWO9bmoSYei6PYRCDe9uOZ6MKzQLcSFNE_H7hqDHFTG2vmHpdyVAPmCCwRBQvOkFh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFAqJc2VjcDI1NmsxoQKlZ6xNTKcz9e_K3bP_SDulzTuc5eZD2XEpD7322gjNRYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.10/tcp/9000/p2p/16Uiu2HAm6ZLwRQee1qKp42zs7FugRSTQb18krF9NADCEBZ97yN6G","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAkzUpRJd1DQ8LC3TpwZ65Jjerw11i691ANCBbZr2Unakpr","enr":"enr:-MS4QNnQO0_uWygTYdQ8iAKd2G67S9909WbER6TgaHu74JQ9UtZ2gBcd0bjtQGOldhofvYBLJZhbERjW8vjEU-7Tt0cBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBCEcXVpY4IjKYlzZWNwMjU2azGhAksYuFN-ZB6meUe_B5NxQr8wvWhKILq4ia6F_RXyQb97iHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo","last_seen_p2p_address":"/ip4/10.0.20.16/udp/9001/quic-v1/p2p/16Uiu2HAkzUpRJd1DQ8LC3TpwZ65Jjerw11i691ANCBbZr2Unakpr","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmTzrtoo9EsY7dyAjwrqYT8BZk9p2gaHL9EohSgrvdjjRJ","enr":"enr:-Ly4QPJq9UA6fmqlQ3E_XgRqipZ8HrMkfzY2ZXYWmKj8g6qcCzOFrRG0Z4Cb-B5uf8t9qtzq-WbwItcXCUQwRr07Ri4Dh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFA2Jc2VjcDI1NmsxoQPj9n0nczI4wI6gQ0sABxt_ecH-r0W464uDe9qDNsEdGYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.13/tcp/9000/p2p/16Uiu2HAmTzrtoo9EsY7dyAjwrqYT8BZk9p2gaHL9EohSgrvdjjRJ","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAkx1ytpjGK3aswEY6mvCoyo4CKoomoaQ71nKMmveBCPMXy","last_seen_p2p_address":"/ip4/10.0.20.12/tcp/57628","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm4HsRm9Mp5nPdqmovPEMpg9dwRKoPPmNCZozv4n5ikmaC","last_seen_p2p_address":"/ip4/10.0.20.5/tcp/56860","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAkzGC4VazBPwjGBqP2eHVr3peLYafFTgr2RFhyq4LtyZrQ","enr":"enr:-Ly4QLZAC1fu95hzwUH4SIEz8OouR7qihgdFVCPD4PAt_zsBSFTf8ZE3pGW-bnrkIt0bcWaKWOlw3fxfDbmRdZLJvLwDh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBKJc2VjcDI1NmsxoQJH3KOnvaAGesYhGYEfTvqIG4QcEu1Wf4LqOJB_sJy4QYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.18/tcp/9000/p2p/16Uiu2HAkzGC4VazBPwjGBqP2eHVr3peLYafFTgr2RFhyq4LtyZrQ","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmTgbJNFvM4SHkT9Y8y15otUqV2h2SqLV8NXvp51Bs8nrf","enr":"enr:-Ly4QAc_eplEQJmA_L3MbcCECk9brHSL5SRqFK_XTwch1f7nPPGCKQmhkynqRphWQBtcKVaxGDwvsM_9nvevVPYwCiADh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFAiJc2VjcDI1NmsxoQPfSGcAx0ljjsPP8JVCumTFgd4n_K2a1HugQrVZZrTCDIhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/9000/p2p/16Uiu2HAmTgbJNFvM4SHkT9Y8y15otUqV2h2SqLV8NXvp51Bs8nrf","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmBKVWsod7jzQt5hLqBeRdYdBVcSHwgyAMrQXau3dxqa7u","enr":"enr:-MS4QIL4bhSRzO3Rr8EEPnMGonWOC2hfIhwghfiIn28j1QDrGM_BDXf-G5dhq6mOkrCPKwvxeMjdneNdgmHhEAwdDT4Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFAuEcXVpY4IjKYlzZWNwMjU2azGhAuwlrJfSnuhidlqna3R6kDrujloBS8EVPTlV02P8qRMkiHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo","last_seen_p2p_address":"/ip4/10.0.20.11/udp/9001/quic-v1/p2p/16Uiu2HAmBKVWsod7jzQt5hLqBeRdYdBVcSHwgyAMrQXau3dxqa7u","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmQY3XSZh3dimduUUw5yFALLGGpGc2qk4M6KaJgrzB9iK4","enr":"enr:-Ly4QPIUYqLFgF59q5HMixlv_5agvE8fcJDC-KlY2vUO-cTlffKMS4lmlCLejoiR5stS-klJaW6RdzCQJFcKLvUdPMoFh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFA-Jc2VjcDI1NmsxoQOwhMR45ADSbfkpYZLV10gmksm65K8NCPtnYqCHVdaKG4hzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.15/tcp/9000/p2p/16Uiu2HAmQY3XSZh3dimduUUw5yFALLGGpGc2qk4M6KaJgrzB9iK4","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAm9ZvLT2KkuujyP7tXcGjjrU9BJLHJxdECPR65Cmyoermy","last_seen_p2p_address":"/ip4/10.0.20.17/tcp/34150","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAmHp1riEgXGaasrngCyxtZAQBtuF927Vx9tpQyXR2hEUvZ","enr":"enr:-MK4QKfnKqJKdBYGRO7CQZj8Bj81ZYbvAcb6_GWLHKBgWTKFPKzmh7r6SJrUa33WxsY8HMThy1b8mXhzD_VLDfDEw3aGAYs1FFcgh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFAmJc2VjcDI1NmsxoQNMml2-4Lortplm2UKHHnSOMcewggOLQDzzS6rGg1UWFohzeW5jbmV0cwCDdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.9/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm4hYYzWpsfPbDKwbAfhrfN43N8rhfwVki1G41hH6cviV1","enr":"enr:-Ly4QDaHedUmiK-mgo72-e5WgoUK5BN6j7E_H9QMxAZ0fRoxXKa5Me7Fa75jV_ftWOyjYQFCY76n4j4uvlndpPTbcqIDh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBeJc2VjcDI1NmsxoQKJyhOtSXA8P_9VCUEgCctdDIc-N7VgbyqBhVnUkWwitIhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.23/tcp/9000/p2p/16Uiu2HAm4hYYzWpsfPbDKwbAfhrfN43N8rhfwVki1G41hH6cviV1","state":"connected","direction":"outbound"},{"peer_id":"16Uiu2HAmMXHcpVD6YNQ87zg5V5PDes5ngh2E6GVYcnVdA1KbUHnU","enr":"enr:-MK4QJJ6D9zx6pJt7rfrrKwbNurEfssaNNgrXM3JPIjIVbA3OrGyKUtXKyUEyrfRQhQ-P97SnBwnNU7_KXbfADOymlaGAYs1FFdEh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFASJc2VjcDI1NmsxoQODwA-B2P7gH4K1l1N5apaJJNybawXEhrrY-Dt_RDv2pYhzeW5jbmV0cwCDdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.4/tcp/9000","state":"connected","direction":"inbound"},{"peer_id":"16Uiu2HAm7HXE58UfHDcv4NNYXMRi53xy7rDVNVSCMKBWVRQ1R8RJ","enr":"enr:-Ly4QFXp16r9rKDoe5k08Xe4G_cbkdNuJ3cur1J9A7nvHJCXXQDmCOY5JyrzLJImci4XfAEUF-tUoqCNJ_fHiU5VD3wFh2F0dG5ldHOI__________-EZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBSJc2VjcDI1NmsxoQKwNbLb9B0rpePJSV4qWAdO1RL2PXe0Gp4yvmCZbztbnYhzeW5jbmV0cw-DdGNwgiMog3VkcIIjKA","last_seen_p2p_address":"/ip4/10.0.20.20/tcp/9000/p2p/16Uiu2HAm7HXE58UfHDcv4NNYXMRi53xy7rDVNVSCMKBWVRQ1R8RJ","state":"connected","direction":"outbound"}],"meta":{"count":19}}`),
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "json: cannot unmarshal array into Go value of type v1.NodePeers",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.NodePeers
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}

func TestNodePeerJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "GoodNoENR",
			input: []byte(`{"peer_id":"16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/43402","state":"connected","direction":"inbound"}`),
		},
		{
			name:  "GoodWithENR",
			input: []byte(`{"peer_id":"16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","enr":"enr:-MS4QExfvXqHhj-nqAqkg1Sn55uV7YgpRtlImGCvMJkrkbnLDo8sGhecAGid9B3NjXzN3UtGxpOOUqHZVcEDQxkniwoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBWEcXVpY4IjKYlzZWNwMjU2azGhA9mr4bIskWVeMt0dEn4IlQJhOFgOqgR9V3gkHTl1lTioiHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo","last_seen_p2p_address":"/ip4/10.0.20.21/udp/9001/quic-v1/p2p/16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","state":"connected","direction":"outbound"}`),
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "json: cannot unmarshal array into Go value of type v1.NodePeer",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.NodePeer
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
