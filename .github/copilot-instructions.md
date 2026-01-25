Repository summary

- This Go program (`what-to-watch`) can be run as an interactive CLI or as an HTTP server, both using the same business logic for consistency.
- Languages / runtimes: Go only (module `what-to-watch`). The repository is small (~25 source files). Key dirs: `db/`, `data/`, `shows/`, `handlers/`, `cmd/cli/`, `cmd/http/`, `plans/`.

Architecture and handler details:

- `handlers/handlers.go` — Core business logic functions:
  - `GetCurrentlyWatchingShows()` — Retrieves currently watching shows
  - `MarkShowWatched(idx)` — Marks a show episode as watched
  - `GetAllFilms()` — Retrieves all films
- `cmd/cli/cli.go` — Interactive CLI interface that calls the handlers
- `cmd/http/http.go` — HTTP REST API that calls the same handlers
- `main.go` — Dispatcher: parses flags, routes to CLI or HTTP mode

Top-level facts the agent should trust (no search needed unless instructions are wrong)

- Go version: 1.25.4 (declared in `go.mod` and in the repo CI at `.github/workflows/go.yml`).
- Build toolchain: use the Go toolchain (`go` CLI). Actions/CI uses `actions/setup-go@v4` with `go-version: '1.25.4'`.
- Test behavior: unit tests live under `shows/` and run fast. Tests do not rely on editing `db/shows.json` (they use in-memory structs).

Build, run, and test (validated commands)

These sequences were run and validated in a Windows PowerShell environment in this repository root.

1) Confirm toolchain (always run this first):

   `go version`

   Observed: `go version go1.25.4 windows/amd64` in my environment. Always target Go 1.25.4 to match CI.

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
     - `GET /api/shows` — Get currently watching shows (JSON)
     - `POST /api/shows/mark?index=1` — Mark show as watched
     - `GET /api/films` — Get all films (JSON)

5) Install (optional):

   `go build` then `go install` — installs the binary to your `GOBIN`/`GOPATH`-based location if desired.

6) Tests:

   `go test ./...`

   Observed results: `ok what-to-watch/shows 3.659s` (tests pass locally). The `data` and `db` packages have no tests. Running `go test ./...` in CI is the expected validation step.

7) Lint / formatting:

   No linter config (golangci-lint, etc.) found in repo.

Important environment/workflow notes

- CI: `.github/workflows/go.yml` is the single GitHub Actions workflow. It runs on `push` and `pull_request` to `main` and uses Go 1.25.4. To avoid surprises, match that Go version locally or use the same action in a test run.
- File I/O: `db.ReadShows()` and `db.ReadFilms()` look for their respective JSON files near the executable first, then fall back to the source `db` directory. This means:
  - Built binaries may expect `shows.json` and `films.json` to live next to the binary.
  - `go run .` and tests will find `db/shows.json` and `db/films.json` in the source tree.
  - When writing (`db.WriteShows`), the code writes atomically (temp file then rename) to the discovered `shows.json` path.
- Tests do not modify on-disk JSON; unit tests use in-memory `data.Show` slices. PRs that change the JSON files should be careful to not accidentally commit runtime-modified files.

Project layout (high-value paths and files to edit)

- `main.go` — dispatcher: parses flags, routes to CLI or HTTP mode
- `handlers/handlers.go` — business logic: `GetCurrentlyWatchingShows()`, `MarkShowWatched()`, `GetAllFilms()`
- `cmd/cli/cli.go` — CLI interface
- `cmd/http/http.go` — HTTP REST API
- `db/db.go` — functions to read/write `shows.json` and `films.json`, plus `getFullPath` logic.
- `data/data.go` — `Show` and `Film` struct definitions used across the project.
- `shows/shows.go` — business logic: `GetCurrentlyWatching`, `MarkEpisodeWatched`.
- `shows/shows_test.go` — unit tests for `shows` package (good examples of expected behavior).
- `db/shows.json` — canonical on-disk data used during `go run .` (do not assume tests use it).
- `plans/` — contains AI-generated implementation plans; not used by code.
- `.github/workflows/go.yml` — CI workflow that must pass for PRs.

Checks run before merge (what CI enforces)

- GitHub Actions `Go` workflow: sets up Go 1.25.4, runs `go build -v ./...` and `go test -v ./...`. A PR should pass both build and tests on this workflow.

Quick validation guidance for the agent making changes

- Always run locally before opening a PR: `go build ./...` then `go test ./...`.
- Ensure your Go tool version matches CI (1.25.4). If you cannot install that version locally, run CI-oriented checks in a container or use `actions/setup-go` locally in a disposable runner.
- If the change touches file I/O, double-check `db.getFullPath` semantics: built binaries and `go run` resolve files differently.
- Unit tests live in `shows/` — read `shows/shows_test.go` to understand expected business behavior. Use those tests as a model for new tests.

Where to search if instructions appear incomplete

- `go.mod` (root) — Go version and module path.
- `.github/workflows/go.yml` — CI config and Go version.
- `main.go`, `db/db.go`, `shows/shows.go`, `data/data.go` — primary behavior and models.
- `shows/shows_test.go` — canonical test expectations.

Trust this file first

- Use this document's Go version and commands as authoritative unless you detect a mismatch in the repo (e.g., `go.mod` changed). Only perform a repo-wide search when you believe the instructions are out of date.

Short content snapshot (high-priority snippets)

- `go.mod`: `go 1.25.4`
- `main.go`: dispatcher with CLI/HTTP routing. CLI: menu for shows/films. HTTP: endpoints for shows/films/mark/health.
- `handlers/handlers.go`: `GetCurrentlyWatchingShows()`, `MarkShowWatched()`, `GetAllFilms()`, `FormatShowsTable()`, `FormatFilmsTable()`
- `db/db.go`: `ReadShows()`, `WriteShows()`, `ReadFilms()` plus `getFullPath` (see above notes about exe vs source lookup).
- `data/data.go`: `Show` struct (with episode tracking) and `Film` struct (simple name/genre/provider).
- `shows/shows.go`: contains `GetCurrentlyWatching` and `MarkEpisodeWatched` business logic (tests in `shows/shows_test.go`).

If anything in this file is inconsistent with the repo state, run `git status` and search the few files listed above before making changes.

End of instructions.
