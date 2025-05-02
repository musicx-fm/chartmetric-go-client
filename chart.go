package chartmetric

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChartPlatform string

const (
	ChartPlatformAirplay    ChartPlatform = "airplay"
	ChartPlatformAmazon     ChartPlatform = "amazon"
	ChartPlatformAppleMusic ChartPlatform = "applemusic"
	ChartPlatformDeezer     ChartPlatform = "deezer"
	ChartPlatformITunes     ChartPlatform = "itunes"
	ChartPlatformShazam     ChartPlatform = "shazam"
	ChartPlatformSoundCloud ChartPlatform = "soundcloud"
	ChartPlatformSpotify    ChartPlatform = "spotify"
	ChartPlatformTikTok     ChartPlatform = "tiktok"
	ChartPlatformYouTube    ChartPlatform = "youtube"
)

type ChartType string

const (
	ChartTypeAlbums  ChartType = "albums"
	ChartTypeArtists ChartType = "artists"
	ChartTypeTracks  ChartType = "tracks"
	ChartTypeTrends  ChartType = "trends"
	ChartTypeVideos  ChartType = "videos"
)

type ChartSubType string

const (
	ChartSubTypeDaily        ChartSubType = "daily"
	ChartSubTypeNewAlbum     ChartSubType = "new_album"
	ChartSubTypeNewTrack     ChartSubType = "new_track"
	ChartSubTypePopularAlbum ChartSubType = "popular_album"
	ChartSubTypePopularTrack ChartSubType = "popular_track"
	ChartSubTypeTop          ChartSubType = "top"
)

type ChartTypeDuration string

const (
	ChartTypeDurationDaily  ChartTypeDuration = "daily"
	ChartTypeDurationWeekly ChartTypeDuration = "weekly"
)

type GetChartCountriesParams struct {
	ChartType
	ChartSubType
	ChartTypeDuration
}

type getChartCountriesResponse struct {
	Obj struct {
		Countries []string `json:"countries"`
	} `json:"obj"`
}

// GetChartCountries can be used to fetch the available chart countries for a particular platform.
// Different platforms require different combinations of params (see https://api.chartmetric.com/apidoc/#api-Charts-GetChartCountriesForPlatform-1.0.0)
func (c *Client) GetChartCountries(ctx context.Context, platform ChartPlatform, params *GetChartCountriesParams) ([]string, error) {
	path := fmt.Sprintf("/charts/%s/countries", platform)

	var queryParams map[string]any
	if params != nil {
		queryParams = make(map[string]any)
		if params.ChartType != "" {
			queryParams["chart_type"] = params.ChartType
		}
		if params.ChartSubType != "" {
			queryParams["type"] = params.ChartSubType
		}
		if params.ChartTypeDuration != "" {
			queryParams["duration"] = params.ChartTypeDuration
		}
	}

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getChartCountriesResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return response.Obj.Countries, nil
}
