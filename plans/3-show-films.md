# Implementation Plan for Issue 3: Show Films

## Overview
Add functionality to display films from `films.json` as an alternative to shows. Users should see a menu on startup to choose between viewing currently watching shows or viewing films.

## Issue Requirements
1. Show users two options when running the application
   - Option 1: View currently watching shows (existing behavior)
   - Option 2: View films
2. Display films in table format to the console
3. Show useful error messages for invalid input (not 1 or 2)
4. Create a new function in the `db` package to read `films.json`

## Implementation Steps

### Step 1: Extend the data model
- **File**: `data/data.go`
- **Change**: Create a new `Film` struct that mirrors the film structure in `films.json`
  - Fields: `Name`, `Genre`, `Provider`
  - Use JSON struct tags matching the JSON file format
- **Rationale**: The `Show` struct has fields specific to show tracking (episodes, current series, etc.) that don't apply to films

### Step 2: Add film reading function to db package
- **File**: `db/db.go`
- **Change**: Implement `ReadFilms()` function following the same pattern as `ReadShows()`
  - Read from `films.json` using `getFullPath("films.json")`
  - Unmarshal JSON into a slice of `Film` structs
  - Return `[]data.Film` and error
- **Rationale**: Reuses existing file path resolution logic; consistent with current patterns

### Step 3: Refactor table rendering logic
- **File**: `main.go`
- **Change**: Extract the table rendering logic into a reusable function (e.g., `printTable()`)
  - Move column width calculation and table formatting into a helper function
  - Accept a slice of films and print in the same table format
- **Rationale**: Avoids code duplication and makes the code more maintainable

### Step 4: Update main.go flow
- **File**: `main.go`
- **Changes**:
  1. Add a menu prompt at the start asking user to choose (1 for shows, 2 for films)
  2. Validate input; show error and exit if not 1 or 2
  3. If user chooses 1: keep existing show display logic
  4. If user chooses 2: call `db.ReadFilms()`, print films table using helper function
- **Rationale**: Provides the user-facing menu as required

### Step 5: Testing (optional enhancement)
- **File**: `shows/shows_test.go` (if needed) or create `db/db_test.go`
- **Change**: Consider adding unit tests for `ReadFilms()` to ensure it correctly parses films.json
- **Rationale**: Ensures robustness; follows existing test patterns in the project

## Acceptance Criteria
- [ ] User sees a menu with options 1 (shows) or 2 (films) on startup
- [ ] Selecting option 1 displays currently watching shows in table format
- [ ] Selecting option 2 displays films from `films.json` in table format
- [ ] Invalid input (not 1 or 2) shows a useful error message and exits gracefully
- [ ] `db.ReadFilms()` successfully reads and parses `films.json`
- [ ] Code follows existing style and patterns (consistent naming)
- [ ] All existing tests still pass
- [ ] New or modified code builds without errors

## Files to Modify
1. `data/data.go` — Add `Film` struct
2. `db/db.go` — Add `ReadFilms()` function
3. `main.go` — Add menu logic, refactor table rendering, call appropriate read function
4. (Optional) `shows/shows_test.go` or new `db/db_test.go` — Add tests for `ReadFilms()`

## Validation Commands
```powershell
go version
go vet ./...
go build ./...
go test ./...
go run .
```

## Notes
- The `films.json` file already exists and contains film data with `name`, `genre`, and `provider` fields
- The `getFullPath()` function in `db.go` handles path resolution for both built binaries and `go run` execution
- Maintain consistency with the existing table formatting style used for shows
- Keep the Film struct simple; no tracking needed (unlike Show which tracks episodes watched)
