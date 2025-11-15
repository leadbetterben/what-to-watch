package shows

import (
	"reflect"
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
				{Name: "Show A", Genre: "Drama", CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2)},
				{Name: "Show B", Genre: "Comedy"},
				{Name: "Show C", Genre: "Sci-Fi", CurrentSeries: intPtr(2)},
				{Name: "Show D", Genre: "Horror", CurrentEpisode: intPtr(3)},
			},
			expected: []data.Show{
				{Name: "Show A", Genre: "Drama", CurrentSeries: intPtr(1), CurrentEpisode: intPtr(2), Series: "1", Episode: "2"},
				{Name: "Show C", Genre: "Sci-Fi", CurrentSeries: intPtr(2), CurrentEpisode: nil, Series: "2", Episode: "-"},
				{Name: "Show D", Genre: "Horror", CurrentSeries: nil, CurrentEpisode: intPtr(3), Series: "-", Episode: "3"},
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

func TestMarkEpisodeWatched(t *testing.T) {
	tests := []struct {
		name        string
		shows       []data.Show
		listIndex   int
		expected    []data.Show
		expectedMsg string
		expectError bool
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
			expectedMsg: "Updated show progress.",
			expectError: false,
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
			expectedMsg: "Updated show progress.",
			expectError: false,
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
			expectedMsg: "Congratulations! You finished Show C.",
			expectError: false,
		},
		{
			name:        "invalid (non-positive) index",
			shows:       []data.Show{},
			listIndex:   0,
			expected:    nil,
			expectedMsg: "",
			expectError: true,
		},
		{
			name: "index out of range",
			shows: []data.Show{
				{Name: "Show D", Episodes: []int{2}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			},
			listIndex:   2,
			expected:    nil,
			expectedMsg: "",
			expectError: true,
		},
		{
			name: "selected show not marked as watching",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}},
			},
			listIndex:   1,
			expected:    nil,
			expectedMsg: "",
			expectError: true,
		},
		{
			name: "selected show not marked as watching for series only",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}, CurrentEpisode: intPtr(1)},
			},
			listIndex:   1,
			expected:    nil,
			expectedMsg: "",
			expectError: true,
		},
		{
			name: "selected show not marked as watching for episode only",
			shows: []data.Show{
				{Name: "Show E", Episodes: []int{3}, CurrentSeries: intPtr(1)},
			},
			listIndex:   1,
			expected:    nil,
			expectedMsg: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, msg, err := MarkEpisodeWatched(tt.shows, tt.listIndex)
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

			if msg != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}

// intPtr is a small test helper to construct *int values inline.
func intPtr(i int) *int {
	return &i
}
