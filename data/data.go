package data

// SpecialCategory identifies a supported seasonal episode category.
type SpecialCategory string

const (
	SpecialCategoryHalloween SpecialCategory = "halloween"
	SpecialCategoryChristmas SpecialCategory = "christmas"
)

// SpecialEpisode identifies a curated seasonal episode within a show.
type SpecialEpisode struct {
	Category SpecialCategory `json:"category"`
	Series   int             `json:"series"`
	Episode  int             `json:"episode"`
}

// ShowSpecialEpisode is a seasonal episode with enough parent-show context to locate it.
type ShowSpecialEpisode struct {
	ShowName string          `json:"showName"`
	Provider string          `json:"provider"`
	Category SpecialCategory `json:"category"`
	Series   int             `json:"series"`
	Episode  int             `json:"episode"`
}

type Show struct {
	Name            string           `json:"name"`
	Genre           string           `json:"genre"`
	Episodes        []int            `json:"episodes"`
	Provider        string           `json:"provider"`
	SpecialEpisodes []SpecialEpisode `json:"specialEpisodes,omitempty"`
	// Rewatch is true when the show is queued for watching again.
	Rewatch bool `json:"rewatch,omitempty"`
	// CurrentSeries is only set if the user is currently watching this show
	CurrentSeries *int   `json:"currentSeries,omitempty"`
	Series        string `json:"-"`
	// CurrentEpisode is only set if the user is currently watching this show
	CurrentEpisode *int   `json:"currentEpisode,omitempty"`
	Episode        string `json:"-"`
}

type Film struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Provider string `json:"provider"`
}
