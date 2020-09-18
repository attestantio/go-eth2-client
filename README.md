# go-eth2-client

[![Tag](https://img.shields.io/github/tag/attestantio/go-eth2-client.svg)](https://github.com/attestantio/go-eth2-client/releases/)
[![License](https://img.shields.io/github/license/attestantio/go-eth2-client.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/attestantio/go-eth2-client?status.svg)](https://godoc.org/github.com/attestantio/go-eth2-client)
![Lint](https://github.com/attestantio/go-eth2-client/workflows/golangci-lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/attestantio/go-eth2-client)](https://goreportcard.com/report/github.com/attestantio/go-eth2-client)

Go library providing an abstraction to multiple Ethereum 2 beacon nodes.  Its external API follows the official [Ethereum 2 APIs](https://github.com/ethereum/eth2.0-APIs) specification.

This library is under development; expect APIs and data structures to change until it reaches version 1.0.

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

`go-eth2-client` supports multiple beacon nodes.  At current it provides support for the following:

  - [Prysm](https://github.com/prysmaticlabs/prysm) using its GRPC interface
  - [Lighthouse](https://github.com/sigp/lighthouse/) using its HTTP interface
  - [Teku](https://github.com/pegasyseng/teku) using its HTTP interface


## Usage

`go-eth2-client` provides independent implementations for each beacon node interface, however it is generally easier to use the `auto` interface, as that will automatically select the correct client given the supplied address.

Please read the [Go documentation for this library](https://godoc.org/github.com/attestantio/go-eth2-client) for interface information.

## Example

Below is a complete annotated example to access a beacon node.

```
package main

import (
    "context"
    "fmt"
    
    eth2client "github.com/attestantio/go-eth2-client"
    "github.com/attestantio/go-eth2-client/auto"
    "github.com/rs/zerolog"
)

func main() {
    // Provide a cancellable context to the creation function.
    ctx, cancel := context.WithCancel(context.Background())
    client, err := auto.New(ctx,
        // WithAddress supplies the address of the beacon node, in host:port format.
        auto.WithAddress("localhost:4000"),
        // LogLevel supplies the level of logging to carry out.
        auto.WithLogLevel(zerolog.WarnLevel),
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Connected to %s\n", client.Name())
    
    // Client functions have their own interfaces.  Not all functions are
    // supported by all clients, so checks should be made for each function when
    // casting the service to the relevant interface.
    if provider, isProvider := client.(eth2client.SlotsPerEpochProvider); isProvider {
        slotsPerEpoch, err := provider.SlotsPerEpoch(ctx)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Slots per epochs is %d\n", slotsPerEpoch)
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

[Apache-2.0](LICENSE) © 2020 Attestant Limited
