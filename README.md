# go-eth2-client

[![Tag](https://img.shields.io/github/tag/attestantio/go-eth2-client.svg)](https://github.com/attestantio/go-eth2-client/releases/)
[![License](https://img.shields.io/github/license/attestantio/go-eth2-client.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/attestantio/go-eth2-client?status.svg)](https://godoc.org/github.com/attestantio/go-eth2-client)
![Lint](https://github.com/attestantio/go-eth2-client/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/attestantio/go-eth2-client)](https://goreportcard.com/report/github.com/attestantio/go-eth2-client)

Go library providing an abstraction to multiple Ethereum 2 beacon nodes.  Its external API follows the official [Ethereum beacon APIs](https://github.com/ethereum/beacon-APIs) specification.

This library is under development; expect APIs and data structures to change until it reaches version 1.0.  In addition, clients' implementations of both their own and the standard API are themselves under development so implementation of the the full API can be incomplete.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-client` is a standard Go module which can be installed with:

```sh
go get github.com/attestantio/go-eth2-client
```

## Support

`go-eth2-client` supports beacon nodes that comply with the standard beacon node API.  To date it has been tested against the following beacon nodes:

  - [Lighthouse](https://github.com/sigp/lighthouse/) minimum version 2.0.0
  - [Nimbus](https://github.com/status-im/nimbus-eth2) minimum version ?
  - [Prysm](https://github.com/prysmaticlabs/prysm) minimum version ?
  - [Teku](https://github.com/consensys/teku) minimum version 21.9.2

Note that the gRPC interface for Prysm is still available in the `grpc` module, however it does not support Altair and is considered deprecated.  The `http` module should be used instead.  gRPC is no longer selected as part of the `auto` module, which is also deprecated in favour of instantiating the `http` module directly.

## Usage

Please read the [Go documentation for this library](https://godoc.org/github.com/attestantio/go-eth2-client) for interface information.

## Example

Below is a complete annotated example to access a beacon node.

```
package main

import (
    "context"
    "fmt"
    
    eth2client "github.com/attestantio/go-eth2-client"
    "github.com/attestantio/go-eth2-client/http"
    "github.com/rs/zerolog"
)

func main() {
    // Provide a cancellable context to the creation function.
    ctx, cancel := context.WithCancel(context.Background())
    client, err := http.New(ctx,
        // WithAddress supplies the address of the beacon node, as a URL.
        http.WithAddress("http://localhost:5052/"),
        // LogLevel supplies the level of logging to carry out.
        http.WithLogLevel(zerolog.WarnLevel),
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Connected to %s\n", client.Name())
    
    // Client functions have their own interfaces.  Not all functions are
    // supported by all clients, so checks should be made for each function when
    // casting the service to the relevant interface.
    if provider, isProvider := client.(eth2client.GenesisProvider); isProvider {
        genesis, err := provider.Genesis(ctx)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Genesis time is %v\n", genesis.GenesisTime)
    }

    // Cancelling the context passed to New() frees up resources held by the
    // client, closes connections, clears handlers, etc.
    cancel()
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/attestantio/go-eth2-client/issues).

## License

[Apache-2.0](LICENSE) © 2020, 2021 Attestant Limited
