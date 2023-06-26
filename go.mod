module github.com/attestantio/go-eth2-client

go 1.20

require (
	github.com/ferranbt/fastssz v0.1.2
	github.com/goccy/go-yaml v1.9.2
	github.com/golang/snappy v0.0.4
	github.com/holiman/uint256 v1.2.2
	github.com/huandu/go-clone/generic v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210809151128-385d8c5e3fb7
	github.com/r3labs/sse/v2 v2.7.4
	github.com/rs/zerolog v1.26.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/huandu/go-clone v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.1.2 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

retract (
	v1.15.2 // Retraction for 1.15.1.
	v1.15.1 // Incorrect release number.
)
