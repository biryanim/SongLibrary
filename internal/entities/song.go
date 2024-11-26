package entities

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"text"`
	Link        string `json:"link"`
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"text"`
	Link        string `json:"link"`
}
