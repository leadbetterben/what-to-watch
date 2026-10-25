Repository summary

- This Go program (`what-to-watch`) can be run as an interactive CLI or as an HTTP server, both using the same business logic for consistency.
- Languages / runtimes: Go only (module `what-to-watch`). The repository is small (~25 source files). Key dirs: `db/`, `data/`, `handlers/`, `cmd/cli/`, `cmd/http/`, `plans/`.

Architecture and handler details:

- `handlers/handlers.go` — Core business logic functions:
  - `GetCurrentlyWatchingShows()` — Retrieves currently watching shows
  - `MarkShowWatched(idx)` — Marks a show episode as watched
  - `GetAllFilms()` — Retrieves all films
  - `GetAvailableGenres()` — Retrieves all unique genres from shows
  - `GetUnwatchedShowsByGenre(genre)` — Retrieves unwatched shows for a specific genre
  - `GetSpecialEpisodes()` — Retrieves curated seasonal episode numbers with show and provider context
  - `GetRandomSpecialEpisode` — Retrieves one curated seasonal episode at random
- `cmd/cli/cli.go` — Interactive CLI interface that calls the handlers
- `cmd/http/http.go` — HTTP REST API that calls the same handlers; refactored with `Handler` interface for dependency injection
- `main.go` — Dispatcher: parses flags, routes to CLI or HTTP mode

Top-level facts the agent should trust (no search needed unless instructions are wrong)

- Go version: 1.27.3 (declared in `go.mod` and in the repo CI at `.github/workflows/go.yml`).
- Build toolchain: use the Go toolchain (`go` CLI). Actions/CI uses `actions/setup-go@v4` with `go-version: '1.27.3'`.
- Test behavior: unit tests live under `shows/` and run fast. Tests do not rely on editing `db/currentShows.json` (they use in-memory structs).

Build, run, and test (validated commands)

These sequences were run and validated in a Windows PowerShell environment in this repository root.

1) Confirm toolchain (always run this first):

   `go version`

   Observed: `go1.27.3 version` reports `go version go1.27.3 windows/amd64` in my environment. Always target Go 1.27.3 to match CI.

2) Build (compile all packages):

   `go build ./...`

   Observed: successful build with no errors.

3) Run (CLI mode):

   `go run .`
   `go run . -mode=cli`

   Behavior: runs the CLI and reads/writes `shows.json` via the `db` package. When run with `go run .` during development, `db.getFullPath` resolves the file relative to the source directory.

4) Run (HTTP mode):

   `go run . -mode=http -port=8080`

   Starts HTTP server. Endpoints:
     - `GET /health` — Health check
     - `GET /shows` — Get currently watching shows (JSON) - optional genre param to filter
     - `POST /shows/watch?index=1` — Mark show as watched
     - `GET /films` — Get all films (JSON)
     - `GET /genres` — Get all available genres (JSON)
     - `GET /specials/halloween` — Get curated Halloween episodes (JSON)
     - `GET /specials/christmas` — Get curated Christmas episodes (JSON)
     - `GET /specials/halloween/random` — Get one random Halloween episode (JSON)
     - `GET /specials/christmas/random` — Get one random Christmas episode (JSON)

5) Install (optional):

   `go build` then `go install` — installs the binary to your `GOBIN`/`GOPATH`-based location if desired.

6) Tests:

   `go test ./...`

   Tests use a table-driven approach. Running `go test ./...` in CI is the expected validation step.

   To check package coverage:

   `go test ./... -cover`

   For a detailed local coverage report:

   `go test ./... -coverprofile=coverage.out`
   `go tool cover -func=coverage.out`

7) Lint / formatting:

   No linter config (golangci-lint, etc.) found in repo.

Important environment/workflow notes

- CI: `.github/workflows/go.yml` is the single GitHub Actions workflow. It runs on `push` and `pull_request` to `main` and uses Go 1.27.3. To avoid surprises, match that Go version locally or use the same action in a test run.
- File I/O: `db.ReadShows()` and `db.ReadFilms()` look for their respective JSON files near the executable first, then fall back to the source `db` directory. This means:
  - Built binaries may expect `shows.json` and `films.json` to live next to the binary.
  - `go run .` and tests will find `db/shows.json` and `db/films.json` in the source tree.
  - When writing (`db.UpdateShows`), the code writes atomically (temp file then rename) to the discovered `shows.json` path.
- Tests do not modify on-disk JSON; unit tests use in-memory `data.Show` slices. PRs that change the JSON files should be careful to not accidentally commit runtime-modified files.

Project layout (high-value paths and files to edit)

- `main.go` — dispatcher: parses flags, routes to CLI or HTTP mode
- `handlers/handlers.go` — business logic: `GetCurrentlyWatchingShows()`, `MarkShowWatched()`, `GetAllFilms()`, `GetAvailableGenres()`, `GetUnwatchedShowsByGenre()`, `GetSpecialEpisodes()`, `GetRandomSpecialEpisode()`
- `cmd/cli/cli.go` — CLI interface
- `cmd/http/http.go` — HTTP REST API
- `db/db.go` — functions to read/write `shows.json` and `films.json`, plus `getFullPath` logic.
- `data/data.go` — `Show` and `Film` struct definitions used across the project.
- `db/shows.json` — canonical on-disk data used during `go run .` (do not assume tests use it).
- `plans/` — contains AI-generated implementation plans; not used by code.
- `.github/workflows/go.yml` — CI workflow that must pass for PRs.

Checks run before merge (what CI enforces)

- GitHub Actions `Go` workflow: sets up Go 1.27.3, runs `go build -v ./...` and `go test -race -v ./...`. A PR should pass both build and tests on this workflow.

Quick validation guidance for the agent making changes

- Always run locally before opening a PR: `go build ./...` then `go test ./...`.
- Ensure your Go tool version matches CI (1.27.3). If you cannot install that version locally, run CI-oriented checks in a container or use `actions/setup-go` locally in a disposable runner.
- If the change touches file I/O, double-check `db.getFullPath` semantics: built binaries and `go run` resolve files differently.

Where to search if instructions appear incomplete

- `go.mod` (root) — Go version and module path.
- `.github/workflows/go.yml` — CI config and Go version.
- `main.go`, `db/db.go`, `shows/shows.go`, `data/data.go` — primary behavior and models.

Trust this file first

- Use this document's Go version and commands as authoritative unless you detect a mismatch in the repo (e.g., `go.mod` changed). Only perform a repo-wide search when you believe the instructions are out of date.

Short content snapshot (high-priority snippets)

- `go.mod`: `go 1.27.3`
- `main.go`: dispatcher with CLI/HTTP routing. CLI: menu for shows/films. HTTP: endpoints for shows/films/mark/health.
- `handlers/handlers.go`: `GetCurrentlyWatchingShows()`, `MarkShowWatched()`, `GetAllFilms()`, `GetAvailableGenres()`, `GetUnwatchedShowsByGenre()`
- `db/db.go`: `ReadShows()`, `UpdateShows()`, `ReadFilms()` plus `getFullPath` (see above notes about exe vs source lookup).
- `data/data.go`: `Show` struct (with episode tracking) and `Film` struct (simple name/genre/provider).
- `cmd/http/http.go`: defines `Handler` interface for dependency injection; `defaultHandler` implements it by calling `handlers` package functions.

If anything in this file is inconsistent with the repo state, run `git status` and search the few files listed above before making changes.

End of instructions.
