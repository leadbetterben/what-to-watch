package data

import "testing"

func intPtr(value int) *int { return &value }

func TestValidateEpisodeCounts(t *testing.T) {
	tests := []struct {
		name    string
		show    Show
		wantErr bool
	}{
		{name: "positive counts", show: Show{Name: "Show", Episodes: []int{10, 8}}},
		{name: "missing counts", show: Show{Name: "Show"}, wantErr: true},
		{name: "zero count", show: Show{Name: "Show", Episodes: []int{10, 0}}, wantErr: true},
		{name: "negative count", show: Show{Name: "Show", Episodes: []int{-1}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEpisodeCounts(tt.show)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateEpisodeCounts() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCurrentPosition(t *testing.T) {
	tests := []struct {
		name    string
		show    Show
		wantErr bool
	}{
		{
			name: "no current position",
			show: Show{Name: "Show", Episodes: []int{10}},
		},
		{
			name: "valid current position",
			show: Show{Name: "Show", Episodes: []int{10, 8}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(8)},
		},
		{
			name:    "missing current series",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentEpisode: intPtr(1)},
			wantErr: true,
		},
		{
			name:    "missing current episode",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentSeries: intPtr(1)},
			wantErr: true,
		},
		{
			name:    "current series must be positive",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentSeries: intPtr(0), CurrentEpisode: intPtr(1)},
			wantErr: true,
		},
		{
			name:    "current series inside bounds",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentSeries: intPtr(2), CurrentEpisode: intPtr(1)},
			wantErr: true,
		},
		{
			name:    "current episode must be positive",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(0)},
			wantErr: true,
		},
		{
			name:    "all episode counts must be positive",
			show:    Show{Name: "Show", Episodes: []int{10, 0}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(1)},
			wantErr: true,
		},
		{
			name:    "current episode inside series bounds",
			show:    Show{Name: "Show", Episodes: []int{10}, CurrentSeries: intPtr(1), CurrentEpisode: intPtr(11)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrentPosition(tt.show)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCurrentPosition() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSpecialEpisodes(t *testing.T) {
	tests := []struct {
		name    string
		show    Show
		wantErr bool
	}{
		{
			name: "no special episodes",
			show: Show{Name: "Show"},
		},
		{
			name: "supported categories and in-bounds episodes",
			show: Show{Name: "Show", Episodes: []int{10, 8}, SpecialEpisodes: []SpecialEpisode{
				{Category: SpecialCategoryHalloween, Series: 1, Episode: 10},
				{Category: SpecialCategoryChristmas, Series: 2, Episode: 8},
			}},
		},
		{
			name:    "series outside bounds when episode counts are missing",
			show:    Show{Name: "Show", SpecialEpisodes: []SpecialEpisode{{Category: SpecialCategoryHalloween, Series: 3, Episode: 12}}},
			wantErr: true,
		},
		{
			name:    "invalid category",
			show:    Show{Name: "Show", SpecialEpisodes: []SpecialEpisode{{Category: "other", Series: 1, Episode: 1}}},
			wantErr: true,
		},
		{
			name:    "series must be positive",
			show:    Show{Name: "Show", SpecialEpisodes: []SpecialEpisode{{Category: SpecialCategoryHalloween, Series: 0, Episode: 1}}},
			wantErr: true,
		},
		{
			name:    "episode must be positive",
			show:    Show{Name: "Show", SpecialEpisodes: []SpecialEpisode{{Category: SpecialCategoryHalloween, Series: 1, Episode: -1}}},
			wantErr: true,
		},
		{
			name:    "series outside bounds",
			show:    Show{Name: "Show", Episodes: []int{10}, SpecialEpisodes: []SpecialEpisode{{Category: SpecialCategoryHalloween, Series: 2, Episode: 1}}},
			wantErr: true,
		},
		{
			name:    "episode outside bounds",
			show:    Show{Name: "Show", Episodes: []int{10}, SpecialEpisodes: []SpecialEpisode{{Category: SpecialCategoryHalloween, Series: 1, Episode: 11}}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSpecialEpisodes(tt.show)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSpecialEpisodes() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
