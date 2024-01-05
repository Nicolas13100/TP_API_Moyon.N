package APi

type Artist struct {
	Name       string
	Type       string
	Popularity int
	Followers  struct {
		Total int
	}
}

type JulAlbum struct {
	Artists []struct {
		AlbumGroup string `json:"album_group"`
		AlbumType  string `json:"album_type"`
		Artists    []struct {
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"artists"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href                 string        `json:"href"`
		ID                   string        `json:"id"`
		Images               []interface{} `json:"images"`
		IsPlayable           bool          `json:"is_playable"`
		Name                 string        `json:"name"`
		ReleaseDate          string        `json:"release_date"`
		ReleaseDatePrecision string        `json:"release_date_precision"`
		TotalTracks          int           `json:"total_tracks"`
		Type                 string        `json:"type"`
		URI                  string        `json:"uri"`
	}
}
