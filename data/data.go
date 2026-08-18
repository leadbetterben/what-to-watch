package data

type Show struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Episodes []int  `json:"episodes"`
	Provider string `json:"provider"`
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
