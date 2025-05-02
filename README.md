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

client := chartmetric.NewClient("<your-refresh-token>", chartmetric.WithRateLimitPerSec(1))
```
- When instantiating the client, an optional _rate limit per second_ argument can be provided (defaults to 1 if not provided). This could correspond to the permitted requests per second of the availed [Developer API plan](https://chartmetric.com/pricing).

### Fetch chart countries

```go
chartCountries, err := client.GetChartCountries(
    ctx,
    chartmetric.ChartPlatformSpotify,
    &chartmetric.GetChartCountriesParams{ChartCountriesType: chartmetric.ChartCountriesTypeTracks},
)
```

### Fetch Spotify chart tracks

```go
spotifyChartTracks, err := client.GetChartTracksSpotify(
    ctx,
    chartmetric.GetChartTracksSpotifyParams{
        Date:        time.Now(),
        CountryCode: chartCountries[0],
        Type:        chartmetric.ChartTypeSpotifyRegional,
        Interval:    chartmetric.ChartIntervalSpotifyWeekly,
    },
)
```

### Fetch TikTok chart entities

```go
tiktokChartEntities, err := client.GetChartEntriesTikTok(
    ctx,
    chartmetric.GetChartEntriesTikTokParams{
        ChartType: chartmetric.ChartTypeTikTokTracks,
        Date:      time.Now().Add(-24 * time.Hour),
        Interval:  chartmetric.Opt(chartmetric.ChartIntervalTikTokDaily),
        Limit:     chartmetric.Opt(10),
    },
)
```