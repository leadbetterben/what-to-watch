package shows

import (
	"reflect"
	"testing"

	"what-to-watch/data"
)

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

func TestGetSpecialEpisodes(t *testing.T) {
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
