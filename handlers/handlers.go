package handlers

import (
	"fmt"

	"what-to-watch/data"
	"what-to-watch/db"
	"what-to-watch/shows"
)

// GetCurrentlyWatchingShows retrieves the list of currently watching shows
func GetCurrentlyWatchingShows() ([]data.Show, error) {
	s, err := db.ReadCurrentShows()
	if err != nil {
		return nil, fmt.Errorf("error reading shows: %w", err)
	}

	cw, err := shows.GetCurrentlyWatching(s)
	if err != nil {
		return nil, fmt.Errorf("error getting currently watching shows: %w", err)
	}

	return cw, nil
}

// MarkShowWatched marks an episode as watched and updates the data store
// idx is 1-based index from the currently watching list
func MarkShowWatched(idx int) (bool, error) {
	s, err := db.ReadCurrentShows()
	if err != nil {
		return false, fmt.Errorf("error reading shows: %w", err)
	}

	updatedShows, isCompleted, err := shows.MarkEpisodeWatched(s, idx)
	if err != nil {
		return false, fmt.Errorf("error updating show: %w", err)
	}

	if err := db.WriteCurrentShows(updatedShows); err != nil {
		return false, fmt.Errorf("error saving updated shows: %w", err)
	}

	return isCompleted, nil
}

// GetAllFilms retrieves the list of all films
func GetAllFilms() ([]data.Film, error) {
	films, err := db.ReadFilms()
	if err != nil {
		return nil, fmt.Errorf("error reading films: %w", err)
	}

	return films, nil
}
