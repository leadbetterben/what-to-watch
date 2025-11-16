package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"what-to-watch/data"
)

// ReadShows reads the shows from the shows.json file and returns a slice of Show structs.
func ReadShows() ([]data.Show, error) {
	fullPath := getFullPath("shows.json")
	if fullPath == "" {
		return nil, fmt.Errorf("ReadShows: could not determine full path to shows.json")
	}

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("readFile: error reading file \n err=%w fullPath=%s", err, fullPath)
	}

	var shows []data.Show
	if err := json.Unmarshal(raw, &shows); err != nil {
		return nil, err
	}

	return shows, nil
}

// ReadFilms reads the films from the films.json file and returns a slice of Film structs.
func ReadFilms() ([]data.Film, error) {
	fullPath := getFullPath("films.json")
	if fullPath == "" {
		return nil, fmt.Errorf("ReadFilms: could not determine full path to films.json")
	}

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("ReadFilms: error reading file \n err=%w fullPath=%s", err, fullPath)
	}

	var films []data.Film
	if err := json.Unmarshal(raw, &films); err != nil {
		return nil, err
	}

	return films, nil
}

// WriteShows writes the provided shows slice to the shows.json file.
// It writes to a temporary file in the same directory and renames it
// to avoid corrupting the file on failure.
func WriteShows(shows []data.Show) error {
	raw, err := json.MarshalIndent(shows, "", "  ")
	if err != nil {
		return err
	}

	fullPath := getFullPath("shows.json")
	if fullPath == "" {
		return fmt.Errorf("WriteShows: could not determine full path to shows.json")
	}

	// create temp file in same directory to ensure atomic rename
	dir := filepath.Dir(fullPath)
	tmpFile, err := os.CreateTemp(dir, "shows-*.json.tmp")
	if err != nil {
		return fmt.Errorf("WriteShows: error creating temp file \n err=%w fullPath=%s", err, fullPath)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// write to temp file
	if _, err := tmpFile.Write(raw); err != nil {
		tmpFile.Close()
		return fmt.Errorf("WriteShows: error writing temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("WriteShows: error closing temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
	}

	// rename temp file to final file
	if err := os.Rename(tmpPath, fullPath); err != nil {
		return fmt.Errorf("WriteShows: error renaming temp file \n err=%w fullPath=%s tmpPath=%s", err, fullPath, tmpPath)
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
