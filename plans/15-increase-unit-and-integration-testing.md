# Issue #15: Increase Unit and Integration Testing

## Overview

GitHub issue #15 asks for more unit and integration testing so future feature work can be made with higher confidence. The current test suite already covers most pure show business logic and a large portion of HTTP handler behavior, but several important boundaries are still untested:

- CLI output formatting and user input flows
- `handlers` package orchestration between `db` and `shows`
- `db` JSON read/write and path resolution behavior
- HTTP routing/server integration beyond direct handler method tests
- `main.go` mode selection and invalid-mode behavior

Current baseline from `go test ./... -cover`:

- `what-to-watch`: 0.0%
- `what-to-watch/cmd/cli`: 0.0%
- `what-to-watch/cmd/http`: 73.1%
- `what-to-watch/data`: no test files
- `what-to-watch/db`: 0.0%
- `what-to-watch/handlers`: 0.0%
- `what-to-watch/shows`: 94.2%

## Goals

- Add targeted unit tests around packages with 0% coverage.
- Add integration tests that exercise real package wiring without mutating canonical JSON files.
- Preserve current behavior and keep production refactors minimal.
- Avoid external dependencies; use Go standard library test tools.

## Non-Goals

- Do not redesign CLI or HTTP APIs as part of this issue.
- Do not add broad snapshot/golden-file infrastructure unless it stays simple and local.
- Do not make tests depend on modifying `db/currentShows.json`, `db/shows.json`, or `db/films.json`.

## Proposed Implementation Steps

### 1. Add a Testability Layer for Data File Paths

Files:

- `db/db.go`
- new `db/db_test.go`

Plan:

- Add a package-level path resolver variable in `db`, defaulting to the existing `getFullPath` behavior.
- Use that resolver from `readFile` and `WriteCurrentShows`.
- In tests, temporarily override the resolver with `t.Cleanup` to point at `t.TempDir()` fixtures.
- Keep `getFullPath` behavior unchanged for production and existing commands.

Unit tests:

- `ReadShows` parses a valid `shows.json`.
- `ReadCurrentShows` parses a valid `currentShows.json`.
- `ReadFilms` parses a valid `films.json`.
- Each reader returns an error for missing files.
- Each reader returns an error for invalid JSON.
- `WriteCurrentShows` writes indented JSON to the resolved path.
- `WriteCurrentShows` replaces an existing file without leaving temp files behind.

Acceptance criteria:

- `db` tests never touch repository JSON files.
- Tests verify both successful JSON parsing and common failure paths.

### 2. Add Handler Integration Tests with Temporary Data

Files:

- `handlers/handlers.go`
- new `handlers/handlers_test.go`

Plan:

- Reuse the `db` test path override pattern to run handlers against temporary JSON fixtures.
- Treat these as package-level integration tests because they exercise `handlers`, `db`, and `shows` together.
- Use small fixtures with clear expected state transitions.

Tests:

- `GetCurrentlyWatchingShows` returns only currently watching shows with computed series/episode display fields.
- `MarkShowWatched` increments an episode and persists the update.
- `MarkShowWatched` handles episode/series completion and persists the update.
- `GetAllFilms` returns films from the temp fixture.
- `GetAvailableGenres` returns unique genres from temp `shows.json`.
- `GetUnwatchedShowsByGenre` returns unwatched shows for a genre.
- Read/write failures are surfaced with handler-specific context.

Acceptance criteria:

- Tests validate persisted data after `MarkShowWatched`.
- Handler tests fail if the handler layer stops using shared business logic.

### 3. Add CLI Formatting Unit Tests

Files:

- new `cmd/cli/format_test.go`

Plan:

- Test pure formatting helpers before adding tests for interactive flows.
- Use exact string expectations for small tables where alignment matters.
- Cover empty-state text because it is user-facing behavior.

Tests:

- `formatShowsTable` with multiple rows and varying column widths.
- `formatShowsTable` empty state.
- `formatShowsByGenreTable` with data.
- `formatShowsByGenreTable` empty state.
- `formatFilmsTable` with data.
- `formatFilmsTable` empty state.

Acceptance criteria:

- CLI table regressions are caught without needing to run the interactive CLI.

### 4. Add CLI Flow Tests with Injectable Handler Functions

Files:

- `cmd/cli/cli.go`
- new `cmd/cli/cli_test.go`

Plan:

- Introduce small package-level function variables in `cmd/cli` for calls currently made directly to `handlers`.
- Default them to the existing handler functions.
- In tests, override them with stubs and restore with `t.Cleanup`.
- Capture `os.Stdout` and provide `bufio.Reader` input from strings.

Tests:

- `viewShows` displays shows and cancels on `0`.
- `viewShows` rejects non-numeric input.
- `viewShows` calls `MarkShowWatched` for valid numeric input and prints completed/non-completed messages.
- `viewFilms` displays films.
- `viewShowsByGenre` displays available genres, handles cancel/invalid selection, and displays selected genre results.
- Error paths print a concise `Error:` line when handler functions fail.

Acceptance criteria:

- Interactive paths are covered without reading or writing real data files.
- Production CLI behavior remains unchanged.

### 5. Expand HTTP Tests into Router-Level Integration Tests

Files:

- `cmd/http/http.go`
- `cmd/http/http_test.go`

Plan:

- Add a `routes()` or `ServeHTTP` helper on `Server` that returns an `http.Handler` using a local `http.ServeMux`.
- Have `Start` call the same route builder and pass it to `http.ListenAndServe`.
- Keep existing direct handler tests.
- Add `httptest.NewServer` integration tests using `NewServerWithHandler`.

Tests:

- `GET /shows` routes to currently watching shows.
- `GET /shows?genre=Drama` routes to genre filtering.
- `POST /shows/watch?index=1` routes to mark watched.
- `GET /films`, `GET /genres`, and `GET /health` return JSON responses.
- Unsupported methods return `405`.
- Unknown routes return `404`.

Acceptance criteria:

- Route registration is tested without binding a real port.
- HTTP tests use the same mux that production uses.

### 6. Add Main Package Tests

Files:

- `main.go`
- new `main_test.go`

Plan:

- Extract mode dispatch into a small function, for example `run(mode string, port int) int`, where return value is an exit code.
- Keep `main()` responsible for flag parsing and `os.Exit`.
- Use injectable CLI and HTTP start functions so tests can verify dispatch without starting a server.

Tests:

- `run("cli", 8080)` calls CLI runner and returns `0`.
- `run("http", 9090)` starts HTTP mode with the requested port and returns `0` on success.
- HTTP start errors return non-zero.
- Invalid mode returns non-zero and writes the current invalid-mode message.

Acceptance criteria:

- Main dispatch is covered without launching interactive input or a long-running server.
- The public command-line behavior remains the same.

### 7. Add Coverage Reporting to Validation Workflow

Files:

- `README.md`
- `AGENTS.md`
- `.github/copilot-instructions.md` if present

Plan:

- Document the coverage command alongside existing validation commands:

```sh
go test ./... -cover
```

- Optionally add a local deeper coverage command for contributors:

```sh
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Acceptance criteria:

- Documentation tells contributors how to check coverage locally.
- Agent instructions remain aligned with repository workflow.

## Suggested Order

1. Add `db` testability layer and `db` unit tests.
2. Add handler integration tests using temporary JSON fixtures.
3. Add CLI formatting tests.
4. Add CLI flow tests with injectable handler functions.
5. Refactor HTTP routing minimally and add router-level integration tests.
6. Extract main dispatch and add main package tests.
7. Update documentation with coverage commands.

## Validation Commands

Run after implementation:

```sh
go vet ./...
go build ./...
go test ./...
go test ./... -cover
```

For detailed coverage:

```sh
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Final Acceptance Criteria

- New tests cover `db`, `handlers`, `cmd/cli`, and `main`.
- HTTP tests include route-level integration coverage.
- Tests do not mutate canonical JSON data under `db/`.
- All validation commands pass.
- Overall package coverage improves materially from the current baseline.
- Any production changes are limited to small testability seams and shared route/dispatch helpers.
