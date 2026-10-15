# Issue #29 — Validate show positions and episode counts

## Overview

Make current show positions consistent across reading, listing, and progress updates. Reject malformed position data and invalid episode counts with errors that identify the show and field, rather than treating malformed data as successful completion.

This plan covers the current-position and episode-count portion of #29. Seasonal episode metadata validation can be implemented independently.

## Acceptance criteria

- A show has either both `CurrentSeries` and `CurrentEpisode` set or neither; partial positions are reported consistently.
- Current series and episode values are positive and within the corresponding `Episodes` bounds.
- Invalid or non-positive episode counts cannot silently mark a show complete.
- Errors identify the affected show and the invalid field/value.
- Valid first and last episode boundaries, series rollover, and normal completion continue to work.
- Validation does not rewrite `shows.json`.

## Files to change

- `shows/shows.go` — add or share current-position validation and use it when listing and updating shows.
- `shows/shows_test.go` — cover valid boundaries, partial positions, non-positive values, and out-of-range positions/counts.
- `handlers/handlers.go` — inspect the current read/list/update paths and propagate contextual validation errors where needed.
- `db/db.go` — inspect readers and determine where malformed persisted positions should be detected so callers cannot disagree about whether a show is current.
- `cmd/cli/cli.go` and `cmd/http/http.go` — update only if needed to preserve and expose actionable errors through existing interfaces.
- `README.md` — document the current-position rule if it changes user-visible data expectations.

## Implementation steps

1. Trace how `ReadCurrentShows`, `GetCurrentlyWatching`, and `MarkEpisodeWatched` are called by handlers and interfaces. Confirm whether validation belongs in the shared business layer or the storage reader, and apply one policy consistently.
2. Define the position invariant: both current fields are absent for a non-current show; otherwise both are present and positive, the series is within `len(Episodes)`, and the episode is within that series' count.
3. Validate episode counts before progress arithmetic. Reject non-positive counts as malformed data; do not clear the current position or report completion on invalid input.
4. Return contextual errors containing the show name and field/value, while preserving existing completion and rollover behavior for valid records.
5. Add focused unit coverage for valid boundaries, each partial-field combination, zero/negative series and episode values, invalid episode counts, and values beyond the available series/episode bounds. Use isolated fixtures and do not modify the canonical `db/shows.json`.
6. Update user-facing documentation if the invariant or error behavior changes the documented data format.

## Validation commands

- `go build ./...`
- `go test ./...`
- `go vet ./...`

## Scope notes

- Seasonal category/series/episode validation is a separate, independent part of #29 and is not included here.
- Keep changes focused; do not alter seasonal records or rewrite local runtime JSON data.

---

*Generated plan for Issue #29 (current positions and episode counts).*
