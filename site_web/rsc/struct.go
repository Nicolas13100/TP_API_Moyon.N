package API

// struct, nothing much to explain here

type Artist struct {
	Name       string
	Type       string
	Popularity int
	Followers  struct {
		Total int
	}
}

type Album struct {
	Name          string
	CoverImage    string
	ReleaseDate   string
	NumberOfSongs int
}

type TrackInfo struct {
	Title       string
	AlbumCover  string
	AlbumName   string
	ArtistName  string
	ReleaseDate string
	SpotifyLink string
}
