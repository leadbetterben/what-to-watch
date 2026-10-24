package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"what-to-watch/data"
	"what-to-watch/db"
)

func TestGetCurrentlyWatchingShows(t *testing.T) {
	tests := []struct {
		name          string
		shows         []data.Show
		want          []data.Show
		wantErrorPart string
	}{
		{
			name:  "no shows",
			shows: []data.Show{},
			want:  nil,
		},
		{
			name: "no currently watching shows",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{10}},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{8}},
			},
			want: nil,
		},
		{
			name: "some currently watching shows",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{8}},
			},
			want: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2), Series: "1", Episode: "2"},
			},
		},
		{
			name:          "series only",
			shows:         []data.Show{{Name: "Partial", Episodes: []int{2}, CurrentSeries: intPtr(1)}},
			wantErrorPart: "currentSeries and currentEpisode",
		},
		{
			name:          "episode only",
			shows:         []data.Show{{Name: "Partial", Episodes: []int{2}, CurrentEpisode: intPtr(1)}},
			wantErrorPart: "currentSeries and currentEpisode",
		},
		{
			name:          "series below bounds",
			shows:         []data.Show{{Name: "Bad Series", Episodes: []int{2}, CurrentSeries: intPtr(0), CurrentEpisode: intPtr(1)}},
			wantErrorPart: "currentSeries",
		},
		{
			name:          "negative series",
			shows:         []data.Show{{Name: "Bad Series", Episodes: []int{2}, CurrentSeries: intPtr(-1), CurrentEpisode: intPtr(1)}},
			wantErrorPart: "currentSeries",
		},
		{
			name:          "negative episode",
			shows:         []data.Show{{Name: "Bad Episode", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(-1)}},
			wantErrorPart: "currentEpisode",
		},
		{
			name:          "episode above bounds",
			shows:         []data.Show{{Name: "Bad Episode", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(3)}},
			wantErrorPart: "currentEpisode",
		},
		{
			name:          "non-positive episode count",
			shows:         []data.Show{{Name: "Bad Count", Episodes: []int{0}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}},
			wantErrorPart: "episodes[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture, err := json.Marshal(tt.shows)
			if err != nil {
				t.Fatalf("marshal shows fixture: %v", err)
			}
			writeFixture(t, dir, "shows.json", string(fixture))

			got, err := GetCurrentlyWatchingShows()
			if tt.wantErrorPart != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrorPart) || !strings.Contains(err.Error(), tt.shows[0].Name) {
					t.Fatalf("expected contextual %q validation error for %q, got %v", tt.wantErrorPart, tt.shows[0].Name, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetCurrentlyWatchingShows returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
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

func TestMarkShowWatchedSerializesConcurrentUpdates(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "shows.json", `[
  {"name":"Current","episodes":[100],"provider":"Netflix","currentSeries":1,"currentEpisode":1},
  {"name":"Unwatched","episodes":[6],"provider":"Prime"}
]`)

	const updates = 24
	start := make(chan struct{})
	errs := make(chan error, updates)
	var wg sync.WaitGroup
	for range updates {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, err := MarkShowWatched(1); err != nil {
				errs <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("MarkShowWatched returned error: %v", err)
	}

	all := readCurrentShowsFixture(t, dir)
	if len(all) != 2 {
		t.Fatalf("expected current and unwatched shows to remain, got %+v", all)
	}
	if all[0].Name != "Unwatched" {
		t.Fatalf("expected unwatched record to be preserved, got %+v", all[0])
	}
	if all[1].Name != "Current" || all[1].CurrentEpisode == nil || *all[1].CurrentEpisode != updates+1 {
		t.Fatalf("expected all %d updates to be retained, got %+v", updates, all[1])
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
	tests := []struct {
		name  string
		shows []data.Show
		want  []string
	}{
		{
			name:  "no shows",
			shows: []data.Show{},
			want:  nil,
		},
		{
			name: "single genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{1}},
				{Name: "Show B", Genre: "Drama", Episodes: []int{1}},
			},
			want: []string{"Drama"},
		},
		{
			name: "multiple genres",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{1}},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{1}},
				{Name: "Show C", Genre: "Sci-Fi", Episodes: []int{1}},
				{Name: "Show D", Genre: "Drama", Episodes: []int{1}},
			},
			want: []string{"Drama", "Comedy", "Sci-Fi"},
		},
		{
			name: "empty genre is ignored",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{1}},
				{Name: "Show B", Episodes: []int{1}},
			},
			want: []string{"Drama"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture, err := json.Marshal(tt.shows)
			if err != nil {
				t.Fatalf("marshal shows fixture: %v", err)
			}
			writeFixture(t, dir, "shows.json", string(fixture))

			got, err := GetAvailableGenres()
			if err != nil {
				t.Fatalf("GetAvailableGenres returned error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d genres, got %d: %v", len(tt.want), len(got), got)
			}
			gotSet := make(map[string]bool, len(got))
			for _, genre := range got {
				gotSet[genre] = true
			}
			for _, genre := range tt.want {
				if !gotSet[genre] {
					t.Errorf("expected genre %q in %v", genre, got)
				}
			}
		})
	}
}

func TestGetUnwatchedShowsByGenre(t *testing.T) {
	tests := []struct {
		name  string
		shows []data.Show
		genre string
		want  []data.Show
	}{
		{
			name:  "no shows",
			shows: []data.Show{},
			genre: "Drama",
			want:  nil,
		},
		{
			name: "no unwatched shows for genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{1}},
			},
			genre: "Drama",
			want:  nil,
		},
		{
			name: "get unwatched shows for genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{2}},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{1}},
				{Name: "Show C", Genre: "Drama", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
				{Name: "Show D", Genre: "Drama", Episodes: []int{3}},
			},
			genre: "Drama",
			want: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{2}},
				{Name: "Show D", Genre: "Drama", Episodes: []int{3}},
			},
		},
		{
			name: "genre does not exist",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{1}},
				{Name: "Show B", Genre: "Comedy", Episodes: []int{1}},
			},
			genre: "Horror",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture, err := json.Marshal(tt.shows)
			if err != nil {
				t.Fatalf("marshal shows fixture: %v", err)
			}
			writeFixture(t, dir, "shows.json", string(fixture))

			got, err := GetUnwatchedShowsByGenre(tt.genre)
			if err != nil {
				t.Fatalf("GetUnwatchedShowsByGenre returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestGetSeasonalSpecialEpisodes(t *testing.T) {
	shows := []data.Show{
		{
			Name: "The Office", Provider: "Netflix", Episodes: []int{4, 5, 4},
			SpecialEpisodes: []data.SpecialEpisode{
				{Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
				{Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4},
			},
		},
		{
			Name: "The Office", Provider: "Prime", Episodes: []int{6},
			SpecialEpisodes: []data.SpecialEpisode{
				{Category: data.SpecialCategoryHalloween, Series: 1, Episode: 6},
			},
		},
	}
	tests := []struct {
		name string
		call func() ([]data.ShowSpecialEpisode, error)
		want []data.ShowSpecialEpisode
	}{
		{
			name: "Halloween",
			call: func() ([]data.ShowSpecialEpisode, error) {
				return GetSpecialEpisodes(data.SpecialCategoryHalloween)
			},
			want: []data.ShowSpecialEpisode{
				{ShowName: "The Office", Provider: "Netflix", Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
				{ShowName: "The Office", Provider: "Prime", Category: data.SpecialCategoryHalloween, Series: 1, Episode: 6},
			},
		},
		{
			name: "Christmas",
			call: func() ([]data.ShowSpecialEpisode, error) {
				return GetSpecialEpisodes(data.SpecialCategoryChristmas)
			},
			want: []data.ShowSpecialEpisode{
				{ShowName: "The Office", Provider: "Netflix", Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture, err := json.Marshal(shows)
			if err != nil {
				t.Fatalf("marshal shows fixture: %v", err)
			}
			writeFixture(t, dir, "shows.json", string(fixture))

			got, err := tt.call()
			if err != nil {
				t.Fatalf("seasonal special episode handler returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestGetRandomSeasonalSpecialEpisode(t *testing.T) {
	GetRandomHalloweenEpisode := func() (data.ShowSpecialEpisode, error) {
		return GetRandomSpecialEpisode(data.SpecialCategoryHalloween)
	}
	GetRandomChristmasEpisode := func() (data.ShowSpecialEpisode, error) {
		return GetRandomSpecialEpisode(data.SpecialCategoryChristmas)
	}

	tests := []struct {
		name          string
		shows         []data.Show
		call          func() (data.ShowSpecialEpisode, error)
		want          data.ShowSpecialEpisode
		wantErrorPart string
	}{
		{
			name: "Halloween",
			shows: []data.Show{{
				Name: "Example Show", Provider: "Netflix", Episodes: []int{4, 5, 4},
				SpecialEpisodes: []data.SpecialEpisode{
					{Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
					{Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4},
				},
			}},
			call: GetRandomHalloweenEpisode,
			want: data.ShowSpecialEpisode{
				ShowName: "Example Show", Provider: "Netflix", Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5,
			},
		},
		{
			name: "Christmas",
			shows: []data.Show{{
				Name: "Example Show", Provider: "Netflix", Episodes: []int{4, 5, 4},
				SpecialEpisodes: []data.SpecialEpisode{
					{Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
					{Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4},
				},
			}},
			call: GetRandomChristmasEpisode,
			want: data.ShowSpecialEpisode{
				ShowName: "Example Show", Provider: "Netflix", Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4,
			},
		},
		{
			name:          "empty Halloween episodes",
			shows:         []data.Show{},
			call:          GetRandomHalloweenEpisode,
			wantErrorPart: "no halloween episodes found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture, err := json.Marshal(tt.shows)
			if err != nil {
				t.Fatalf("marshal shows fixture: %v", err)
			}
			writeFixture(t, dir, "shows.json", string(fixture))

			got, err := tt.call()
			if tt.wantErrorPart != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrorPart) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrorPart, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("random episode handler returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
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
