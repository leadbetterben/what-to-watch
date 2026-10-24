package handlers

import (
	"fmt"
	"math/rand"
	"strconv"

	"what-to-watch/data"
	"what-to-watch/db"
)

// GetCurrentlyWatchingShows retrieves the list of currently watching shows
func GetCurrentlyWatchingShows() ([]data.Show, error) {
	shows, err := db.ReadCurrentShows()
	if err != nil {
		return nil, fmt.Errorf("error reading shows: %w", err)
	}

	var watching []data.Show
	for _, s := range shows {
		if err := data.ValidateCurrentPosition(s); err != nil {
			return nil, fmt.Errorf("error validating shows: %w", err)
		}
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

// MarkShowWatched marks an episode as watched and updates the data store
// idx is 1-based index from the currently watching list
func MarkShowWatched(idx int) (bool, error) {
	completed, err := db.UpdateCurrentShows(func(shows []data.Show) ([]data.Show, bool, error) {
		if idx <= 0 {
			return nil, false, fmt.Errorf("invalid index: %d", idx)
		}

		if idx > len(shows) {
			return nil, false, fmt.Errorf("index out of range")
		}

		s := &shows[idx-1]

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
	})
	if err != nil {
		return false, err
	}

	return completed, nil
}

// GetAllFilms retrieves the list of all films
func GetAllFilms() ([]data.Film, error) {
	films, err := db.ReadFilms()
	if err != nil {
		return nil, fmt.Errorf("error reading films: %w", err)
	}

	return films, nil
}

// GetAvailableGenres retrieves a list of unique genres from all shows
func GetAvailableGenres() ([]string, error) {
	shows, err := db.ReadAllShows()
	if err != nil {
		return nil, fmt.Errorf("GetAvailableGenres: error reading shows: %w", err)
	}

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
	return genres, nil
}

// GetUnwatchedShowsByGenre retrieves all unwatched shows for a given genre
func GetUnwatchedShowsByGenre(genre string) ([]data.Show, error) {
	shows, err := db.ReadUnwatchedShows()
	if err != nil {
		return nil, fmt.Errorf("GetUnwatchedShowsByGenre: error reading shows: %w", err)
	}

	var unwatched []data.Show
	for _, s := range shows {
		if s.Genre == genre && s.CurrentSeries == nil && s.CurrentEpisode == nil && !s.Rewatch {
			unwatched = append(unwatched, s)
		}
	}
	return unwatched, nil
}

func GetRandomSpecialEpisode(category data.SpecialCategory) (data.ShowSpecialEpisode, error) {
	episodes, err := GetSpecialEpisodes(category)
	if err != nil {
		return data.ShowSpecialEpisode{}, err
	}
	if len(episodes) == 0 {
		return data.ShowSpecialEpisode{}, fmt.Errorf("no %s episodes found", category)
	}

	return episodes[rand.Intn(len(episodes))], nil
}

func GetSpecialEpisodes(category data.SpecialCategory) ([]data.ShowSpecialEpisode, error) {
	shows, err := db.ReadAllShows()
	if err != nil {
		return nil, fmt.Errorf("error reading shows: %w", err)
	}

	episodes := make([]data.ShowSpecialEpisode, 0)
	for _, show := range shows {
		for _, special := range show.SpecialEpisodes {
			if special.Category != category {
				continue
			}
			episodes = append(episodes, data.ShowSpecialEpisode{
				ShowName: show.Name,
				Provider: show.Provider,
				Category: special.Category,
				Series:   special.Series,
				Episode:  special.Episode,
			})
		}
	}
	return episodes, nil
}
