# Issue #4: Display New Shows Based on Genre

## Overview

Users should have the option to see shows they have not yet watched, filtered by genre. The feature should:

1. Present a list of available genres
2. Allow the user to select a genre
3. Display all shows from that genre that have not been watched

## Current State

- `data/data.go` defines `Show` and `Film` structs
- `shows/shows.go` contains business logic for show operations
- `handlers/handlers.go` provides handler functions for CLI and HTTP
- `cmd/cli/cli.go` implements the interactive CLI menu
- `cmd/http/http.go` implements HTTP REST endpoints

## Implementation Steps

### 1. Analyze Data Structure

- Review `data/data.go` to understand the `Show` struct
- Identify how genres are stored in shows
- Determine how watched status is tracked

### 2. Add Business Logic

- In `shows/shows.go`, add function to extract unique genres from all shows
- Add function to filter shows by genre and watched status
- Add unit tests in `shows/shows_test.go` for both functions

### 3. Add Handler Functions

- In `handlers/handlers.go`, add handler to get available genres
- Add handler to get shows filtered by genre (excluding watched)
- Ensure handlers follow the existing pattern

### 4. Update CLI Interface

- In `cmd/cli/cli.go`, add menu option to "View shows by genre"
- Implement interactive flow: display genres → user selects → display shows

### 5. Add HTTP Endpoints

- In `cmd/http/http.go`, add `GET /genres` endpoint to list available genres
- Add `GET /shows/genre?name=<genre_name>` endpoint to get shows by genre
- Add corresponding tests in `cmd/http/http_test.go`

### 6. Test & Validate

- Run unit tests: `go test ./...`
- Test CLI mode: `go run .`
- Test HTTP mode: `go run . -mode=http -port=8080`
- Verify both CLI and HTTP return correct filtered results

### Update Documentation

- Update README.md intended for a human audience
- Update AGENTS.md intended for an AI/LLM audience
- Update .github/copilot-instructions.md intended for an AI/LLM audience

## Acceptance Criteria

- ✓ CLI displays genre selection menu
- ✓ CLI displays all unwatched shows from selected genre
- ✓ HTTP endpoint `/genres` returns list of all genres
- ✓ HTTP endpoint `/shows/genre?name=<genre_name>` returns unwatched shows for that genre
- ✓ All unit tests pass
- ✓ No breaking changes to existing functionality
- ✓ Documentation updated with new functionality
