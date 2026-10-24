package data

import "fmt"

type Film struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Provider string `json:"provider"`
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

// ValidateCurrentPosition checks the stored position of a show that is marked as currently watching.
// A non-current show has no position to validate.
func ValidateCurrentPosition(show Show) error {
	if show.CurrentSeries == nil && show.CurrentEpisode == nil {
		return nil
	}
	if show.CurrentSeries == nil || show.CurrentEpisode == nil {
		return fmt.Errorf("show %q must have both currentSeries and currentEpisode", show.Name)
	}

	for i, count := range show.Episodes {
		if count <= 0 {
			return fmt.Errorf("show %q has invalid episodes[%d] count %d: must be positive", show.Name, i+1, count)
		}
	}
	return validateEpisodePosition(show.Name, "currentSeries", "currentEpisode", *show.CurrentSeries, *show.CurrentEpisode, show.Episodes, false)
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

		if err := validateEpisodePosition(show.Name, field+".series", field+".episode", special.Series, special.Episode, show.Episodes, true); err != nil {
			return err
		}
	}
	return nil
}

// validateEpisodePosition checks a series and episode position against the episode counts.
// Some callers can validate positive positions even when the show has no episode counts available.
func validateEpisodePosition(showName, seriesField, episodeField string, series, episode int, episodes []int, allowUnknownCounts bool) error {
	if series <= 0 {
		return fmt.Errorf("show %q has invalid %s %d: must be positive", showName, seriesField, series)
	}
	if len(episodes) == 0 && allowUnknownCounts {
		if episode <= 0 {
			return fmt.Errorf("show %q has invalid %s %d: must be positive", showName, episodeField, episode)
		}
		return nil
	}
	if series > len(episodes) {
		return fmt.Errorf("show %q has %s %d outside episodes bounds 1..%d", showName, seriesField, series, len(episodes))
	}
	if episode <= 0 {
		return fmt.Errorf("show %q has invalid %s %d: must be positive", showName, episodeField, episode)
	}
	episodeCount := episodes[series-1]
	if episode > episodeCount {
		return fmt.Errorf("show %q has %s %d outside series %d episode bounds 1..%d", showName, episodeField, episode, series, episodeCount)
	}
	return nil
}
