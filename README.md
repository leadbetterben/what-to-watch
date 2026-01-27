# what-to-watch

Finds a series, TV show, or film to watch

This program can be run as an interactive CLI or as an HTTP server, both using the same business logic for consistency.

**Shows**: View and update TV shows you are currently watching, including provider, current series, and episode. Mark episodes as watched.

**Films**: View your film collection with genre and provider information.

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
Enter your choice (1 or 2):
```

Select option 1 to view and update currently watching shows, or option 2 to view your films collection.

### HTTP Mode

Start an HTTP server to interact with the API:

```bash
go run . -mode=http -port=8080
```

The port can be customized with the `-port` flag (default: 8080).

#### Available Endpoints

- `GET /health` — Health check
- `GET /shows` — Get currently watching shows (JSON)
- `POST /shows/watch?index=1` — Mark show as watched
- `GET /films` — Get all films (JSON)

#### Example API Calls

```bash
# Check server health
curl http://localhost:8080/health

# Get currently watching shows
curl http://localhost:8080/shows

# Mark show #1 as watched
curl -X POST http://localhost:8080/shows/mark?index=1

# Get all films
curl http://localhost:8080/films
```

## Architecture

The program uses consistent handler functions that can be called by either interface:

- **`handlers/handlers.go`** — Core business logic functions:
  - `GetCurrentlyWatchingShows()` — Retrieves currently watching shows
  - `MarkShowWatched(idx)` — Marks a show episode as watched
  - `GetAllFilms()` — Retrieves all films
- **`cmd/cli/cli.go`** — Interactive CLI interface that calls the handlers
- **`cmd/http/http.go`** — HTTP REST API that calls the same handlers

Both modes use the same underlying business logic, ensuring consistency across interfaces.

## Build & Run

1. Run `go build`
2. Run `go install`
3. `what-to-watch` should be available to run from any directory
