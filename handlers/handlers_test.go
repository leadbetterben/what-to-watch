package handlers

import (
	"reflect"
	"strings"
	"testing"

	"what-to-watch/data"
)

func TestGetCurrentlyWatchingShows(t *testing.T) {
	tests := []struct {
		name  string
		shows []data.Show
		want  []data.Show
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubStorage(t, tt.shows, nil)
			got, err := GetCurrentlyWatchingShows()
			if err != nil {
				t.Fatalf("GetCurrentlyWatchingShows returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestMarkShowWatched(t *testing.T) {
	tests := []struct {
		name          string
		shows         []data.Show
		index         int
		wantShows     []data.Show
		wantCompleted bool
		wantError     bool
	}{
		{
			name:      "increment episode within series",
			shows:     []data.Show{{Name: "Show A", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}},
			index:     1,
			wantShows: []data.Show{{Name: "Show A", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)}},
		},
		{
			name:      "rollover to next series",
			shows:     []data.Show{{Name: "Show B", Episodes: []int{2, 3}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)}},
			index:     1,
			wantShows: []data.Show{{Name: "Show B", Episodes: []int{2, 3}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)}},
		},
		{
			name:          "finish show when past last series",
			shows:         []data.Show{{Name: "Show C", Episodes: []int{1, 1}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)}},
			index:         1,
			wantShows:     []data.Show{{Name: "Show C", Episodes: []int{1, 1}}},
			wantCompleted: true,
		},
		{
			name:      "invalid non-positive index",
			shows:     []data.Show{},
			index:     0,
			wantShows: []data.Show{},
			wantError: true,
		},
		{
			name:      "index out of range",
			shows:     []data.Show{{Name: "Show D", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}},
			index:     2,
			wantError: true,
		},
		{
			name:      "selected show not marked as watching",
			shows:     []data.Show{{Name: "Show E", Episodes: []int{3}}},
			index:     1,
			wantError: true,
		},
		{
			name:      "selected show not marked as watching for series only",
			shows:     []data.Show{{Name: "Show E", Episodes: []int{3}, CurrentEpisode: intPtr(1)}},
			index:     1,
			wantError: true,
		},
		{
			name:      "selected show not marked as watching for episode only",
			shows:     []data.Show{{Name: "Show E", Episodes: []int{3}, CurrentSeries: intPtr(1)}},
			index:     1,
			wantError: true,
		},
		{
			name:      "invalid position does not complete show",
			shows:     []data.Show{{Name: "Bad Series", Episodes: []int{3}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)}},
			index:     1,
			wantError: true,
		},
		{
			name:      "non-positive episode count does not complete show",
			shows:     []data.Show{{Name: "Bad Count", Episodes: []int{0}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}},
			index:     1,
			wantError: true,
		},
		{
			name:      "negative episode count does not complete show",
			shows:     []data.Show{{Name: "Bad Count", Episodes: []int{-1}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}},
			index:     1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stored := append([]data.Show(nil), tt.shows...)
			stubStorage(t, nil, nil)
			updateCurrentShows = func(update func([]data.Show) ([]data.Show, bool, error)) (bool, error) {
				updated, completed, err := update(stored)
				if err == nil {
					stored = updated
				}
				return completed, err
			}
			completed, err := MarkShowWatched(tt.index)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("MarkShowWatched returned error: %v", err)
			}
			if completed != tt.wantCompleted {
				t.Fatalf("expected completed=%v, got %v", tt.wantCompleted, completed)
			}

			gotShows := stored
			if !reflect.DeepEqual(gotShows, tt.wantShows) {
				t.Fatalf("expected persisted shows %+v, got %+v", tt.wantShows, gotShows)
			}
		})
	}
}

func TestGetAllFilms(t *testing.T) {
	expected := []data.Film{
		{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"},
		{Name: "Paddington 2", Genre: "Comedy", Provider: "BBC iPlayer"},
	}
	stubStorage(t, nil, expected)
	films, err := GetAllFilms()
	if err != nil {
		t.Fatalf("GetAllFilms returned error: %v", err)
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
			stubStorage(t, tt.shows, nil)
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
			stubStorage(t, tt.shows, nil)
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
			stubStorage(t, shows, nil)
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
			stubStorage(t, tt.shows, nil)
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

func stubStorage(t *testing.T, shows []data.Show, films []data.Film) {
	t.Helper()
	oldReadShows, oldReadFilms, oldUpdate := readShows, readFilms, updateCurrentShows
	readShows = func() ([]data.Show, error) { return shows, nil }
	readFilms = func() ([]data.Film, error) { return films, nil }
	updateCurrentShows = func(update func([]data.Show) ([]data.Show, bool, error)) (bool, error) {
		_, completed, err := update(shows)
		return completed, err
	}
	t.Cleanup(func() { readShows, readFilms, updateCurrentShows = oldReadShows, oldReadFilms, oldUpdate })
}

func intPtr(i int) *int {
	return &i
}
