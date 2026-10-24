package db

import (
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
    "specialEpisodes": [
      {"category": "halloween", "series": 1, "episode": 1},
      {"category": "christmas", "series": 2, "episode": 13}
    ],
    "currentSeries": 1,
    "currentEpisode": 2
  }
]`)

	shows, err := ReadShows()
	if err != nil {
		t.Fatalf("ReadAllShows returned error: %v", err)
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
	if len(shows[0].SpecialEpisodes) != 2 || shows[0].SpecialEpisodes[0].Episode != 1 || shows[0].SpecialEpisodes[1].Episode != 13 {
		t.Fatalf("expected seasonal metadata at valid episode boundaries, got %+v", shows[0].SpecialEpisodes)
	}
}

func TestReadShowsRejectsInvalidSeasonalMetadata(t *testing.T) {
	tests := []struct {
		name  string
		entry string
		field string
	}{
		{name: "unsupported category", entry: `{"category":"thanksgiving","series":1,"episode":1}`, field: ".category"},
		{name: "zero series", entry: `{"category":"halloween","series":0,"episode":1}`, field: ".series"},
		{name: "negative series", entry: `{"category":"halloween","series":-1,"episode":1}`, field: ".series"},
		{name: "zero episode", entry: `{"category":"christmas","series":1,"episode":0}`, field: ".episode"},
		{name: "negative episode", entry: `{"category":"christmas","series":1,"episode":-1}`, field: ".episode"},
		{name: "series above bounds", entry: `{"category":"halloween","series":3,"episode":1}`, field: ".series"},
		{name: "episode above bounds", entry: `{"category":"christmas","series":2,"episode":14}`, field: ".episode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			fixture := `[{"name":"Example","episodes":[10,13],"specialEpisodes":[` + tt.entry + `]}]`
			writeFixture(t, dir, "shows.json", fixture)
			if _, err := ReadShows(); err == nil || !strings.Contains(err.Error(), "Example") || !strings.Contains(err.Error(), tt.field) {
				t.Fatalf("expected contextual validation error for %s, got %v", tt.name, err)
			}
			raw, err := os.ReadFile(filepath.Join(dir, "shows.json"))
			if err != nil {
				t.Fatalf("reading fixture after validation: %v", err)
			}
			if string(raw) != fixture {
				t.Fatalf("validation changed source data: got %s", raw)
			}
		})
	}
}

func TestReadShowsRejectsInvalidCurrentPositions(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{name: "partial position", json: `[{"name":"Partial","episodes":[2],"currentSeries":1}]`},
		{name: "series outside bounds", json: `[{"name":"Bad Series","episodes":[2],"currentSeries":2,"currentEpisode":1}]`},
		{name: "episode outside bounds", json: `[{"name":"Bad Episode","episodes":[2],"currentSeries":1,"currentEpisode":3}]`},
		{name: "non-positive episode count", json: `[{"name":"Bad Count","episodes":[0],"currentSeries":1,"currentEpisode":1}]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			writeFixture(t, dir, "shows.json", tt.json)
			if _, err := ReadShows(); err == nil {
				t.Fatalf("expected validation error for %s", tt.name)
			}
		})
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

func TestReadFileReturnsErrorWhenPathCannotBeResolved(t *testing.T) {
	t.Cleanup(SetFullPathResolverForTest(func(string) string { return "" }))

	if _, err := readFile("shows.json"); err == nil {
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
