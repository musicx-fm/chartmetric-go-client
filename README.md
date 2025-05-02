# Chartmetric Go Client

Golang client for [Chartmetric API](https://api.chartmetric.com/apidoc/).

## Prerequisites

- Go version >= 1.24

## Installation

```
go get -u github.com/musicx-fm/chartmetric-go-client
```

## Usage

```go
import "github.com/musicx-fm/chartmetric-go-client"

cmClient := chartmetric.NewClient("<your-refresh-token>", chartmetric.WithRateLimitPerSec(1))
```
- When instantiating the client, an optional _rate limit per second_ argument can be provided (defaults to 1 if not provided). This could correspond to the permitted requests per second of the availed [Developer API plan](https://chartmetric.com/pricing).
