# Issues #12 and #13: Add Halloween and Christmas episode listings

## Overview

Allow users to retrieve all known Halloween or Christmas episodes for shows in
the collection. Store the curated episode metadata on each show, then expose
separate handler functions and HTTP endpoints that return a flat list of
matching episodes with enough show context to identify where each episode
belongs.

The episode list is user-curated data. Do not attempt to infer special episodes
from the episode counts in `Show.Episodes`.

## Data format

Add a dedicated category type with constants for the two supported categories,
then use that type on a reusable special-episode type in `data/data.go`:

```go
type SpecialCategory string

const (
 SpecialCategoryHalloween SpecialCategory = "halloween"
 SpecialCategoryChristmas SpecialCategory = "christmas"
)

type SpecialEpisode struct {
 Category SpecialCategory `json:"category"`
 Series   int             `json:"series"`
 Episode  int             `json:"episode"`
}

type Show struct {
 // existing fields...
 SpecialEpisodes []SpecialEpisode `json:"specialEpisodes,omitempty"`
}
```

Represent a curated episode as an entry nested in its parent show in
`db/shows.json`:

```json
{
  "name": "Example Show",
  "genre": "comedy",
  "episodes": [6, 8],
  "provider": "Example Service",
  "specialEpisodes": [
    {
      "category": "halloween",
      "series": 2,
      "episode": 5
    },
    {
      "category": "christmas",
      "series": 3,
      "episode": 4
    }
  ]
}
```

- Use the two category constants when creating or querying entries; their JSON
  values are exactly `halloween` and `christmas`.
- Store only the special category, series number, and episode number. The
  purpose is to locate and watch the episode, so a title field is unnecessary.
- If the curated data identifies a special without a normal series or episode
  number, revise the model to use optional number pointers before importing it;
  do not invent numbers.
- Preserve multiple specials for a show and do not deduplicate solely by title.
  The show name and provider identify the parent record, since duplicate show
  names can exist on different providers.

## Implementation steps

### 1. Confirm and map the curated list

- Use the user's collated list as the source of truth for all records in issues
  #12 and #13.
- Map each source row to a show by name and provider, then map category, series,
  and episode number.
- Record unmatched or ambiguous rows for resolution instead of silently
  skipping them or assigning them to a similarly named show.
- Check whether any listed episode belongs to both categories and preserve both
  classifications if the source says so.

### 2. Add data types and local data

- Add `data.SpecialEpisode` and `Show.SpecialEpisodes` in `data/data.go`.
- Add the curated entries to the corresponding show objects in local
  `db/shows.json`, preserving all existing show fields and order.
- `db/*.json` is ignored by Git by design. Do not force-add or otherwise track
  `db/shows.json`; commit the schema, functions, tests, docs, and plan only.
- Validate the resulting JSON and confirm each imported row maps to exactly
  the intended show and episode.

### 3. Add reusable query logic and handler functions

- Add a query function in `shows/shows.go` that takes all shows and a requested
  category, flattens matching special entries, and attaches parent show name
  and provider to each returned record.
- Define a result type for the flattened response, such as
  `data.ShowSpecialEpisode` containing show name, provider, category, series,
  and episode number.
- Accept `data.SpecialCategory` in the query function and validate it against
  the two constants. Unsupported values decoded from JSON or passed by callers
  should produce a clear error rather than an empty successful response.
- Add handler functions for Halloween and Christmas, or a shared category
  handler with thin `GetHalloweenEpisodes()` and `GetChristmasEpisodes()`
  wrappers. Read all shows through `db.ReadAllShows()`.

### 4. Add HTTP endpoints

- Extend the HTTP `Handler` interface and `defaultHandler` with the new
  functions.
- Add `GET /specials/halloween` and `GET /specials/christmas` endpoints.
- Return JSON arrays of flattened episode records. Return an empty JSON array
  when there are no matches, and use the existing error response conventions
  for read or query failures.
- Reject unsupported methods with the existing method-error response.

### 5. Test and document

- Add unit tests for category filtering, multiple specials on one show,
  duplicate show names with different providers, empty results, and invalid
  categories.
- Add handler tests that verify data flows from the all-shows reader into the
  category results.
- Add table-driven HTTP tests for both endpoints, method errors, handler errors,
  and JSON response shape.
- Update `README.md` with endpoint descriptions and example responses.
- Update `AGENTS.md` with the new special-episode model and retrieval functions.
- Do not update `copilot-instructions.md` unless the file exists and contains
  relevant repository guidance.

## Acceptance criteria

- Halloween and Christmas metadata is represented as structured entries on
  the relevant show records, with no inferred or invented episode details.
- Halloween retrieval returns every and only curated Halloween episode;
  Christmas retrieval returns every and only curated Christmas episode.
- Each result includes the parent show's name and provider, plus category,
  series number, and episode number. No episode title is required.
- Both HTTP GET endpoints return the documented JSON format and follow the
  existing error and method handling conventions.
- The ignored `db/shows.json` remains untracked.
- `go vet ./...`, `go build ./...`, `go test ./...`, and
  `go test ./... -cover` pass.
