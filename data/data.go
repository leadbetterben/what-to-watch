package data

import "fmt"

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

// ValidateSpecialEpisodes checks the seasonal episode metadata on a show.
// Episode-count bounds are checked when counts are available.
func ValidateSpecialEpisodes(show Show) error {
	for i, special := range show.SpecialEpisodes {
		field := fmt.Sprintf("specialEpisodes[%d]", i+1)
		switch special.Category {
		case SpecialCategoryHalloween, SpecialCategoryChristmas:
		default:
			return fmt.Errorf("show %q has invalid %s.category %q", show.Name, field, special.Category)
		}
		if special.Series <= 0 {
			return fmt.Errorf("show %q has invalid %s.series %d: must be positive", show.Name, field, special.Series)
		}
		if special.Episode <= 0 {
			return fmt.Errorf("show %q has invalid %s.episode %d: must be positive", show.Name, field, special.Episode)
		}
		if len(show.Episodes) == 0 {
			continue
		}
		if special.Series > len(show.Episodes) {
			return fmt.Errorf("show %q has %s.series %d outside episodes bounds 1..%d", show.Name, field, special.Series, len(show.Episodes))
		}
		episodeCount := show.Episodes[special.Series-1]
		if special.Episode > episodeCount {
			return fmt.Errorf("show %q has %s.episode %d outside series %d episode bounds 1..%d", show.Name, field, special.Episode, special.Series, episodeCount)
		}
	}
	return nil
}

type Film struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Provider string `json:"provider"`
}
