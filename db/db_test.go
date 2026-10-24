package db

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"what-to-watch/data"
)

func TestReadShows(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    []data.Show
		wantErr []string
	}{
		{
			name:    "success",
			fixture: `[{"name":"The Expanse","genre":"Sci-Fi","episodes":[10,13],"provider":"Prime","specialEpisodes":[{"category":"halloween","series":1,"episode":1},{"category":"christmas","series":2,"episode":13}],"currentSeries":1,"currentEpisode":2}]`,
			want: []data.Show{{
				Name: "The Expanse", Genre: "Sci-Fi", Episodes: []int{10, 13}, Provider: "Prime",
				SpecialEpisodes: []data.SpecialEpisode{{Category: data.SpecialCategoryHalloween, Series: 1, Episode: 1}, {Category: data.SpecialCategoryChristmas, Series: 2, Episode: 13}},
				CurrentSeries:   intPointer(1), CurrentEpisode: intPointer(2),
			}},
		},
		{name: "missing file", wantErr: []string{"shows.json"}},
		{name: "invalid JSON", fixture: `{bad json`, wantErr: []string{"invalid character"}},
		{name: "unsupported seasonal category", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"thanksgiving","series":1,"episode":1}]}]`, wantErr: []string{"Example", ".category"}},
		{name: "zero seasonal series", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"halloween","series":0,"episode":1}]}]`, wantErr: []string{"Example", ".series"}},
		{name: "negative seasonal series", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"halloween","series":-1,"episode":1}]}]`, wantErr: []string{"Example", ".series"}},
		{name: "zero seasonal episode", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"christmas","series":1,"episode":0}]}]`, wantErr: []string{"Example", ".episode"}},
		{name: "negative seasonal episode", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"christmas","series":1,"episode":-1}]}]`, wantErr: []string{"Example", ".episode"}},
		{name: "seasonal series above bounds", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"halloween","series":3,"episode":1}]}]`, wantErr: []string{"Example", ".series"}},
		{name: "seasonal episode above bounds", fixture: `[{"name":"Example","episodes":[10,13],"specialEpisodes":[{"category":"christmas","series":2,"episode":14}]}]`, wantErr: []string{"Example", ".episode"}},
		{name: "partial current position", fixture: `[{"name":"Partial","episodes":[2],"currentSeries":1}]`, wantErr: []string{"Partial", "both currentSeries and currentEpisode"}},
		{name: "current series outside bounds", fixture: `[{"name":"Bad Series","episodes":[2],"currentSeries":2,"currentEpisode":1}]`, wantErr: []string{"Bad Series", "currentSeries", "outside episodes bounds"}},
		{name: "current episode outside bounds", fixture: `[{"name":"Bad Episode","episodes":[2],"currentSeries":1,"currentEpisode":3}]`, wantErr: []string{"Bad Episode", "currentEpisode", "outside series"}},
		{name: "non-positive episode count", fixture: `[{"name":"Bad Count","episodes":[0],"currentSeries":1,"currentEpisode":1}]`, wantErr: []string{"Bad Count", "invalid episodes[1] count 0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			if tt.fixture != "" {
				writeFixture(t, dir, "shows.json", tt.fixture)
			}
			shows, err := ReadShows()
			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				for _, fragment := range tt.wantErr {
					if !strings.Contains(err.Error(), fragment) {
						t.Fatalf("expected error containing %q, got %q", fragment, err.Error())
					}
				}
				raw, readErr := os.ReadFile(filepath.Join(dir, "shows.json"))
				if readErr == nil && string(raw) != tt.fixture {
					t.Fatalf("validation changed source data: got %s", raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadShows returned error: %v", err)
			}
			if len(shows) != len(tt.want) {
				t.Fatalf("expected %d shows, got %d", len(tt.want), len(shows))
			}
			for i, want := range tt.want {
				if !reflect.DeepEqual(shows[i], want) {
					t.Fatalf("unexpected show at index %d: %+v", i, shows[i])
				}
			}
		})
	}
}

func TestReadFilms(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    []data.Film
		wantErr string
	}{
		{
			name:    "success",
			fixture: `[{"name":"Arrival","genre":"Sci-Fi","provider":"Netflix"}]`,
			want:    []data.Film{{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"}},
		},
		{
			name:    "missing file",
			fixture: "",
			want:    nil,
			wantErr: "ReadFilms: error reading file: readFile: error reading file films.json: open",
		},
		{
			name:    "invalid JSON",
			fixture: `{"bad json"`,
			want:    nil,
			wantErr: "ReadFilms: error unmarshalling JSON: unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			if tt.fixture != "" {
				writeFixture(t, dir, "films.json", tt.fixture)
			}
			films, err := ReadFilms()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadFilms returned error: %v", err)
			}
			if len(films) != len(tt.want) {
				t.Fatalf("expected %d films, got %d", len(tt.want), len(films))
			}
			for i, want := range tt.want {
				if films[i] != want {
					t.Fatalf("unexpected film at index %d: %+v", i, films[i])
				}
			}
		})
	}
}

func useTempDataDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	originalResolver := fullPathResolver
	fullPathResolver = func(path string) string {
		return filepath.Join(dir, path)
	}
	t.Cleanup(func() {
		fullPathResolver = originalResolver
	})

	return dir
}

func writeFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("writing fixture %s: %v", name, err)
	}
}

func intPointer(value int) *int {
	return &value
}
