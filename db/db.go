package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"what-to-watch/data"
)

var fullPathResolver = getFullPath
var writeShowsFn = writeShowsFile

// showsWriteMu serializes show read-modify-write operations in this process.
// Multiple processes writing the same shows.json file are not supported.
var showsWriteMu sync.Mutex

func ReadShows() ([]data.Show, error) {
	raw, err := readFile("shows.json")
	if err != nil {
		return nil, fmt.Errorf("readShows: error reading file: %w", err)
	}

	var shows []data.Show
	if err := json.Unmarshal(raw, &shows); err != nil {
		return nil, err
	}
	for _, show := range shows {
		if err := data.ValidateEpisodeCounts(show); err != nil {
			return nil, err
		}
		if err := data.ValidateSpecialEpisodes(show); err != nil {
			return nil, err
		}
		if err := data.ValidateCurrentPosition(show); err != nil {
			return nil, err
		}
	}

	return shows, nil
}

// ReadFilms reads the films from the films.json file and returns a slice of Film structs.
func ReadFilms() ([]data.Film, error) {
	raw, err := readFile("films.json")
	if err != nil {
		return nil, fmt.Errorf("ReadFilms: error reading file: %w", err)
	}

	var films []data.Film
	if err := json.Unmarshal(raw, &films); err != nil {
		return nil, fmt.Errorf("ReadFilms: error unmarshalling JSON: %w", err)
	}

	return films, nil
}

// UpdateShows reads the latest complete shows, applies the provided update function and then writes the updated list to the store.
// It holds a lock on the data whilst updating.
func UpdateShows(update func([]data.Show) ([]data.Show, error)) error {
	showsWriteMu.Lock()
	defer showsWriteMu.Unlock()

	shows, err := ReadShows()
	if err != nil {
		return fmt.Errorf("UpdateShows: error reading shows.json: %w", err)
	}

	updated, err := update(shows)
	if err != nil {
		return err
	}

	if err := writeShowsFn(updated); err != nil {
		return fmt.Errorf("UpdateShows: error writing shows.json: %w", err)
	}

	return nil
}

// writeShows atomically replaces shows.json with the supplied complete list.
// Callers that perform read-modify-write operations must hold showsWriteMu.
func writeShowsFile(shows []data.Show) error {
	raw, err := json.MarshalIndent(shows, "", "  ")
	if err != nil {
		return err
	}

	fullPath := fullPathResolver("shows.json")
	if fullPath == "" {
		return fmt.Errorf("could not determine full path to shows.json")
	}

	// create temp file in same directory to ensure atomic rename
	dir := filepath.Dir(fullPath)
	tmpFile, err := os.CreateTemp(dir, "shows-*.json.tmp")
	if err != nil {
		return fmt.Errorf("error creating temp file at %s: %w", fullPath, err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// write to temp file
	if _, err := tmpFile.Write(raw); err != nil {
		tmpFile.Close()
		return fmt.Errorf("error writing to temp file at %s: %w", tmpPath, err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("error closing temp file at %s: %w", tmpPath, err)
	}

	// rename temp file to final file
	if err := os.Rename(tmpPath, fullPath); err != nil {
		return fmt.Errorf("error renaming temp file at %s to %s: %w", tmpPath, fullPath, err)
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
		return nil, fmt.Errorf("readFile: error reading file %s: %w", path, err)
	}

	return raw, nil
}
