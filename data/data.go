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
	// Rewatch is true when the show is available for watching again.
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

// ValidateCurrentPosition checks the stored position of a show that is marked as currently watching.
// A non-current show has no position to validate.
func ValidateCurrentPosition(show Show) error {
	if show.CurrentSeries == nil && show.CurrentEpisode == nil {
		return nil
	}
	if show.CurrentSeries == nil || show.CurrentEpisode == nil {
		return fmt.Errorf("show %q must have both currentSeries and currentEpisode", show.Name)
	}

	if *show.CurrentSeries <= 0 {
		return fmt.Errorf("show %q has invalid currentSeries %d: must be positive", show.Name, *show.CurrentSeries)
	}
	if *show.CurrentSeries > len(show.Episodes) {
		return fmt.Errorf("show %q has currentSeries %d outside episodes bounds 1..%d", show.Name, *show.CurrentSeries, len(show.Episodes))
	}
	if *show.CurrentEpisode <= 0 {
		return fmt.Errorf("show %q has invalid currentEpisode %d: must be positive", show.Name, *show.CurrentEpisode)
	}

	for i, count := range show.Episodes {
		if count <= 0 {
			return fmt.Errorf("show %q has invalid episodes[%d] count %d: must be positive", show.Name, i+1, count)
		}
	}

	episodesInSeries := show.Episodes[*show.CurrentSeries-1]
	if *show.CurrentEpisode > episodesInSeries {
		return fmt.Errorf("show %q has currentEpisode %d outside series %d episode bounds 1..%d", show.Name, *show.CurrentEpisode, *show.CurrentSeries, episodesInSeries)
	}
	return nil
}

type Film struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Provider string `json:"provider"`
}
