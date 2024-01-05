package APi

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
