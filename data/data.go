package data

type Show struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Episodes []int  `json:"episodes"`
	Provider string `json:"provider"`
	// CurrentSeries is only set if the user is currently watching this show
	CurrentSeries *int `json:"currentSeries,omitempty"`
	Series        string
	// CurrentEpisode is only set if the user is currently watching this show
	CurrentEpisode *int `json:"currentEpisode,omitempty"`
	Episode        string
}
