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
	raw, err := readFile("shows.json")
	if err != nil {
		return nil, err
	}

	var shows []data.Show
	if err := json.Unmarshal(raw, &shows); err != nil {
		return nil, err
	}

	return shows, nil
}

// readFile tries to read the file at the given path.
// It first attempts to locate the file relative to the executable's directory.
// If that fails, it falls back to locating the file relative to the source code directory.
func readFile(path string) ([]byte, error) {
	// Try to get path relative to executable first (for built binaries)
	var fullPath string
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

	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("readFile: error reading file \n err=%w path=%s fullPath=%s", err, path, fullPath)
	}

	return raw, nil
}
