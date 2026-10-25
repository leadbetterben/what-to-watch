package db

import (
	"errors"
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

func TestUpdateShows(t *testing.T) {
	originalShows := []data.Show{
		{Name: "Current", Episodes: []int{2}, CurrentSeries: intPointer(1), CurrentEpisode: intPointer(1)},
		{Name: "Queued", Episodes: []int{3}},
	}
	updatedShows := []data.Show{
		{Name: "Current", Episodes: []int{2}, CurrentSeries: intPointer(1), CurrentEpisode: intPointer(2)},
		{Name: "Queued", Episodes: []int{3}},
	}
	fixture := `[{"name":"Current","episodes":[2],"currentSeries":1,"currentEpisode":1},{"name":"Queued","episodes":[3]}]`

	tests := []struct {
		name         string
		fixture      string
		callbackErr  error
		writeErr     error
		wantErr      string
		wantCallback bool
		wantWrite    bool
	}{
		{
			name:    "read error",
			wantErr: "UpdateShows: error reading shows.json",
		},
		{
			name:         "callback error",
			fixture:      fixture,
			callbackErr:  errors.New("callback failed"),
			wantErr:      "callback failed",
			wantCallback: true,
		},
		{
			name:         "write error",
			fixture:      fixture,
			writeErr:     errors.New("write failed"),
			wantErr:      "UpdateShows: error writing shows.json: write failed",
			wantCallback: true,
			wantWrite:    true,
		},
		{
			name:         "success",
			fixture:      fixture,
			wantCallback: true,
			wantWrite:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := useTempDataDir(t)
			if tt.fixture != "" {
				writeFixture(t, dir, "shows.json", tt.fixture)
			}

			originalWriteShowsFn := writeShowsFn
			writeCalled := false
			writeShowsFn = func(shows []data.Show) error {
				writeCalled = true
				if !reflect.DeepEqual(shows, updatedShows) {
					t.Errorf("write callback received unexpected shows: %+v", shows)
				}
				if tt.writeErr != nil {
					return tt.writeErr
				}
				return writeShowsFile(shows)
			}
			t.Cleanup(func() { writeShowsFn = originalWriteShowsFn })

			callbackCalled := false
			err := UpdateShows(func(shows []data.Show) ([]data.Show, error) {
				callbackCalled = true
				if !reflect.DeepEqual(shows, originalShows) {
					t.Errorf("callback received unexpected shows: %+v", shows)
				}
				if tt.callbackErr != nil {
					return nil, tt.callbackErr
				}
				return updatedShows, nil
			})

			if callbackCalled != tt.wantCallback {
				t.Fatalf("callback called = %v, want %v", callbackCalled, tt.wantCallback)
			}
			if writeCalled != tt.wantWrite {
				t.Fatalf("write called = %v, want %v", writeCalled, tt.wantWrite)
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateShows returned error: %v", err)
			}
			got, err := ReadShows()
			if err != nil {
				t.Fatalf("ReadShows returned error after update: %v", err)
			}
			if !reflect.DeepEqual(got, updatedShows) {
				t.Fatalf("expected persisted shows %+v, got %+v", updatedShows, got)
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
