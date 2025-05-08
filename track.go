package chartmetric

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type TrackPlatform string

const (
	TrackPlatformAmazon      TrackPlatform = "amazon"
	TrackPlatformBeatport    TrackPlatform = "beatport"
	TrackPlatformDeezer      TrackPlatform = "deezer"
	TrackPlatformGenius      TrackPlatform = "genius"
	TrackPlatformITunes      TrackPlatform = "itunes"
	TrackPlatformShazam      TrackPlatform = "shazam"
	TrackPlatformSoundCloud  TrackPlatform = "soundcloud"
	TrackPlatformSpotify     TrackPlatform = "spotify"
	TrackPlatformTikTok      TrackPlatform = "tiktok"
	TrackPlatformQQ          TrackPlatform = "qq"
	TrackPlatformYouTube     TrackPlatform = "youtube"
	TrackPlatformChartmetric TrackPlatform = "chartmetric"
	TrackPlatformISRC        TrackPlatform = "isrc"
)

type getTrackIDsResponse struct {
	Obj []TrackIDs `json:"obj"`
}

type TrackIDs struct {
	ISRC           string   `json:"isrc"`
	ChartmetricIDs []int    `json:"chartmetric_ids"`
	SpotifyIDs     []string `json:"spotify_ids"`
	ITunesIDs      []string `json:"itunes_ids"`
	DeezerIDs      []string `json:"deezer_ids"`
	AmazonIDs      []string `json:"amazon_ids"`
	YouTubeIDs     []string `json:"youtube_ids"`
	SoundCloudIDs  []string `json:"soundcloud_ids"`
	ShazamIDs      []string `json:"shazam_ids"`
	TikTokIDs      []string `json:"tiktok_ids"`
	BeatportIDs    []int    `json:"beatport_ids"`
	QQIDs          []int    `json:"qq_ids"`
	GeniusIDs      []int    `json:"genius_ids"`
}

// GetTrackIDs accepts a platform and a track's ID on that platform, then returns
// the track IDs across different platforms for that same track.
// See https://api.chartmetric.com/apidoc/#api-Track-getTrackIDs.
func (c *Client) GetTrackIDs(ctx context.Context, platform TrackPlatform, id string) (*TrackIDs, error) {
	path := fmt.Sprintf("/track/%s/%s/get-ids", platform, url.PathEscape(id))

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getTrackIDsResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if len(response.Obj) == 0 {
		return nil, fmt.Errorf("no track IDs found for platform %s and ID %s", platform, id)
	}

	return &response.Obj[0], nil
}
