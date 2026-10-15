package handlers

import (
	"fmt"
	"math/rand"

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
	completed, err := db.UpdateCurrentShows(func(current []data.Show) ([]data.Show, bool, error) {
		updatedShows, isCompleted, err := shows.MarkEpisodeWatched(current, idx)
		if err != nil {
			return nil, false, fmt.Errorf("error updating show: %w", err)
		}
		return updatedShows, isCompleted, nil
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
	s, err := db.ReadAllShows()
	if err != nil {
		return nil, fmt.Errorf("GetAvailableGenres: error reading shows: %w", err)
	}

	return shows.GetUniqueGenres(s), nil
}

// GetUnwatchedShowsByGenre retrieves all unwatched shows for a given genre
func GetUnwatchedShowsByGenre(genre string) ([]data.Show, error) {
	s, err := db.ReadUnwatchedShows()
	if err != nil {
		return nil, fmt.Errorf("GetUnwatchedShowsByGenre: error reading shows: %w", err)
	}

	return shows.GetUnwatchedShowsByGenre(s, genre), nil
}

// GetHalloweenEpisodes retrieves all curated Halloween episodes.
func GetHalloweenEpisodes() ([]data.ShowSpecialEpisode, error) {
	return getSpecialEpisodes(data.SpecialCategoryHalloween)
}

// GetChristmasEpisodes retrieves all curated Christmas episodes.
func GetChristmasEpisodes() ([]data.ShowSpecialEpisode, error) {
	return getSpecialEpisodes(data.SpecialCategoryChristmas)
}

// GetRandomHalloweenEpisode retrieves one curated Halloween episode at random.
func GetRandomHalloweenEpisode() (data.ShowSpecialEpisode, error) {
	return getRandomSpecialEpisode(data.SpecialCategoryHalloween)
}

// GetRandomChristmasEpisode retrieves one curated Christmas episode at random.
func GetRandomChristmasEpisode() (data.ShowSpecialEpisode, error) {
	return getRandomSpecialEpisode(data.SpecialCategoryChristmas)
}

func getRandomSpecialEpisode(category data.SpecialCategory) (data.ShowSpecialEpisode, error) {
	episodes, err := getSpecialEpisodes(category)
	if err != nil {
		return data.ShowSpecialEpisode{}, err
	}
	if len(episodes) == 0 {
		return data.ShowSpecialEpisode{}, fmt.Errorf("no %s episodes found", category)
	}

	return episodes[rand.Intn(len(episodes))], nil
}

func getSpecialEpisodes(category data.SpecialCategory) ([]data.ShowSpecialEpisode, error) {
	allShows, err := db.ReadAllShows()
	if err != nil {
		return nil, fmt.Errorf("error reading shows: %w", err)
	}

	return shows.GetSpecialEpisodes(allShows, category), nil
}
