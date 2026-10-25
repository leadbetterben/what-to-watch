# Issue #30 — Return genres in deterministic sorted order

## Overview

Make `handlers.GetAvailableGenres` return unique genres in lexical order so the CLI menu numbering and `/genres` API output are stable for identical show data. Genre aggregation now lives in `handlers/handlers.go`; there is no `shows/` package in the current repository layout.

## Acceptance criteria

- `handlers.GetAvailableGenres` returns unique genre strings in ascending lexical order.
- Repeated calls with identical inputs yield the same order.
- CLI genre menu numbers follow that order, and HTTP `/genres` returns it unchanged.
- Function comments and tests clearly state and verify the ordering contract.
- Empty input and duplicate genres continue to behave correctly.

## Files to change

- `handlers/handlers.go` — sort the unique genre slice before returning it and update the function comment.
- `handlers/handlers_test.go` — assert an explicitly sorted result, duplicate removal, empty genres behavior, and empty input behavior.
- `cmd/cli/cli_test.go` — adjust the genre menu case to use unsorted handler input and assert the expected lexical numbering through the shared handler behavior.
- `cmd/http/http_test.go` — confirm `/genres` preserves the order returned by `GetAvailableGenres`; avoid sorting in the HTTP layer.
- `README.md` — document sorted genre output only if the API/menu documentation should state this user-visible contract.

## Implementation steps

1. Inspect `handlers.GetAvailableGenres` and its tests; keep sorting in the shared handler rather than duplicating it in CLI or HTTP layers.
2. Sort the deduplicated genre slice using Go's lexical string order and update the function comment.
3. Add focused handler tests with intentionally unsorted and repeated genres, asserting exact output order and empty-input behavior.
4. Review CLI and HTTP tests to ensure they consume the returned order without reordering it; strengthen the CLI assertion if necessary to make stable numbering clear.
5. Update README only if needed to make the stable output contract discoverable to users.

## Validation commands

- `go test ./shows ./cmd/cli ./cmd/http`
- `go build ./...`
- `go test ./...`
- `go vet ./...`

## Scope notes

- Use ordinary lexical ordering; no locale-aware or custom genre ranking is required by the issue.
- Do not modify `db/shows.json` or `db/films.json`.

---

*Generated plan for Issue #30.*
