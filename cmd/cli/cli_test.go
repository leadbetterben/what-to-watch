package cli

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"what-to-watch/data"
)

func TestViewShows(t *testing.T) {
	tests := []struct {
		name         string
		shows        []data.Show
		showErr      error
		completed    bool
		completedErr error
		input        string
		wantIndex    int
		wantOutput   []string
	}{
		{
			name:       "cancel",
			shows:      []data.Show{{Name: "Slow Horses", Genre: "Drama", Provider: "Apple TV+", Series: "1", Episode: "5"}},
			input:      "0\n",
			wantIndex:  -1,
			wantOutput: []string{"Slow Horses", "No changes made."},
		},
		{
			name:       "rejects non-numeric input",
			shows:      []data.Show{{Name: "Severance", Genre: "Sci-Fi", Provider: "Apple TV+", Series: "1", Episode: "3"}},
			input:      "abc\n",
			wantIndex:  -1,
			wantOutput: []string{"Invalid input: abc"},
		},
		{
			name:       "not completed",
			shows:      []data.Show{{Name: "Andor", Genre: "Sci-Fi", Provider: "Disney+", Series: "2", Episode: "4"}},
			input:      "1\n",
			wantIndex:  1,
			wantOutput: []string{"Show 1 marked as watched."},
		},
		{
			name:       "completed",
			shows:      []data.Show{{Name: "Andor", Genre: "Sci-Fi", Provider: "Disney+", Series: "2", Episode: "4"}},
			completed:  true,
			input:      "1\n",
			wantIndex:  1,
			wantOutput: []string{"Show 1 marked as watched and completed!"},
		},
		{
			name:       "get shows error",
			shows:      nil,
			showErr:    errors.New("cannot read shows"),
			wantIndex:  -1,
			wantOutput: []string{"Error: cannot read shows"},
		},
		{
			name:         "mark watched error",
			shows:        []data.Show{{Name: "Andor", Genre: "Sci-Fi", Provider: "Disney+", Series: "2", Episode: "4"}},
			completedErr: errors.New("cannot save shows"),
			input:        "1\n",
			wantIndex:    1,
			wantOutput:   []string{"Error: cannot save shows"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getCurrentlyWatchingShows = func() ([]data.Show, error) {
				return tt.shows, tt.showErr
			}

			gotIndex := -1
			markShowWatched = func(idx int) (bool, error) {
				gotIndex = idx
				return tt.completed, tt.completedErr
			}

			output := captureOutput(t, func() {
				viewShows(bufio.NewReader(strings.NewReader(tt.input)))
			})

			if gotIndex != tt.wantIndex {
				t.Fatalf("expected MarkShowWatched index %d, got %d", tt.wantIndex, gotIndex)
			}
			for _, want := range tt.wantOutput {
				assertContains(t, output, want)
			}
		})
	}
}

func TestViewFilms(t *testing.T) {
	tests := []struct {
		name       string
		films      []data.Film
		filmsErr   error
		wantOutput []string
	}{
		{
			name:       "normal case",
			films:      []data.Film{{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"}},
			wantOutput: []string{"Arrival", "Sci-Fi", "Netflix"},
		},
		{
			name:       "error",
			filmsErr:   errors.New("cannot read films"),
			wantOutput: []string{"Error: cannot read films"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getAllFilms = func() ([]data.Film, error) {
				return tt.films, tt.filmsErr
			}

			output := captureOutput(t, viewFilms)

			for _, want := range tt.wantOutput {
				assertContains(t, output, want)
			}
		})
	}
}

func TestViewShowsByGenre(t *testing.T) {
	tests := []struct {
		name       string
		genres     []string
		genresErr  error
		shows      []data.Show
		showsErr   error
		input      string
		wantOutput []string
	}{
		{
			name:       "normal case",
			genres:     []string{"Drama", "Comedy"},
			shows:      []data.Show{{Name: "Unstarted Comedy", Provider: "Netflix"}},
			input:      "2\n",
			wantOutput: []string{"Available genres:", "1. Drama", "2. Comedy", "Unwatched shows in genre 'Comedy':", "Unstarted Comedy"},
		},
		{
			name:       "get genres error",
			genresErr:  errors.New("cannot read genres"),
			input:      "2\n",
			wantOutput: []string{"Error: cannot read genres"},
		},
		{
			name:       "get shows by genre error",
			genres:     []string{"Drama", "Comedy"},
			showsErr:   errors.New("cannot read shows by genre"),
			input:      "2\n",
			wantOutput: []string{"Error: cannot read shows by genre"},
		},
		{
			name:       "empty genres",
			genres:     []string{},
			wantOutput: []string{"No genres available."},
		},
		{
			name:       "cancel",
			genres:     []string{"Drama"},
			input:      "0\n",
			wantOutput: []string{"Available genres:", "1. Drama", "No selection made."},
		},
		{
			name:       "invalid text",
			genres:     []string{"Drama"},
			input:      "abc\n",
			wantOutput: []string{"Available genres:", "1. Drama", "Invalid input: abc"},
		},
		{
			name:       "out of range",
			genres:     []string{"Drama"},
			input:      "2\n",
			wantOutput: []string{"Available genres:", "1. Drama", "Invalid input: 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getAvailableGenres = func() ([]string, error) {
				return tt.genres, tt.genresErr
			}
			getUnwatchedShowsByGenre = func(string) ([]data.Show, error) {
				return tt.shows, tt.showsErr
			}

			output := captureOutput(t, func() {
				viewShowsByGenre(bufio.NewReader(strings.NewReader(tt.input)))
			})

			for _, want := range tt.wantOutput {
				assertContains(t, output, want)
			}
		})
	}
}

func TestViewSeasonalSpecials(t *testing.T) {
	tests := []struct {
		name        string
		episodes    []data.ShowSpecialEpisode
		episodesErr error
		wantOutput  []string
	}{
		{
			name:       "success",
			episodes:   []data.ShowSpecialEpisode{{ShowName: "Example Show", Provider: "Netflix", Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5}},
			wantOutput: []string{"Example Show", "Netflix", "TEST episodes:"},
		},
		{
			name:       "empty episodes",
			episodes:   []data.ShowSpecialEpisode{},
			wantOutput: []string{"No test episodes found."},
		},
		{
			name:        "error",
			episodesErr: errors.New("cannot read episodes"),
			wantOutput:  []string{"Error: cannot read episodes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getEpisodes := func() ([]data.ShowSpecialEpisode, error) {
				return tt.episodes, tt.episodesErr
			}

			output := captureOutput(t, func() {
				viewSeasonalSpecials("TEST", getEpisodes)
			})

			for _, want := range tt.wantOutput {
				assertContains(t, output, want)
			}
		})
	}
}

func TestViewRandomSeasonalSpecial(t *testing.T) {
	tests := []struct {
		name       string
		episode    data.ShowSpecialEpisode
		episodeErr error
		wantOutput []string
	}{
		{
			name:       "success",
			episode:    data.ShowSpecialEpisode{ShowName: "Example Show", Provider: "Netflix", Category: data.SpecialCategoryHalloween, Series: 2, Episode: 5},
			wantOutput: []string{"Random test episode: Example Show (Netflix), Series 2 Episode 5"},
		},
		{
			name:       "error",
			episodeErr: errors.New("cannot read episode"),
			wantOutput: []string{"Error: cannot read episode"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getEpisode := func() (data.ShowSpecialEpisode, error) {
				return tt.episode, tt.episodeErr
			}

			output := captureOutput(t, func() {
				viewRandomSeasonalSpecial("TEST", getEpisode)
			})

			for _, want := range tt.wantOutput {
				assertContains(t, output, want)
			}
		})
	}
}

func stubCLIHandlers() func() {
	originalGetCurrentlyWatchingShows := getCurrentlyWatchingShows
	originalMarkShowWatched := markShowWatched
	originalGetAllFilms := getAllFilms
	originalGetAvailableGenres := getAvailableGenres
	originalGetUnwatchedShowsByGenre := getUnwatchedShowsByGenre
	originalGetHalloweenEpisodes := getHalloweenEpisodes
	originalGetChristmasEpisodes := getChristmasEpisodes
	originalGetRandomHalloweenEpisode := getRandomHalloweenEpisode
	originalGetRandomChristmasEpisode := getRandomChristmasEpisode

	getCurrentlyWatchingShows = func() ([]data.Show, error) { return nil, nil }
	markShowWatched = func(int) (bool, error) { return false, nil }
	getAllFilms = func() ([]data.Film, error) { return nil, nil }
	getAvailableGenres = func() ([]string, error) { return nil, nil }
	getUnwatchedShowsByGenre = func(string) ([]data.Show, error) { return nil, nil }
	getHalloweenEpisodes = func() ([]data.ShowSpecialEpisode, error) { return nil, nil }
	getChristmasEpisodes = func() ([]data.ShowSpecialEpisode, error) { return nil, nil }
	getRandomHalloweenEpisode = func() (data.ShowSpecialEpisode, error) { return data.ShowSpecialEpisode{}, nil }
	getRandomChristmasEpisode = func() (data.ShowSpecialEpisode, error) { return data.ShowSpecialEpisode{}, nil }

	return func() {
		getCurrentlyWatchingShows = originalGetCurrentlyWatchingShows
		markShowWatched = originalMarkShowWatched
		getAllFilms = originalGetAllFilms
		getAvailableGenres = originalGetAvailableGenres
		getUnwatchedShowsByGenre = originalGetUnwatchedShowsByGenre
		getHalloweenEpisodes = originalGetHalloweenEpisodes
		getChristmasEpisodes = originalGetChristmasEpisodes
		getRandomHalloweenEpisode = originalGetRandomHalloweenEpisode
		getRandomChristmasEpisode = originalGetRandomChristmasEpisode
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}

	os.Stdout = writer
	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("closing writer: %v", err)
	}
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("closing reader: %v", err)
	}

	return string(output)
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("expected output to contain %q, got %q", want, got)
	}
}
