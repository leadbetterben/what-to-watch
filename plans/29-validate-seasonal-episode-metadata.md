# Issue #29 — Validate seasonal episode metadata

## Overview

Validate each curated `specialEpisodes` entry when show data is read. Accept only supported seasonal categories and positive series and episode numbers. When the parent show's `episodes` counts are available, reject seasonal positions beyond those bounds. Errors should identify the show and field, and validation must not rewrite the source data.

This plan covers the seasonal metadata portion of #29. Current watching-position and episode-count validation is handled separately.

## Acceptance criteria

- Only `halloween` and `christmas` are accepted as seasonal categories.
- Series and episode numbers must be positive.
- When episode counts are available for the show, the seasonal series and episode are within those counts.
- Invalid entries return actionable errors identifying the show and invalid field/value.
- Valid curated entries remain unchanged and are returned through the existing seasonal list and random selection paths.
- Reading invalid metadata does not rewrite `shows.json` or silently discard the entry.
- Tests cover valid first/last boundaries, unsupported categories, zero/negative values, and out-of-range series/episode values.

## Files to change

- `data/data.go` — add shared validation for `SpecialEpisode` values in the context of their parent show.
- `db/db.go` — validate seasonal metadata in the common show reader so handlers and commands see consistent errors.
- `db/db_test.go` — use isolated fixtures to test malformed persisted seasonal metadata and confirm valid data still loads.
- `shows/shows_test.go` — cover seasonal listing behavior for valid boundaries and ensure curated values are preserved.
- `handlers/handlers_test.go` — verify storage validation errors propagate through seasonal handlers, if not already covered.
- `README.md` — document supported categories and metadata requirements if the JSON data format is described.

## Implementation steps

1. Trace seasonal read paths from `db.ReadAllShows` through handlers and the existing list/random functions. Keep validation in the shared persisted-data path so all callers use the same rule.
2. Define a validator that checks the category against `SpecialCategoryHalloween` and `SpecialCategoryChristmas`, and checks that series and episode values are positive.
3. If the parent show's `Episodes` slice provides counts, validate the series against the slice and the episode against that series' count. Decide and document behavior for missing counts; do not invent a bound when source data cannot establish one.
4. Return errors with the parent show name and field/value. Do not mutate or write the source JSON during validation.
5. Add tests for valid lower and upper bounds, each unsupported or non-positive value, series and episode overflow, and error propagation from persisted reads.
6. Confirm existing valid Halloween/Christmas entries and random/list behavior are unchanged, and update README data-format guidance where applicable.

## Validation commands

- `go build ./...`
- `go test ./...`
- `go vet ./...`

## Scope notes

- This work is independent from current watching-position validation and should remain based on `main`.
- Preserve curated episode numbers; report invalid records with context instead of silently correcting them.
- Do not modify local runtime `db/shows.json` as part of implementation or tests.

---

*Generated plan for Issue #29 (seasonal episode metadata).*
