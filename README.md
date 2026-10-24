# what-to-watch

Finds a series, TV show, or film to watch

This program can be run as an interactive CLI or as an HTTP server, both using the same business logic for consistency.

**Shows**: View and update TV shows you are currently watching, including provider, current series, and episode. Mark episodes as watched.

**Films**: View your film collection with genre and provider information.

Show data is stored in `db/shows.json` as one array. Shows with
`currentSeries`/`currentEpisode` are currently watching; shows with
`rewatch: true` are queued for rewatch; records with neither are unwatched.

The `plans/` directory contains AI-generated plans for implementations.

See GitHub Issues for future plans. New Issues and Pull Requests are welcome.

## Usage

### CLI Mode (Default)

Run the program in interactive command-line mode:

```bash
go run .
# or with explicit mode flag
go run . -mode=cli
```

This launches an interactive menu:

``` text
What would you like to view?
1. Currently watching shows
2. Films
3. Shows by genre
4. Halloween episodes
5. Christmas episodes
6. Random Halloween episode
7. Random Christmas episode
Enter your choice (1, 2, 3, 4, 5, 6, or 7):
```

Select options 1–3 for the existing show, film, and genre views, 4 or 5 to list curated Halloween or Christmas episodes, or 6 or 7 to get one at random.

### HTTP Mode

Start an HTTP server to interact with the API:

```bash
go run . -mode=http -port=8080
```

The server listens on `127.0.0.1` only, so it is reachable from this computer and not from other machines on the network. The port can be customized with the `-port` flag (default: 8080).

HTTP mode is intended for local use. Remote deployment and access from web or app clients are not currently supported. The server applies a 5-second header read timeout, a 10-second read timeout, a 30-second write timeout, and a 60-second idle connection timeout.

Show progress updates are serialized across concurrent requests within one running process. Running multiple application processes that write to the same `shows.json` file is not supported.

#### Available Endpoints

- `GET /health` — Health check
- `GET /shows` — Get currently watching shows (JSON) - optional genre param to filter
- `POST /shows/watch?index=1` — Mark show as watched
- `GET /films` — Get all films (JSON)
- `GET /genres` — Get all available genres (JSON)
- `GET /specials/halloween` — Get curated Halloween episodes across shows
- `GET /specials/christmas` — Get curated Christmas episodes across shows
- `GET /specials/halloween/random` — Get one random Halloween episode
- `GET /specials/christmas/random` — Get one random Christmas episode

#### Example API Calls

```bash
# Check server health
curl http://127.0.0.1:8080/health

# Get currently watching shows
curl http://127.0.0.1:8080/shows

# Get unwatched shows in Drama genre
curl http://127.0.0.1:8080/shows?genre=drama

# Mark show #1 as watched
curl -X POST http://127.0.0.1:8080/shows/watch?index=1

# Get all films
curl http://127.0.0.1:8080/films

# Get available genres
curl http://127.0.0.1:8080/genres

# Get curated Halloween episodes
curl http://127.0.0.1:8080/specials/halloween

# Get curated Christmas episodes
curl http://127.0.0.1:8080/specials/christmas

# Get one random Halloween episode
curl http://127.0.0.1:8080/specials/halloween/random

# Get one random Christmas episode
curl http://127.0.0.1:8080/specials/christmas/random

```

## Architecture

The program uses consistent handler functions that can be called by either interface:

- **`handlers/handlers.go`** — Core business logic functions:
  - `GetCurrentlyWatchingShows()` — Retrieves currently watching shows
  - `MarkShowWatched(idx)` — Marks a show episode as watched
  - `GetAllFilms()` — Retrieves all films
  - `GetAvailableGenres()` — Retrieves all unique genres from shows
  - `GetUnwatchedShowsByGenre(genre)` — Retrieves unwatched shows for a specific genre
  - `GetHalloweenEpisodes()` / `GetChristmasEpisodes()` — Retrieves curated seasonal episode numbers with show and provider context
  - `GetRandomHalloweenEpisode()` / `GetRandomChristmasEpisode()` — Retrieves one curated seasonal episode at random
- **`cmd/cli/cli.go`** — Interactive CLI interface that calls the handlers
- **`cmd/http/http.go`** — HTTP REST API that calls the same handlers

Both modes use the same underlying business logic, ensuring consistency across interfaces.

## Build & Run

1. Run `go build`
2. Run `go install`
3. `what-to-watch` should be available to run from any directory

## Testing

Run the full test suite:

```bash
go test ./...
```

Check package coverage:

```bash
go test ./... -cover
```

For a detailed coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```
