package APi

type Artist struct {
	Name       string
	Type       string
	Popularity int
	Followers  struct {
		Total int
	}
}
