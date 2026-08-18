package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"what-to-watch/data"
	"what-to-watch/db"
)

func TestGetCurrentlyWatchingShows(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "shows.json", `[
  {
    "name": "Slow Horses",
    "genre": "Drama",
    "episodes": [6],
    "provider": "Apple TV+",
    "currentSeries": 1,
    "currentEpisode": 5
  },
  {
    "name": "Unstarted Comedy",
    "genre": "Comedy",
    "episodes": [8],
    "provider": "Netflix"
  }
]`)

	shows, err := GetCurrentlyWatchingShows()
	if err != nil {
		t.Fatalf("GetCurrentlyWatchingShows returned error: %v", err)
	}

	expected := []data.Show{
		{
			Name:           "Slow Horses",
			Genre:          "Drama",
			Episodes:       []int{6},
			Provider:       "Apple TV+",
			CurrentSeries:  intPtr(1),
			CurrentEpisode: intPtr(5),
			Series:         "1",
			Episode:        "5",
		},
	}
	if !reflect.DeepEqual(shows, expected) {
		t.Fatalf("expected %+v, got %+v", expected, shows)
	}
}

func TestMarkShowWatchedIncrementsAndPersistsEpisode(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "shows.json", `[
  {
    "name": "Severance",
    "genre": "Sci-Fi",
    "episodes": [9],
    "provider": "Apple TV+",
    "currentSeries": 1,
    "currentEpisode": 3
  }
]`)

	completed, err := MarkShowWatched(1)
	if err != nil {
		t.Fatalf("MarkShowWatched returned error: %v", err)
	}
	if completed {
		t.Fatalf("expected show not to be completed")
	}

	shows := readCurrentShowsFixture(t, dir)
	if len(shows) != 1 {
		t.Fatalf("expected 1 persisted show, got %d", len(shows))
	}
	if shows[0].CurrentEpisode == nil || *shows[0].CurrentEpisode != 4 {
		t.Fatalf("expected persisted episode 4, got %+v", shows[0])
	}
	if shows[0].CurrentSeries == nil || *shows[0].CurrentSeries != 1 {
		t.Fatalf("expected persisted series 1, got %+v", shows[0])
	}
}

func TestMarkShowWatchedCompletesAndPersistsShow(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "shows.json", `[
  {
    "name": "One Episode Show",
    "genre": "Drama",
    "episodes": [1],
    "provider": "BBC",
    "currentSeries": 1,
    "currentEpisode": 1
  }
]`)

	completed, err := MarkShowWatched(1)
	if err != nil {
		t.Fatalf("MarkShowWatched returned error: %v", err)
	}
	if !completed {
		t.Fatalf("expected show to be completed")
	}

	shows := readCurrentShowsFixture(t, dir)
	if len(shows) != 1 {
		t.Fatalf("expected 1 persisted show, got %d", len(shows))
	}
	if shows[0].CurrentSeries != nil || shows[0].CurrentEpisode != nil {
		t.Fatalf("expected persisted show to clear current position, got %+v", shows[0])
	}
}

func TestGetAllFilms(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "films.json", `[
  {
    "name": "Arrival",
    "genre": "Sci-Fi",
    "provider": "Netflix"
  },
  {
    "name": "Paddington 2",
    "genre": "Comedy",
    "provider": "BBC iPlayer"
  }
]`)

	films, err := GetAllFilms()
	if err != nil {
		t.Fatalf("GetAllFilms returned error: %v", err)
	}

	expected := []data.Film{
		{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"},
		{Name: "Paddington 2", Genre: "Comedy", Provider: "BBC iPlayer"},
	}
	if !reflect.DeepEqual(films, expected) {
		t.Fatalf("expected %+v, got %+v", expected, films)
	}
}

func TestGetAvailableGenres(t *testing.T) {
	dir := useTempDataDir(t)
	writeShowsFixture(t, dir)

	genres, err := GetAvailableGenres()
	if err != nil {
		t.Fatalf("GetAvailableGenres returned error: %v", err)
	}

	expected := map[string]bool{"Drama": true, "Comedy": true, "Sci-Fi": true}
	if len(genres) != len(expected) {
		t.Fatalf("expected %d genres, got %d: %v", len(expected), len(genres), genres)
	}
	for _, genre := range genres {
		if !expected[genre] {
			t.Fatalf("unexpected genre %q in %v", genre, genres)
		}
	}
}

func TestGetUnwatchedShowsByGenre(t *testing.T) {
	dir := useTempDataDir(t)
	writeShowsFixture(t, dir)

	shows, err := GetUnwatchedShowsByGenre("Drama")
	if err != nil {
		t.Fatalf("GetUnwatchedShowsByGenre returned error: %v", err)
	}

	expected := []data.Show{
		{Name: "Unstarted Drama", Genre: "Drama", Episodes: []int{6}, Provider: "Netflix"},
	}
	if !reflect.DeepEqual(shows, expected) {
		t.Fatalf("expected %+v, got %+v", expected, shows)
	}
}

func TestHandlersSurfaceReadErrors(t *testing.T) {
	useTempDataDir(t)

	tests := []struct {
		name      string
		call      func() error
		wantError string
	}{
		{
			name: "GetCurrentlyWatchingShows",
			call: func() error {
				_, err := GetCurrentlyWatchingShows()
				return err
			},
			wantError: "error reading shows",
		},
		{
			name: "MarkShowWatched",
			call: func() error {
				_, err := MarkShowWatched(1)
				return err
			},
			wantError: "error reading shows",
		},
		{
			name: "GetAllFilms",
			call: func() error {
				_, err := GetAllFilms()
				return err
			},
			wantError: "error reading films",
		},
		{
			name: "GetAvailableGenres",
			call: func() error {
				_, err := GetAvailableGenres()
				return err
			},
			wantError: "GetAvailableGenres: error reading shows",
		},
		{
			name: "GetUnwatchedShowsByGenre",
			call: func() error {
				_, err := GetUnwatchedShowsByGenre("Drama")
				return err
			},
			wantError: "GetUnwatchedShowsByGenre: error reading shows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %q", tt.wantError, err.Error())
			}
		})
	}
}

func useTempDataDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Cleanup(db.SetFullPathResolverForTest(func(path string) string {
		return filepath.Join(dir, path)
	}))

	return dir
}

func writeShowsFixture(t *testing.T, dir string) {
	t.Helper()

	writeFixture(t, dir, "shows.json", `[
  {
    "name": "Unstarted Drama",
    "genre": "Drama",
    "episodes": [6],
    "provider": "Netflix"
  },
  {
    "name": "Current Drama",
    "genre": "Drama",
    "episodes": [8],
    "provider": "Apple TV+",
    "currentSeries": 1,
    "currentEpisode": 2
  },
  {
    "name": "Unstarted Comedy",
    "genre": "Comedy",
    "episodes": [10],
    "provider": "Disney+"
  },
  {
    "name": "Unstarted Sci-Fi",
    "genre": "Sci-Fi",
    "episodes": [9],
    "provider": "Prime"
  },
  {
    "name": "Missing Genre",
    "episodes": [4],
    "provider": "Channel 4"
  }
]`)
}

func writeFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writing fixture %s: %v", name, err)
	}
}

func readCurrentShowsFixture(t *testing.T, dir string) []data.Show {
	t.Helper()

	raw, err := os.ReadFile(filepath.Join(dir, "shows.json"))
	if err != nil {
		t.Fatalf("reading current shows fixture: %v", err)
	}

	var shows []data.Show
	if err := json.Unmarshal(raw, &shows); err != nil {
		t.Fatalf("unmarshaling current shows fixture: %v", err)
	}

	return shows
}

func intPtr(i int) *int {
	return &i
}
