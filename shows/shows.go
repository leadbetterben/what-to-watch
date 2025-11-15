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
// corresponds to the index displayed by `GetCurrentlyWatching()`.
// It returns the updated shows slice, a non-empty congratulations message if the
// show was completed, and an error.
func MarkEpisodeWatched(shows []data.Show, listIndex int) ([]data.Show, string, error) {
	if listIndex <= 0 {
		return nil, "", fmt.Errorf("invalid index: %d", listIndex)
	}

	// Build mapping from currently-watching list back to original shows slice
	var watchingIndices []int
	for i, s := range shows {
		if s.CurrentSeries == nil && s.CurrentEpisode == nil {
			continue
		}
		watchingIndices = append(watchingIndices, i)
	}

	if listIndex > len(watchingIndices) {
		return nil, "", fmt.Errorf("index out of range")
	}

	origIdx := watchingIndices[listIndex-1]
	s := &shows[origIdx]

	if s.CurrentSeries == nil || s.CurrentEpisode == nil {
		return nil, "", fmt.Errorf("selected show is not currently being watched")
	}

	curSeries := *s.CurrentSeries
	curEpisode := *s.CurrentEpisode

	// Determine episodes in current series. `Episodes` slice holds episode counts per series.
	if curSeries <= 0 || curSeries > len(s.Episodes) {
		// Defensive: if data is malformed, treat as finished
		s.CurrentSeries = nil
		s.CurrentEpisode = nil
		return shows, fmt.Sprintf("Congratulations! You finished %s.", s.Name), nil
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
			return shows, fmt.Sprintf("Congratulations! You finished %s.", s.Name), nil
		}
	}

	// set updated values
	s.CurrentSeries = &curSeries
	s.CurrentEpisode = &curEpisode
	return shows, "Updated show progress.", nil
}
