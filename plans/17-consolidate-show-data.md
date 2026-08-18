# Issue #17: Consolidate show data into one file

## Overview

Replace the three show data files (`db/shows.json`, `db/currentShows.json`, and
`db/rewatchShows.json`) with one canonical `db/shows.json` file. The file will
store every show in one flat array, while the `db` package continues to expose
focused functions that return only the requested category.

The category is derived from fields on each `data.Show`: current-series and
current-episode values identify currently watched shows, and a rewatch flag
identifies rewatch entries. A show with neither status is unwatched.

## Target data format

Keep `db/shows.json` as a single JSON array, but add a rewatch marker to the
`Show` type:

```json
[
  {
    "name": "Example",
    "genre": "comedy",
    "episodes": [6],
    "provider": "BBC iPlayer",
    "rewatch": true
  }
]
```

- Merge every existing entry from all three files into the array.
- Add a `Rewatch bool` field to `data.Show` with the struct tag
  `json:"rewatch,omitempty"`, setting it only on entries migrated from
  `rewatchShows.json`.
- Preserve `currentSeries` and `currentEpisode` on entries migrated from
  `currentShows.json`; leave those pointers unset for unwatched and rewatch
  entries unless the data explicitly requires both statuses.

## Implementation steps

### 1. Model the consolidated document

- In `data/data.go`, add the `Rewatch` boolean field with an `omitempty` JSON
  tag; keep `Show` as the single persisted record type.
- Document the classification rule in code: current means both current-series
  and current-episode are set, rewatch means `Rewatch == true`, and unwatched
  means neither condition is true. Define precedence if a record has both
  current pointers and the rewatch flag.

### 2. Update database reads and writes

- In `db/db.go`, add one private helper that reads and unmarshals the single
  `shows.json` array. Retain the existing path
  resolution behaviour, so a built binary can still use a data file beside the
  executable and `go run` still falls back to `db/shows.json` in source.
- Expose these distinct read functions, each returning `[]data.Show`:
  - `ReadAllShows()` - returns the complete array.
  - `ReadUnwatchedShows()` - filters records with neither current pointers nor
    the rewatch flag.
  - `ReadCurrentShows()` - filters records with the current pointers set.
  - `ReadRewatchShows()` - filters records with `Rewatch == true`.
- Replace the ambiguous existing `ReadShows()` call sites with
  `ReadUnwatchedShows()` and remove `ReadShows()` once all in-repository callers
  have migrated.
- Change `WriteCurrentShows()` to read the complete array, replace the records
  classified as currently watching with the supplied current records, and
  atomically write the complete array back to `shows.json`. Preserve unwatched
  and rewatch records. Define and test a stable matching key (for example,
  name plus provider) for updates and additions because `Show` has no ID.
- Factor the atomic JSON write into a private helper if it avoids duplication;
  continue creating the temporary file in the target directory before rename.

### 3. Migrate application callers

- Update `handlers/handlers.go` so genre and unwatched-show operations call
  `db.ReadUnwatchedShows()`; keep currently-watching operations on
  `db.ReadCurrentShows()`.
- Inspect CLI and HTTP layers for direct database access or terminology that
  assumes separate JSON files, and update any affected tests or messages.
- Do not add a user-facing rewatch feature as part of this issue; provide the
  database reader so it is available to that future work.

### 4. Migrate fixtures and tests

- Update `db/db_test.go` fixtures to write one flat `shows.json` array and test
  all four reader functions, including classification boundaries and complete
  `ReadAllShows()` results.
- Update missing-file and invalid-JSON table tests to exercise the consolidated
  file through each reader.
- Update write tests to verify `WriteCurrentShows()` replaces current records,
  preserves unwatched and rewatch records, produces indented valid JSON,
  replaces an existing file atomically, and cleans up its temporary file.
- Update `handlers/handlers_test.go` fixtures to use the new JSON document;
  retain coverage for marking an episode watched and verify it leaves
  unwatched and rewatch records intact.

### 5. Replace repository data and documentation

- Build the new flat `db/shows.json` by merging the three current files and
  setting `rewatch: true` on migrated rewatch entries, then remove
  `db/currentShows.json` and
  `db/rewatchShows.json` after the combined file has been verified.
- Update `README.md` with the consolidated data-file format if it documents
  show storage or deployment data files.
- Update `AGENTS.md` to describe `db/shows.json` as the canonical show data
  file and revise its database/path-resolution and write notes.
- Update `copilot-instructions.md` only if it contains the old file names or
  database workflow.

## Acceptance criteria

- `db/shows.json` is the only on-disk show data file; it is a flat array
  containing every entry from the former files.
- `ReadAllShows`, `ReadUnwatchedShows`, `ReadCurrentShows`, and
  `ReadRewatchShows` each return exactly their documented data.
- Marking a current episode persists the updated current records and does not
  discard unwatched or rewatch data.
- The CLI and HTTP behaviour for existing currently-watching and unwatched
  flows remains unchanged.
- Database and handler tests use only the consolidated file format and pass.
- `go vet ./...`, `go build ./...`, `go test ./...`, and
  `go test ./... -cover` succeed.
