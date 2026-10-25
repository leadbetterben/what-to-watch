package data

import "fmt"

// ValidateEpisodeCounts checks that a show has a positive episode count for every series.
func ValidateEpisodeCounts(show Show) error {
	if len(show.Episodes) == 0 {
		return fmt.Errorf("show %q must have episode counts", show.Name)
	}
	for i, count := range show.Episodes {
		if count <= 0 {
			return fmt.Errorf("show %q has invalid episodes[%d] count %d: must be positive", show.Name, i+1, count)
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

	if err := ValidateEpisodeCounts(show); err != nil {
		return err
	}
	return validateEpisodePosition(show.Name, "currentSeries", "currentEpisode", *show.CurrentSeries, *show.CurrentEpisode, show.Episodes)
}

// ValidateSpecialEpisodes checks the seasonal episode metadata on a show.
// Episode-count bounds are checked against the show's episode counts.
func ValidateSpecialEpisodes(show Show) error {
	for i, special := range show.SpecialEpisodes {
		field := fmt.Sprintf("specialEpisodes[%d]", i+1)
		switch special.Category {
		case SpecialCategoryHalloween, SpecialCategoryChristmas:
		default:
			return fmt.Errorf("show %q has invalid %s.category %q", show.Name, field, special.Category)
		}

		if err := validateEpisodePosition(show.Name, field+".series", field+".episode", special.Series, special.Episode, show.Episodes); err != nil {
			return err
		}
	}
	return nil
}

// validateEpisodePosition checks a series and episode position against the episode counts.
func validateEpisodePosition(showName, seriesField, episodeField string, series, episode int, episodes []int) error {
	if series <= 0 {
		return fmt.Errorf("show %q has invalid %s %d: must be positive", showName, seriesField, series)
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
