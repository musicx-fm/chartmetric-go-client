package chartmetric

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

type ChartCountriesType string

const (
	ChartCountriesTypeAlbums  ChartCountriesType = "albums"
	ChartCountriesTypeArtists ChartCountriesType = "artists"
	ChartCountriesTypeTracks  ChartCountriesType = "tracks"
	ChartCountriesTypeTrends  ChartCountriesType = "trends"
	ChartCountriesTypeVideos  ChartCountriesType = "videos"
)

type ChartCountriesSubType string

const (
	ChartCountriesSubTypeDaily        ChartCountriesSubType = "daily"
	ChartCountriesSubTypeNewAlbum     ChartCountriesSubType = "new_album"
	ChartCountriesSubTypeNewTrack     ChartCountriesSubType = "new_track"
	ChartCountriesSubTypePopularAlbum ChartCountriesSubType = "popular_album"
	ChartCountriesSubTypePopularTrack ChartCountriesSubType = "popular_track"
	ChartCountriesSubTypeTop          ChartCountriesSubType = "top"
)

type ChartCountriesDuration string

const (
	ChartCountriesDurationDaily  ChartCountriesDuration = "daily"
	ChartCountriesDurationWeekly ChartCountriesDuration = "weekly"
)

type GetChartCountriesParams struct {
	ChartCountriesType
	ChartCountriesSubType
	ChartCountriesDuration
}

type getChartCountriesResponse struct {
	Obj struct {
		Countries []string `json:"countries"`
	} `json:"obj"`
}

// GetChartCountries fetches the available chart countries for a particular platform.
// Different platforms require different combinations of params (see https://api.chartmetric.com/apidoc/#api-Charts-GetChartCountriesForPlatform-1.0.0)
func (c *Client) GetChartCountries(ctx context.Context, platform ChartPlatform, params *GetChartCountriesParams) ([]string, error) {
	path := fmt.Sprintf("/charts/%s/countries", platform)

	var queryParams map[string]any
	if params != nil {
		queryParams = make(map[string]any)
		if params.ChartCountriesType != "" {
			queryParams["chart_type"] = params.ChartCountriesType
		}
		if params.ChartCountriesSubType != "" {
			queryParams["type"] = params.ChartCountriesSubType
		}
		if params.ChartCountriesDuration != "" {
			queryParams["duration"] = params.ChartCountriesDuration
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

// ==================================================

type ChartTypeTracksSpotify string

const (
	ChartTypeTracksSpotifyRegional ChartTypeTracksSpotify = "regional"
	ChartTypeTracksSpotifyViral    ChartTypeTracksSpotify = "viral"
)

type ChartIntervalTracksSpotify string

const (
	ChartIntervalTracksSpotifyDaily  ChartIntervalTracksSpotify = "daily"
	ChartIntervalTracksSpotifyWeekly ChartIntervalTracksSpotify = "weekly"
)

type GetChartTracksSpotifyParams struct {
	Date        time.Time
	CountryCode string
	Type        ChartTypeTracksSpotify
	Interval    ChartIntervalTracksSpotify
	Offset      Optional[int]
	Latest      Optional[bool]
}

type getChartTracksSpotifyResponse struct {
	Obj struct {
		Length int                 `json:"length"`
		Data   []ChartTrackSpotify `json:"data"`
	} `json:"obj"`
}

type ChartTrackSpotify struct {
	Name                 string                 `json:"name"`
	ISRC                 string                 `json:"isrc"`
	SpotifyTrackID       string                 `json:"spotify_track_id"`
	SpotifyAlbumID       string                 `json:"spotify_album_id"`
	ImageURL             string                 `json:"image_url"`
	SpotifyPopularity    int                    `json:"spotify_popularity"`
	ChartmetricTrackID   int                    `json:"cm_track"`
	SpotifyArtistNames   []string               `json:"spotify_artist_names"`
	Genre                string                 `json:"track_genre"`
	ChartmetricArtistIDs []int                  `json:"cm_artist"`
	ArtistNames          []string               `json:"artist_names"`
	CountryCodes         []string               `json:"code2s"`
	ArtistImages         []string               `json:"artist_images"`
	ArtistCovers         []string               `json:"artist_covers"`
	SpotifyArtistIDs     []string               `json:"spotify_artist_ids"`
	SpotifyTrackIDs      []string               `json:"spotify_track_ids"`
	SpotifyAlbumIDs      []string               `json:"spotify_album_ids"`
	SpotifyDurationMs    int                    `json:"spotify_duration_ms"`
	ChartmetricAlbumIDs  []int                  `json:"album_ids"`
	AlbumNames           []string               `json:"album_names"`
	AlbumUPC             []string               `json:"album_upc"`
	AlbumLabel           []string               `json:"album_label"`
	ReleaseDates         []Date                 `json:"release_dates"`
	Rank                 int                    `json:"rank"`
	AddedAt              time.Time              `json:"added_at"`
	ChartType            ChartTypeTracksSpotify `json:"chart_type"`
	ChartName            string                 `json:"chart_name"`
	CurrentPlays         int                    `json:"current_plays"`
	CountryCode          string                 `json:"code2"`
	Velocity             float64                `json:"velocity"`
	PreRank              int                    `json:"pre_rank"`
	PeakRank             int                    `json:"peak_rank"`
	PeakDate             time.Time              `json:"peak_date"`
	TimeOnChart          int                    `json:"time_on_chart"`
	RankStats            struct {
		Plays     int       `json:"plays"`
		Rank      int       `json:"rank"`
		Timestamp time.Time `json:"timestp"`
	} `json:"rank_stats"`
}

// GetChartTracksSpotify fetches insights for tracks on Spotify charts.
// See https://api.chartmetric.com/apidoc/#api-Charts-GetSpotifyTracksChart-1.0.0.
func (c *Client) GetChartTracksSpotify(ctx context.Context, params GetChartTracksSpotifyParams) ([]ChartTrackSpotify, error) {
	path := "/charts/spotify"

	queryParams := make(map[string]any)
	queryParams["date"] = params.Date.Format(DateFormat)
	queryParams["country_code"] = params.CountryCode
	queryParams["type"] = params.Type
	queryParams["interval"] = params.Interval
	if params.Offset != nil {
		queryParams["offset"] = *params.Offset
	}
	if params.Latest != nil {
		queryParams["latest"] = *params.Latest
	}

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getChartTracksSpotifyResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return response.Obj.Data, nil
}

// ==================================================

type ChartTypeArtistsSpotify string

const (
	ChartTypeArtistsSpotifyMonthlyListeners ChartTypeArtistsSpotify = "monthly_listeners"
	ChartTypeArtistsSpotifyPopularity       ChartTypeArtistsSpotify = "popularity"
	ChartTypeArtistsSpotifyFollowers        ChartTypeArtistsSpotify = "followers"
	ChartTypeArtistsSpotifyPlaylistCount    ChartTypeArtistsSpotify = "playlist_count"
	ChartTypeArtistsSpotifyPlaylistReach    ChartTypeArtistsSpotify = "playlist_reach"
)

type ChartIntervalArtistsSpotify string

const (
	ChartIntervalArtistsSpotifyDaily   ChartIntervalArtistsSpotify = "daily"
	ChartIntervalArtistsSpotifyWeekly  ChartIntervalArtistsSpotify = "weekly"
	ChartIntervalArtistsSpotifyMonthly ChartIntervalArtistsSpotify = "monthly"
)

type GetChartArtistSpotifyParams struct {
	Date     time.Time
	Type     ChartTypeArtistsSpotify
	Interval ChartIntervalArtistsSpotify
	Offset   Optional[int]
	Latest   Optional[bool]
}

type getChartArtistsSpotifyResponse struct {
	Obj struct {
		Length int                  `json:"length"`
		Data   []ChartArtistSpotify `json:"data"`
	} `json:"obj"`
}

type ChartArtistSpotify struct {
	ID                     int      `json:"id"`
	Name                   string   `json:"name"`
	ImageURL               string   `json:"image_url"`
	ISNI                   string   `json:"isni"`
	CountryCode            string   `json:"code2"`
	HometownCity           string   `json:"hometown_city"`
	CurrentCity            string   `json:"current_city"`
	Followers              int      `json:"sp_followers"`
	Popularity             int      `json:"sp_popularity"`
	MonthlyListeners       int      `json:"sp_monthly_listeners"`
	DeezerFans             int      `json:"deezer_fans"`
	Tags                   []string `json:"tags"`
	SpotifyArtistIDs       []string `json:"spotify_artist_ids"`
	ITunesArtistIDs        []int    `json:"itunes_artist_ids"`
	DeezerArtistIDs        []string `json:"deezer_artist_ids"`
	AmazonArtistIDs        []string `json:"amazon_artist_ids"`
	ChartmetricArtistRank  int      `json:"cm_artist_rank"`
	ChartmetricArtistScore float64  `json:"cm_artist_score"`
	Timestamp              Date     `json:"timestp"`
	MonthlyPlaylistReach   int      `json:"monthly_playlist_reach"`
	Rank                   int      `json:"rank"`
	Change                 string   `json:"change"`
	Streak                 int      `json:"streak"`
	PeakDate               Date     `json:"peak_date"`
	PeakRank               int      `json:"peak_rank"`
	TimeOnChart            int      `json:"time_on_chart"`
	RankStats              []struct {
		Rank int  `json:"rank"`
		Date Date `json:"date"`
	} `json:"rankStats"`
}

// GetChartArtistsSpotify fetches insights for artists on Spotify charts.
// See https://api.chartmetric.com/apidoc/#api-Charts-GetSpotifyArtistsChart-1.0.0.
func (c *Client) GetChartArtistsSpotify(ctx context.Context, params GetChartArtistSpotifyParams) ([]ChartArtistSpotify, error) {
	path := "/charts/spotify/artists"

	queryParams := make(map[string]any)
	queryParams["date"] = params.Date.Format(DateFormat)
	queryParams["type"] = params.Type
	queryParams["interval"] = params.Interval
	if params.Offset != nil {
		queryParams["offset"] = *params.Offset
	}
	if params.Latest != nil {
		queryParams["latest"] = *params.Latest
	}

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getChartArtistsSpotifyResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return response.Obj.Data, nil
}

// ==================================================

type ChartTypeTikTok string

const (
	ChartTypeTikTokTracks ChartTypeTikTok = "tracks"
	ChartTypeTikTokVideos ChartTypeTikTok = "videos"
	ChartTypeTikTokUsers  ChartTypeTikTok = "users"
)

type ChartIntervalTikTok string

const (
	ChartIntervalTikTokDaily   ChartIntervalTikTok = "daily"
	ChartIntervalTikTokWeekly  ChartIntervalTikTok = "weekly"
	ChartIntervalTikTokAllTime ChartIntervalTikTok = "all_time"
)

type ChartUserTypeTikTok string

const (
	ChartUserTypeTikTokLikes     ChartUserTypeTikTok = "likes"
	ChartUserTypeTikTokFollowers ChartUserTypeTikTok = "followers"
)

type GetChartEntriesTikTokParams struct {
	ChartType    ChartTypeTikTok
	Date         time.Time
	Interval     Optional[ChartIntervalTikTok]
	UserType     Optional[ChartUserTypeTikTok]
	Limit        Optional[int]
	Offset       Optional[int]
	Latest       Optional[bool]
	CountryChart Optional[bool]
	CountryCode  Optional[string]
}

type getChartEntriesTikTokResponse struct {
	Obj struct {
		Length int                `json:"length"`
		Data   []ChartEntryTikTok `json:"data"`
	} `json:"obj"`
}

type ChartEntryTikTok struct {
	TikTokTrackID       string    `json:"tiktok_track_id"`
	Name                string    `json:"name"`
	ISRC                string    `json:"isrc"`
	ImageURL            string    `json:"image_url"`
	ChartmetricTrackID  int       `json:"cm_track"`
	TikTokArtistNames   []string  `json:"tiktok_artist_names"`
	TikTokTrackIDs      []string  `json:"tiktok_track_ids"`
	ChartmetricAlbumIDs []int     `json:"album_ids"`
	AlbumNames          []string  `json:"album_names"`
	AlbumUPC            []string  `json:"album_upc"`
	AlbumLabel          []string  `json:"album_label"`
	ReleaseDates        []Date    `json:"release_dates"`
	Rank                int       `json:"rank"`
	Posts               int       `json:"posts"`
	AddedAt             time.Time `json:"added_at"`
	PreRank             int       `json:"pre_rank"`
	PeakRank            int       `json:"peak_rank"`
	PeakDate            time.Time `json:"peak_date"`
	TimeOnChart         int       `json:"time_on_chart"`
	RankStats           []struct {
		Rank      int       `json:"rank"`
		Posts     int       `json:"posts"`
		Timestamp time.Time `json:"timestp"`
	} `json:"rankStats"`
	Views      int `json:"views"`
	ViewsDelta int `json:"views_prev"`
}

// GetChartEntriesTikTok fetches information for some TikTok chart.
// See https://api.chartmetric.com/apidoc/#api-Charts-GetTiktokTracksChart-1.0.0.
func (c *Client) GetChartEntriesTikTok(ctx context.Context, params GetChartEntriesTikTokParams) ([]ChartEntryTikTok, error) {
	path := fmt.Sprintf("/charts/tiktok/%s", params.ChartType)

	queryParams := make(map[string]any)
	queryParams["date"] = params.Date.Format(DateFormat)
	if params.Interval != nil {
		queryParams["interval"] = *params.Interval
	}
	if params.UserType != nil {
		queryParams["type"] = *params.UserType
	}
	if params.Limit != nil {
		queryParams["limit"] = *params.Limit
	}
	if params.Offset != nil {
		queryParams["offset"] = *params.Offset
	}
	if params.Latest != nil {
		queryParams["latest"] = *params.Latest
	}
	if params.CountryChart != nil {
		queryParams["country_chart"] = *params.CountryChart
	}
	if params.CountryCode != nil {
		queryParams["code2"] = *params.CountryCode
	}

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getChartEntriesTikTokResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return response.Obj.Data, nil
}

// ==================================================

type ChartTypeAppleMusic string

const (
	ChartTypeAppleMusicAlbums ChartTypeAppleMusic = "albums"
	ChartTypeAppleMusicTracks ChartTypeAppleMusic = "tracks"
	ChartTypeAppleMusicVideos ChartTypeAppleMusic = "videos"
)

type ChartTracksTypeAppleMusic string

const (
	ChartTracksTypeAppleMusicDaily ChartTracksTypeAppleMusic = "daily"
	ChartTracksTypeAppleMusicTop   ChartTracksTypeAppleMusic = "top"
)

type GetChartEntriesAppleMusicParams struct {
	ChartType   ChartTypeAppleMusic
	Type        Optional[ChartTracksTypeAppleMusic]
	CountryCode string
	CityID      Optional[string]
	Date        time.Time
	Genre       Optional[string]
	Offset      Optional[int]
	Latest      Optional[bool]
}

type getChartEntriesAppleMusicResponse struct {
	Obj struct {
		Length int                    `json:"length"`
		Data   []ChartEntryAppleMusic `json:"data"`
	} `json:"obj"`
}

type ChartEntryAppleMusic struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	ISRC                 string    `json:"isrc"`
	ITunesAlbumID        any       `json:"itunes_album_id"`
	ComposerName         string    `json:"composer_name"`
	ImageURL             string    `json:"image_url"`
	ChartmetricTrackID   int       `json:"cm_track"`
	TrackGenre           string    `json:"track_genre"`
	ChartmetricArtistIDs []int     `json:"cm_artist"`
	ArtistNames          []string  `json:"artist_names"`
	CountryCodes         []string  `json:"code2s"`
	ArtistImages         []string  `json:"artist_images"`
	ArtistCovers         []string  `json:"artist_covers"`
	ITunesArtistIDs      []int     `json:"itunes_artist_ids"`
	ITunesArtistNames    []string  `json:"itunes_artist_names"`
	ITunesTrackIDs       []string  `json:"itunes_track_ids"`
	ITunesAlbumIDs       []string  `json:"itunes_album_ids"`
	Storefronts          []string  `json:"storefronts"`
	ChartmetricAlbumIDs  []int     `json:"album_ids"`
	AlbumNames           []string  `json:"album_names"`
	AlbumUPC             []string  `json:"album_upc"`
	AlbumLabel           []string  `json:"album_label"`
	UPC                  string    `json:"upc"`
	ReleaseDates         []Date    `json:"release_dates"`
	Rank                 int       `json:"rank"`
	AddedAt              time.Time `json:"added_at"`
	CountryCode          string    `json:"code2"`
	Country              string    `json:"country"`
	ITunes               int       `json:"itunes"`
	Velocity             float64   `json:"velocity"`
	PreRank              int       `json:"pre_rank"`
	PeakRank             int       `json:"peak_rank"`
	PeakDate             time.Time `json:"peak_date"`
	TimeOnChart          int       `json:"time_on_chart"`
	RankStats            struct {
		Rank      int       `json:"rank"`
		Timestamp time.Time `json:"timestp"`
	} `json:"rank_stats"`
}

// GetChartEntriesAppleMusic fetches insights for albums on Apple Music charts.
// See https://api.chartmetric.com/apidoc/#api-Charts-GetAppleMusicChart.
func (c *Client) GetChartEntriesAppleMusic(ctx context.Context, params GetChartEntriesAppleMusicParams) ([]ChartEntryAppleMusic, error) {
	path := fmt.Sprintf("/charts/applemusic/%s", params.ChartType)

	queryParams := make(map[string]any)
	if params.Type != nil {
		queryParams["type"] = *params.Type
	}
	queryParams["country_code"] = params.CountryCode
	if params.CityID != nil {
		queryParams["city_id"] = *params.CityID
	}
	queryParams["date"] = params.Date.Format(DateFormat)
	if params.Genre != nil {
		queryParams["genre"] = *params.Genre
	}
	if params.Offset != nil {
		queryParams["offset"] = *params.Offset
	}
	if params.Latest != nil {
		queryParams["latest"] = *params.Latest
	}

	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	var response getChartEntriesAppleMusicResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return response.Obj.Data, nil
}
