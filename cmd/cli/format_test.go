package cli

import (
	"testing"

	"what-to-watch/data"
)

func TestFormatShowsTable(t *testing.T) {
	tests := []struct {
		name  string
		shows []data.Show
		want  string
	}{
		{
			name: "normal case",
			shows: []data.Show{
				{Name: "A", Genre: "Drama", Provider: "BBC", Series: "1", Episode: "2"},
				{Name: "Long Show", Genre: "Comedy", Provider: "Netflix", Series: "12", Episode: "10"},
			},
			want: "" +
				"Index  Name       Genre   Provider  Series  Episode\n" +
				"-----  ---------  ------  --------  ------  -------\n" +
				"1      A          Drama   BBC       1       2      \n" +
				"2      Long Show  Comedy  Netflix   12      10     \n",
		},
		{
			name:  "empty shows",
			shows: []data.Show{},
			want:  "No shows currently being watched.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatShowsTable(tt.shows); got != tt.want {
				t.Fatalf("formatShowsTable() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatShowsByGenreTable(t *testing.T) {
	tests := []struct {
		name  string
		shows []data.Show
		want  string
	}{
		{
			name: "normal case",
			shows: []data.Show{
				{Name: "A", Provider: "BBC"},
				{Name: "Long Show", Provider: "Netflix"},
			},
			want: "" +
				"Index  Name       Provider\n" +
				"-----  ---------  --------\n" +
				"1      A          BBC     \n" +
				"2      Long Show  Netflix \n",
		},
		{
			name:  "empty shows",
			shows: []data.Show{},
			want:  "No unwatched shows in this genre.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatShowsByGenreTable(tt.shows); got != tt.want {
				t.Fatalf("formatShowsByGenreTable() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFilmsTable(t *testing.T) {
	test := []struct {
		name  string
		films []data.Film
		want  string
	}{
		{
			name: "normal case",
			films: []data.Film{
				{Name: "A", Genre: "Drama", Provider: "BBC"},
				{Name: "Long Film", Genre: "Comedy", Provider: "Netflix"},
			},
			want: "" +
				"Index  Name       Genre   Provider\n" +
				"-----  ---------  ------  --------\n" +
				"1      A          Drama   BBC     \n" +
				"2      Long Film  Comedy  Netflix \n",
		},
		{
			name:  "empty films",
			films: []data.Film{},
			want:  "No films found.\n",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatFilmsTable(tt.films); got != tt.want {
				t.Fatalf("formatFilmsTable() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatSpecialEpisodesTable(t *testing.T) {
	tests := []struct {
		name     string
		category string
		episodes []data.ShowSpecialEpisode
		want     string
	}{
		{
			name:     "normal case",
			category: "Halloween",
			episodes: []data.ShowSpecialEpisode{
				{ShowName: "A", Provider: "BBC", Series: 1, Episode: 2},
				{ShowName: "Long Show", Provider: "Netflix", Series: 12, Episode: 10},
			},
			want: "" +
				"Index  Show       Provider  Series  Episode\n" +
				"-----  ---------  --------  ------  -------\n" +
				"1      A          BBC       1       2      \n" +
				"2      Long Show  Netflix   12      10     \n",
		},
		{
			name:     "empty episodes",
			category: "Halloween",
			episodes: []data.ShowSpecialEpisode{},
			want:     "No halloween episodes found.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatSpecialEpisodesTable(tt.category, tt.episodes); got != tt.want {
				t.Fatalf("formatSpecialEpisodesTable() = %q, want %q", got, tt.want)
			}
		})
	}
}
