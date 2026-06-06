package db

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"what-to-watch/data"
)

func TestReadShows(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "shows.json", `[
  {
    "name": "The Expanse",
    "genre": "Sci-Fi",
    "episodes": [10, 13],
    "provider": "Prime",
    "currentSeries": 1,
    "currentEpisode": 2
  }
]`)

	shows, err := ReadShows()
	if err != nil {
		t.Fatalf("ReadShows returned error: %v", err)
	}

	if len(shows) != 1 {
		t.Fatalf("expected 1 show, got %d", len(shows))
	}
	if shows[0].Name != "The Expanse" || shows[0].Genre != "Sci-Fi" || shows[0].Provider != "Prime" {
		t.Fatalf("unexpected show: %+v", shows[0])
	}
	if shows[0].CurrentSeries == nil || *shows[0].CurrentSeries != 1 {
		t.Fatalf("expected current series 1, got %+v", shows[0].CurrentSeries)
	}
	if shows[0].CurrentEpisode == nil || *shows[0].CurrentEpisode != 2 {
		t.Fatalf("expected current episode 2, got %+v", shows[0].CurrentEpisode)
	}
}

func TestReadCurrentShows(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "currentShows.json", `[
  {
    "name": "Slow Horses",
    "genre": "Drama",
    "episodes": [6],
    "provider": "Apple TV+",
    "currentSeries": 1,
    "currentEpisode": 5
  }
]`)

	shows, err := ReadCurrentShows()
	if err != nil {
		t.Fatalf("ReadCurrentShows returned error: %v", err)
	}

	if len(shows) != 1 {
		t.Fatalf("expected 1 current show, got %d", len(shows))
	}
	if shows[0].Name != "Slow Horses" || shows[0].CurrentEpisode == nil || *shows[0].CurrentEpisode != 5 {
		t.Fatalf("unexpected current show: %+v", shows[0])
	}
}

func TestReadFilms(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "films.json", `[
  {
    "name": "Arrival",
    "genre": "Sci-Fi",
    "provider": "Netflix"
  }
]`)

	films, err := ReadFilms()
	if err != nil {
		t.Fatalf("ReadFilms returned error: %v", err)
	}

	if len(films) != 1 {
		t.Fatalf("expected 1 film, got %d", len(films))
	}
	if films[0] != (data.Film{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"}) {
		t.Fatalf("unexpected film: %+v", films[0])
	}
}

func TestReadersReturnErrorForMissingFiles(t *testing.T) {
	useTempDataDir(t)

	tests := []struct {
		name string
		read func() error
	}{
		{name: "ReadShows", read: func() error { _, err := ReadShows(); return err }},
		{name: "ReadCurrentShows", read: func() error { _, err := ReadCurrentShows(); return err }},
		{name: "ReadFilms", read: func() error { _, err := ReadFilms(); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.read(); err == nil {
				t.Fatalf("expected error for missing file")
			}
		})
	}
}

func TestReadersReturnErrorForInvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		read     func() error
	}{
		{name: "ReadShows", fileName: "shows.json", read: func() error { _, err := ReadShows(); return err }},
		{name: "ReadCurrentShows", fileName: "currentShows.json", read: func() error { _, err := ReadCurrentShows(); return err }},
		{name: "ReadFilms", fileName: "films.json", read: func() error { _, err := ReadFilms(); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			writeFixture(t, dir, tt.fileName, `{bad json`)

			if err := tt.read(); err == nil {
				t.Fatalf("expected error for invalid JSON")
			}
		})
	}
}

func TestWriteCurrentShowsWritesIndentedJSON(t *testing.T) {
	dir := useTempDataDir(t)

	shows := []data.Show{
		{
			Name:           "Severance",
			Genre:          "Sci-Fi",
			Episodes:       []int{9},
			Provider:       "Apple TV+",
			CurrentSeries:  intPtr(1),
			CurrentEpisode: intPtr(3),
		},
	}

	if err := WriteCurrentShows(shows); err != nil {
		t.Fatalf("WriteCurrentShows returned error: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "currentShows.json"))
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}

	if !strings.Contains(string(raw), "\n  {") {
		t.Fatalf("expected indented JSON, got %q", string(raw))
	}

	var written []data.Show
	if err := json.Unmarshal(raw, &written); err != nil {
		t.Fatalf("written file is not valid JSON: %v", err)
	}
	if len(written) != 1 || written[0].Name != "Severance" {
		t.Fatalf("unexpected written shows: %+v", written)
	}
}

func TestWriteCurrentShowsReplacesExistingFileAndRemovesTempFile(t *testing.T) {
	dir := useTempDataDir(t)
	writeFixture(t, dir, "currentShows.json", `[]`)

	if err := WriteCurrentShows([]data.Show{{Name: "Andor", Genre: "Sci-Fi"}}); err != nil {
		t.Fatalf("WriteCurrentShows returned error: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "currentShows.json"))
	if err != nil {
		t.Fatalf("reading replaced file: %v", err)
	}
	if !strings.Contains(string(raw), `"name": "Andor"`) {
		t.Fatalf("expected file to be replaced with new show, got %s", string(raw))
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*.tmp"))
	if err != nil {
		t.Fatalf("checking temp files: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no temp files, found %v", matches)
	}
}

func TestReadFileReturnsErrorWhenPathCannotBeResolved(t *testing.T) {
	t.Cleanup(SetFullPathResolverForTest(func(string) string { return "" }))

	if _, err := readFile("shows.json"); err == nil {
		t.Fatalf("expected error when path cannot be resolved")
	}
}

func TestWriteCurrentShowsReturnsErrorWhenPathCannotBeResolved(t *testing.T) {
	t.Cleanup(SetFullPathResolverForTest(func(string) string { return "" }))

	if err := WriteCurrentShows([]data.Show{}); err == nil {
		t.Fatalf("expected error when path cannot be resolved")
	}
}

func useTempDataDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Cleanup(SetFullPathResolverForTest(func(path string) string {
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

func intPtr(i int) *int {
	return &i
}
