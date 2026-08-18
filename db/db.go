package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"what-to-watch/data"
)

var fullPathResolver = getFullPath

// SetFullPathResolverForTest replaces the data file path resolver and returns
// a restore function. It is intended for tests that need isolated fixture files.
func SetFullPathResolverForTest(resolver func(string) string) func() {
	originalResolver := fullPathResolver
	fullPathResolver = resolver
	return func() {
		fullPathResolver = originalResolver
	}
}

func readShows() ([]data.Show, error) {
	raw, err := readFile("shows.json")
	if err != nil {
		return nil, fmt.Errorf("readShows: error reading file \n err=%w", err)
	}

	var shows []data.Show
	if err := json.Unmarshal(raw, &shows); err != nil {
		return nil, err
	}

	return shows, nil
}

// ReadAllShows reads every show from shows.json.
func ReadAllShows() ([]data.Show, error) {
	return readShows()
}

// ReadUnwatchedShows returns shows that have no current position and are not marked for rewatch.
func ReadUnwatchedShows() ([]data.Show, error) {
	shows, err := readShows()
	if err != nil {
		return nil, fmt.Errorf("ReadUnwatchedShows: %w", err)
	}
	filtered := make([]data.Show, 0)
	for _, show := range shows {
		if show.CurrentSeries == nil && show.CurrentEpisode == nil && !show.Rewatch {
			filtered = append(filtered, show)
		}
	}
	return filtered, nil
}

// ReadCurrentShows returns shows with a current series or episode position.
func ReadCurrentShows() ([]data.Show, error) {
	shows, err := readShows()
	if err != nil {
		return nil, fmt.Errorf("ReadCurrentShows: %w", err)
	}
	filtered := make([]data.Show, 0)
	for _, show := range shows {
		if show.CurrentSeries != nil || show.CurrentEpisode != nil {
			filtered = append(filtered, show)
		}
	}
	return filtered, nil
}

// ReadRewatchShows returns shows marked for rewatch.
func ReadRewatchShows() ([]data.Show, error) {
	shows, err := readShows()
	if err != nil {
		return nil, fmt.Errorf("ReadRewatchShows: %w", err)
	}
	filtered := make([]data.Show, 0)
	for _, show := range shows {
		if show.Rewatch {
			filtered = append(filtered, show)
		}
	}
	return filtered, nil
}

// ReadFilms reads the films from the films.json file and returns a slice of Film structs.
func ReadFilms() ([]data.Film, error) {
	raw, err := readFile("films.json")
	if err != nil {
		return nil, fmt.Errorf("ReadFilms: error reading file \n err=%w", err)
	}

	var films []data.Film
	if err := json.Unmarshal(raw, &films); err != nil {
		return nil, err
	}

	return films, nil
}

// WriteCurrentShows replaces the current-position records in shows.json with the provided slice.
// It writes to a temporary file in the same directory and renames it
// to avoid corrupting the file on failure.
func WriteCurrentShows(shows []data.Show) error {
	existing, err := readShows()
	if err != nil {
		return fmt.Errorf("WriteCurrentShows: error reading shows.json \n err=%w", err)
	}

	merged := make([]data.Show, 0, len(existing)+len(shows))
	for _, show := range existing {
		if show.CurrentSeries == nil && show.CurrentEpisode == nil {
			merged = append(merged, show)
		}
	}
	merged = append(merged, shows...)

	raw, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}

	fullPath := fullPathResolver("shows.json")
	if fullPath == "" {
		return fmt.Errorf("WriteCurrentShows: could not determine full path to shows.json")
	}

	// create temp file in same directory to ensure atomic rename
	dir := filepath.Dir(fullPath)
	tmpFile, err := os.CreateTemp(dir, "shows-*.json.tmp")
	if err != nil {
		return fmt.Errorf("WriteCurrentShows: error creating temp file \n err=%w fullPath=%s", err, fullPath)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// write to temp file
	if _, err := tmpFile.Write(raw); err != nil {
		tmpFile.Close()
		return fmt.Errorf("WriteCurrentShows: error writing temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("WriteCurrentShows: error closing temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
	}

	// rename temp file to final file
	if err := os.Rename(tmpPath, fullPath); err != nil {
		return fmt.Errorf("WriteCurrentShows: error renaming temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
	}

	return nil
}

// getFullPath attempts to determine the full path to the given file.
func getFullPath(path string) (fullPath string) {
	// Try to get path relative to executable first (for built binaries)
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidatePath := filepath.Join(exeDir, path)
		if _, err := os.Stat(candidatePath); err == nil {
			fullPath = candidatePath
		}
	}

	// Fall back to source directory (for go run during development)
	if fullPath == "" {
		_, currentFile, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(currentFile)
		fullPath = filepath.Join(sourceDir, path)
	}

	return
}

// readFile reads the contents of the file at the given path.
func readFile(path string) ([]byte, error) {
	fullPath := fullPathResolver(path)
	if fullPath == "" {
		return nil, fmt.Errorf("readFile: could not determine full path to %s", path)
	}

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("readFile: error reading file \n err=%w path=%s fullPath=%s", err, path, fullPath)
	}

	return raw, nil
}
