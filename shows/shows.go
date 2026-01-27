package shows

import (
	"fmt"
	"strconv"

	"what-to-watch/data"
)

// GetCurrentlyWatching returns a slice of shows that the user is currently watching,
// including their current series and episode information.
func GetCurrentlyWatching(shows []data.Show) ([]data.Show, error) {
	var watching []data.Show
	for _, s := range shows {
		if s.CurrentSeries == nil && s.CurrentEpisode == nil {
			continue
		}

		series := "-"
		if s.CurrentSeries != nil {
			series = strconv.Itoa(*s.CurrentSeries)
		}
		s.Series = series

		episode := "-"
		if s.CurrentEpisode != nil {
			episode = strconv.Itoa(*s.CurrentEpisode)
		}
		s.Episode = episode

		watching = append(watching, s)
	}

	return watching, nil
}

// MarkEpisodeWatched updates the provided shows slice when the user reports they've
// watched the next episode of a show. The parameter `listIndex` is 1-based and
// corresponds to the index displayed by currently watching shows.
// It returns the updated shows slice, a boolean if the show was completed, and an error.
func MarkEpisodeWatched(shows []data.Show, listIndex int) ([]data.Show, bool, error) {
	if listIndex <= 0 {
		return nil, false, fmt.Errorf("invalid index: %d", listIndex)
	}

	if listIndex > len(shows) {
		return nil, false, fmt.Errorf("index out of range")
	}

	s := &shows[listIndex-1]

	if s.CurrentSeries == nil || s.CurrentEpisode == nil {
		return nil, false, fmt.Errorf("selected show is not currently being watched")
	}

	curSeries := *s.CurrentSeries
	curEpisode := *s.CurrentEpisode

	// Determine episodes in current series. `Episodes` slice holds episode counts per series.
	if curSeries <= 0 || curSeries > len(s.Episodes) {
		// Defensive: if data is malformed, treat as finished
		s.CurrentSeries = nil
		s.CurrentEpisode = nil
		return shows, true, nil
	}

	episodesInSeries := s.Episodes[curSeries-1]

	// increment episode
	curEpisode++

	if curEpisode > episodesInSeries {
		// rollover to next series
		curEpisode = 1
		curSeries++
		if curSeries > len(s.Episodes) {
			// finished the show
			s.CurrentSeries = nil
			s.CurrentEpisode = nil
			return shows, true, nil
		}
	}

	// set updated values
	s.CurrentSeries = &curSeries
	s.CurrentEpisode = &curEpisode
	return shows, false, nil
}

// GetUniqueGenres returns a sorted list of unique genres from all shows
func GetUniqueGenres(shows []data.Show) []string {
	genreMap := make(map[string]bool)
	for _, s := range shows {
		if s.Genre != "" {
			genreMap[s.Genre] = true
		}
	}

	var genres []string
	for genre := range genreMap {
		genres = append(genres, genre)
	}
	return genres
}

// GetUnwatchedShowsByGenre returns all shows from a given genre that haven't been watched
// (i.e., shows without CurrentSeries and CurrentEpisode set)
func GetUnwatchedShowsByGenre(shows []data.Show, genre string) []data.Show {
	var unwatched []data.Show
	for _, s := range shows {
		if s.Genre == genre && s.CurrentSeries == nil && s.CurrentEpisode == nil {
			unwatched = append(unwatched, s)
		}
	}
	return unwatched
}
