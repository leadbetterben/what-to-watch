package shows

import (
	"fmt"
	"strconv"

	"what-to-watch/data"
	"what-to-watch/db"
)

// GetCurrentlyWatching returns a slice of shows that the user is currently watching,
// including their current series and episode information.
func GetCurrentlyWatching() ([]data.Show, error) {
	shows, err := db.ReadShows()
	if err != nil {
		return nil, fmt.Errorf("GetCurrentlyWatching: error reading shows \n err=%w", err)
	}

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
