package shows

import (
	"reflect"
	"strings"
	"testing"

	"what-to-watch/data"
)

func TestGetCurrentlyWatching(t *testing.T) {
	tests := []struct {
		name     string
		shows    []data.Show
		expected []data.Show
	}{
		{
			name:     "no shows",
			shows:    []data.Show{},
			expected: nil,
		},
		{
			name: "no currently watching shows",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: "Comedy"},
			},
			expected: nil,
		},
		{
			name: "some currently watching shows",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)},
				{Name: "Show B", Genre: "Comedy"},
			},
			expected: []data.Show{
				{Name: "Show A", Genre: "Drama", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2), Series: "1", Episode: "2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCurrentlyWatching(tt.shows)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d shows, got %d", len(tt.expected), len(result))
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetCurrentlyWatchingRejectsInvalidPositions(t *testing.T) {
	tests := []struct {
		name  string
		show  data.Show
		field string
	}{
		{name: "series only", show: data.Show{Name: "Partial", Episodes: []int{2}, CurrentSeries: intPtr(1)}, field: "currentSeries and currentEpisode"},
		{name: "episode only", show: data.Show{Name: "Partial", Episodes: []int{2}, CurrentEpisode: intPtr(1)}, field: "currentSeries and currentEpisode"},
		{name: "series below bounds", show: data.Show{Name: "Bad Series", Episodes: []int{2}, CurrentSeries: intPtr(0), CurrentEpisode: intPtr(1)}, field: "currentSeries"},
		{name: "negative series", show: data.Show{Name: "Bad Series", Episodes: []int{2}, CurrentSeries: intPtr(-1), CurrentEpisode: intPtr(1)}, field: "currentSeries"},
		{name: "negative episode", show: data.Show{Name: "Bad Episode", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(-1)}, field: "currentEpisode"},
		{name: "episode above bounds", show: data.Show{Name: "Bad Episode", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(3)}, field: "currentEpisode"},
		{name: "non-positive episode count", show: data.Show{Name: "Bad Count", Episodes: []int{0}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)}, field: "episodes[1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetCurrentlyWatching([]data.Show{tt.show}); err == nil || !strings.Contains(err.Error(), tt.field) || !strings.Contains(err.Error(), tt.show.Name) {
				t.Fatalf("expected contextual %q validation error for %q, got %v", tt.field, tt.show.Name, err)
			}
		})
	}
}

func TestMarkEpisodeWatched(t *testing.T) {
	tests := []struct {
		name           string
		shows          []data.Show
		listIndex      int
		expected       []data.Show
		expectedFinish bool
		expectError    bool
	}{
		{
			name: "increment episode within series",
			shows: []data.Show{
				{Name: "Show A", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			},
			listIndex: 1,
			expected: []data.Show{
				{Name: "Show A", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)},
			},
			expectedFinish: false,
			expectError:    false,
		},
		{
			name: "rollover to next series",
			shows: []data.Show{
				{Name: "Show B", Episodes: []int{2, 3}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)},
			},
			listIndex: 1,
			expected: []data.Show{
				{Name: "Show B", Episodes: []int{2, 3}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)},
			},
			expectedFinish: false,
			expectError:    false,
		},
		{
			name: "finish show when past last series",
			shows: []data.Show{
				{Name: "Show C", Episodes: []int{1, 1}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)},
			},
			listIndex: 1,
			expected: []data.Show{
				{Name: "Show C", Episodes: []int{1, 1}, CurrentSeries: nil, CurrentEpisode: nil},
			},
			expectedFinish: true,
			expectError:    false,
		},
		{
			name:           "invalid (non-positive) index",
			shows:          []data.Show{},
			listIndex:      0,
			expected:       nil,
			expectedFinish: false,
			expectError:    true,
		},
		{
			name: "index out of range",
			shows: []data.Show{
				{Name: "Show D", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			},
			listIndex:      2,
			expected:       nil,
			expectedFinish: false,
			expectError:    true,
		},
		{
			name: "selected show not marked as watching",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}},
			},
			listIndex:      1,
			expected:       nil,
			expectedFinish: false,
			expectError:    true,
		},
		{
			name: "selected show not marked as watching for series only",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}, CurrentEpisode: intPtr(1)},
			},
			listIndex:      1,
			expected:       nil,
			expectedFinish: false,
			expectError:    true,
		},
		{
			name: "selected show not marked as watching for episode only",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}, CurrentSeries: intPtr(1)},
			},
			listIndex:      1,
			expected:       nil,
			expectedFinish: false,
			expectError:    true,
		},
		{
			name: "invalid position does not complete show",
			shows: []data.Show{
				{Name: "Bad Series", Episodes: []int{3}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)},
			},
			listIndex:   1,
			expectError: true,
		},
		{
			name: "non-positive episode count does not complete show",
			shows: []data.Show{
				{Name: "Bad Count", Episodes: []int{0}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			},
			listIndex:   1,
			expectError: true,
		},
		{
			name: "negative episode count does not complete show",
			shows: []data.Show{
				{Name: "Bad Count", Episodes: []int{-1}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			},
			listIndex:   1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, finish, err := MarkEpisodeWatched(tt.shows, tt.listIndex)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}

			if finish != tt.expectedFinish {
				t.Errorf("expected finish %v, got %v", tt.expectedFinish, finish)
			}
		})
	}
}

func TestGetUniqueGenres(t *testing.T) {
	tests := []struct {
		name     string
		shows    []data.Show
		expected []string
	}{
		{
			name:     "no shows",
			shows:    []data.Show{},
			expected: []string(nil),
		},
		{
			name: "single genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: "Drama"},
			},
			expected: []string{"Drama"},
		},
		{
			name: "multiple genres",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: "Comedy"},
				{Name: "Show C", Genre: "Sci-Fi"},
				{Name: "Show D", Genre: "Drama"},
			},
			expected: []string{"Drama", "Comedy", "Sci-Fi"},
		},
		{
			name: "empty genre is ignored",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: ""},
			},
			expected: []string{"Drama"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUniqueGenres(tt.shows)
			// Note: order of genres is not guaranteed, so we use a map-based comparison
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d genres, got %d", len(tt.expected), len(result))
			}
			resultMap := make(map[string]bool)
			for _, g := range result {
				resultMap[g] = true
			}
			for _, g := range tt.expected {
				if !resultMap[g] {
					t.Errorf("expected genre %s not found in result", g)
				}
			}
		})
	}
}

func TestGetUnwatchedShowsByGenre(t *testing.T) {
	tests := []struct {
		name     string
		shows    []data.Show
		genre    string
		expected []data.Show
	}{
		{
			name:     "no shows",
			shows:    []data.Show{},
			genre:    "Drama",
			expected: []data.Show(nil),
		},
		{
			name: "no unwatched shows for genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama", CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
				{Name: "Show B", Genre: "Comedy"},
			},
			genre:    "Drama",
			expected: []data.Show(nil),
		},
		{
			name: "get unwatched shows for genre",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: "Comedy"},
				{Name: "Show C", Genre: "Drama", CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
				{Name: "Show D", Genre: "Drama"},
			},
			genre: "Drama",
			expected: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show D", Genre: "Drama"},
			},
		},
		{
			name: "genre does not exist",
			shows: []data.Show{
				{Name: "Show A", Genre: "Drama"},
				{Name: "Show B", Genre: "Comedy"},
			},
			genre:    "Horror",
			expected: []data.Show(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUnwatchedShowsByGenre(tt.shows, tt.genre)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetSpecialEpisodes(t *testing.T) {
	shows := []data.Show{
		{
			Name: "The Office", Provider: "Netflix",
			SpecialEpisodes: []data.SpecialEpisode{
				{Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
				{Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4},
			},
		},
		{
			Name: "The Office", Provider: "Prime",
			SpecialEpisodes: []data.SpecialEpisode{
				{Category: data.SpecialCategoryHalloween, Series: 1, Episode: 6},
			},
		},
	}

	got := GetSpecialEpisodes(shows, data.SpecialCategoryHalloween)
	want := []data.ShowSpecialEpisode{
		{ShowName: "The Office", Provider: "Netflix", Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
		{ShowName: "The Office", Provider: "Prime", Category: data.SpecialCategoryHalloween, Series: 1, Episode: 6},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %+v, got %+v", want, got)
	}

	christmas := GetSpecialEpisodes(shows, data.SpecialCategoryChristmas)
	wantChristmas := []data.ShowSpecialEpisode{{
		ShowName: "The Office", Provider: "Netflix", Category: data.SpecialCategoryChristmas, Series: 3, Episode: 4,
	}}
	if !reflect.DeepEqual(christmas, wantChristmas) {
		t.Fatalf("expected %+v, got %+v", wantChristmas, christmas)
	}
}

// intPtr is a small test helper to construct *int values inline.
func intPtr(i int) *int {
	return &i
}
