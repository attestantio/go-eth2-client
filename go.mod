module github.com/attestantio/go-eth2-client

go 1.14

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/ferranbt/fastssz v0.1.2
	github.com/goccy/go-yaml v1.9.2
	github.com/golang/snappy v0.0.4
	github.com/klauspost/cpuid/v2 v2.1.2 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210809151128-385d8c5e3fb7
	github.com/r3labs/sse/v2 v2.7.4
	github.com/rs/zerolog v1.26.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gotest.tools v2.2.0+incompatible
)

retract (
	v1.15.2 // Retraction for 1.15.1.
	v1.15.1 // Incorrect release number.
)
