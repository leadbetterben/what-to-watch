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

func TestViewShowsCancel(t *testing.T) {
	restore := stubCLIHandlers()
	t.Cleanup(restore)

	getCurrentlyWatchingShows = func() ([]data.Show, error) {
		return []data.Show{{Name: "Slow Horses", Genre: "Drama", Provider: "Apple TV+", Series: "1", Episode: "5"}}, nil
	}
	markCalled := false
	markShowWatched = func(int) (bool, error) {
		markCalled = true
		return false, nil
	}

	output := captureOutput(t, func() {
		viewShows(bufio.NewReader(strings.NewReader("0\n")))
	})

	if markCalled {
		t.Fatalf("expected MarkShowWatched not to be called")
	}
	assertContains(t, output, "Slow Horses")
	assertContains(t, output, "No changes made.")
}

func TestViewShowsRejectsNonNumericInput(t *testing.T) {
	restore := stubCLIHandlers()
	t.Cleanup(restore)

	getCurrentlyWatchingShows = func() ([]data.Show, error) {
		return []data.Show{{Name: "Severance", Genre: "Sci-Fi", Provider: "Apple TV+", Series: "1", Episode: "3"}}, nil
	}

	output := captureOutput(t, func() {
		viewShows(bufio.NewReader(strings.NewReader("abc\n")))
	})

	assertContains(t, output, "Invalid input: abc")
}

func TestViewShowsMarksSelectedShowWatched(t *testing.T) {
	tests := []struct {
		name       string
		completed  bool
		wantOutput string
	}{
		{name: "not completed", completed: false, wantOutput: "Show 2 marked as watched."},
		{name: "completed", completed: true, wantOutput: "Show 2 marked as watched and completed!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getCurrentlyWatchingShows = func() ([]data.Show, error) {
				return []data.Show{{Name: "Andor", Genre: "Sci-Fi", Provider: "Disney+", Series: "2", Episode: "4"}}, nil
			}
			var gotIndex int
			markShowWatched = func(idx int) (bool, error) {
				gotIndex = idx
				return tt.completed, nil
			}

			output := captureOutput(t, func() {
				viewShows(bufio.NewReader(strings.NewReader("2\n")))
			})

			if gotIndex != 2 {
				t.Fatalf("expected MarkShowWatched index 2, got %d", gotIndex)
			}
			assertContains(t, output, tt.wantOutput)
		})
	}
}

func TestViewShowsPrintsHandlerErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		input string
	}{
		{
			name: "get shows error",
			setup: func() {
				getCurrentlyWatchingShows = func() ([]data.Show, error) {
					return nil, errors.New("cannot read shows")
				}
			},
		},
		{
			name: "mark watched error",
			setup: func() {
				getCurrentlyWatchingShows = func() ([]data.Show, error) {
					return []data.Show{{Name: "Andor", Genre: "Sci-Fi", Provider: "Disney+", Series: "2", Episode: "4"}}, nil
				}
				markShowWatched = func(int) (bool, error) {
					return false, errors.New("cannot save shows")
				}
			},
			input: "1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)
			tt.setup()

			output := captureOutput(t, func() {
				viewShows(bufio.NewReader(strings.NewReader(tt.input)))
			})

			assertContains(t, output, "Error:")
		})
	}
}

func TestViewFilmsDisplaysFilms(t *testing.T) {
	restore := stubCLIHandlers()
	t.Cleanup(restore)

	getAllFilms = func() ([]data.Film, error) {
		return []data.Film{{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"}}, nil
	}

	output := captureOutput(t, viewFilms)

	assertContains(t, output, "Arrival")
	assertContains(t, output, "Sci-Fi")
}

func TestViewFilmsPrintsHandlerError(t *testing.T) {
	restore := stubCLIHandlers()
	t.Cleanup(restore)

	getAllFilms = func() ([]data.Film, error) {
		return nil, errors.New("cannot read films")
	}

	output := captureOutput(t, viewFilms)

	assertContains(t, output, "Error: cannot read films")
}

func TestViewShowsByGenreDisplaysSelectedGenre(t *testing.T) {
	restore := stubCLIHandlers()
	t.Cleanup(restore)

	getAvailableGenres = func() ([]string, error) {
		return []string{"Drama", "Comedy"}, nil
	}
	var gotGenre string
	getUnwatchedShowsByGenre = func(genre string) ([]data.Show, error) {
		gotGenre = genre
		return []data.Show{{Name: "Unstarted Comedy", Provider: "Netflix"}}, nil
	}

	output := captureOutput(t, func() {
		viewShowsByGenre(bufio.NewReader(strings.NewReader("2\n")))
	})

	if gotGenre != "Comedy" {
		t.Fatalf("expected selected genre Comedy, got %q", gotGenre)
	}
	assertContains(t, output, "Available genres:")
	assertContains(t, output, "Unwatched shows in genre 'Comedy':")
	assertContains(t, output, "Unstarted Comedy")
}

func TestViewShowsByGenreHandlesCancelInvalidAndEmptyGenres(t *testing.T) {
	tests := []struct {
		name       string
		genres     []string
		input      string
		wantOutput string
	}{
		{name: "empty genres", genres: nil, wantOutput: "No genres available."},
		{name: "cancel", genres: []string{"Drama"}, input: "0\n", wantOutput: "No selection made."},
		{name: "invalid text", genres: []string{"Drama"}, input: "abc\n", wantOutput: "Invalid input: abc"},
		{name: "out of range", genres: []string{"Drama"}, input: "2\n", wantOutput: "Invalid input: 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)

			getAvailableGenres = func() ([]string, error) {
				return tt.genres, nil
			}

			output := captureOutput(t, func() {
				viewShowsByGenre(bufio.NewReader(strings.NewReader(tt.input)))
			})

			assertContains(t, output, tt.wantOutput)
		})
	}
}

func TestViewShowsByGenrePrintsHandlerErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		input string
		want  string
	}{
		{
			name: "get genres error",
			setup: func() {
				getAvailableGenres = func() ([]string, error) {
					return nil, errors.New("cannot read genres")
				}
			},
			want: "Error: cannot read genres",
		},
		{
			name: "get shows by genre error",
			setup: func() {
				getAvailableGenres = func() ([]string, error) {
					return []string{"Drama"}, nil
				}
				getUnwatchedShowsByGenre = func(string) ([]data.Show, error) {
					return nil, errors.New("cannot read shows")
				}
			},
			input: "1\n",
			want:  "Error: cannot read shows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := stubCLIHandlers()
			t.Cleanup(restore)
			tt.setup()

			output := captureOutput(t, func() {
				viewShowsByGenre(bufio.NewReader(strings.NewReader(tt.input)))
			})

			assertContains(t, output, tt.want)
		})
	}
}

func stubCLIHandlers() func() {
	originalGetCurrentlyWatchingShows := getCurrentlyWatchingShows
	originalMarkShowWatched := markShowWatched
	originalGetAllFilms := getAllFilms
	originalGetAvailableGenres := getAvailableGenres
	originalGetUnwatchedShowsByGenre := getUnwatchedShowsByGenre

	getCurrentlyWatchingShows = func() ([]data.Show, error) { return nil, nil }
	markShowWatched = func(int) (bool, error) { return false, nil }
	getAllFilms = func() ([]data.Film, error) { return nil, nil }
	getAvailableGenres = func() ([]string, error) { return nil, nil }
	getUnwatchedShowsByGenre = func(string) ([]data.Show, error) { return nil, nil }

	return func() {
		getCurrentlyWatchingShows = originalGetCurrentlyWatchingShows
		markShowWatched = originalMarkShowWatched
		getAllFilms = originalGetAllFilms
		getAvailableGenres = originalGetAvailableGenres
		getUnwatchedShowsByGenre = originalGetUnwatchedShowsByGenre
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
