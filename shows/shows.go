package shows

import (
	"fmt"

	"what-to-watch/data"
)

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

	if err := data.ValidateCurrentPosition(*s); err != nil {
		return nil, false, err
	}
	if s.CurrentSeries == nil && s.CurrentEpisode == nil {
		return nil, false, fmt.Errorf("selected show is not currently being watched")
	}

	curSeries := *s.CurrentSeries
	curEpisode := *s.CurrentEpisode

	// Determine episodes in current series. `Episodes` slice holds episode counts per series.
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
