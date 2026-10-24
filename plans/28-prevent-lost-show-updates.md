# Issue #28: Prevent lost updates when multiple requests write show data

## Overview

Protect the complete show progress update from lost updates. Today
`handlers.MarkShowWatched` reads current shows, applies the episode update, then
calls `db.WriteCurrentShows`, which independently rereads and merges the file.
Atomic rename prevents a partially written JSON file, but does not make this
read–modify–write sequence atomic.

For this iteration, guarantee serialization within one running application
process. Multiple application processes that write the same `shows.json` are
unsupported; document this deployment limit. A future cross-process lock can
extend the guarantee if simultaneous instances become a requirement.

## Implementation steps

1. **Move the update transaction into `db`.** Add an update function that
   acquires a package-level mutex, reads the latest complete shows array,
   selects current shows, invokes the episode update operation, merges the
   updated current records with all non-current records, and writes the full
   replacement before releasing the mutex. Keep the existing same-directory
   temporary file and rename behavior for safe replacement.
2. **Route watched-episode writes through the transaction.** Refactor
   `handlers.MarkShowWatched` so it no longer reads data before entering the
   protected operation. Preserve its index semantics and returned completion
   status/error behavior. Keep pure show progression logic in `shows`.
3. **Define behavior for other writers.** Ensure every in-process operation
   that can replace current show data uses the same lock and transaction
   boundary. Avoid leaving a public write path that can interleave with the
   transaction and reintroduce lost updates.
4. **Document deployment scope.** Update `README.md` to state that concurrent
   requests within one process are serialized and that multiple processes
   sharing the same data file are not supported. Keep the documented atomic
   replacement behavior.
5. **Add focused concurrency coverage.** Use an isolated temporary shows file
   and launch simultaneous updates to different shows/episodes through the
   supported update path. Assert that both successful updates are present,
   non-current records remain intact, and the final file parses as valid JSON.

## Acceptance criteria

- Concurrent in-process update attempts do not silently discard a successful
  update.
- The lock covers reading the latest data, applying the update, merging all
  records, and replacing the file.
- Existing non-current records are preserved.
- File replacement remains safe against partial-file corruption.
- A focused concurrency test exercises simultaneous update attempts.
- README clearly states the single-process guarantee and multi-process limit.

## Files expected to change

- `db/db.go` and `db/db_test.go` — serialized update transaction and storage
  concurrency coverage.
- `handlers/handlers.go` and relevant handler tests — use the transaction API.
- `README.md` — deployment/concurrency scope.

## Validation

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `git diff --check`

## Scope note

If multiple processes must be supported, replace the process-local mutex with
a cross-process lock or a conflict-detection/retry approach before claiming
that deployment guarantee. Any lock must coordinate on the resolved data-file
path so all writers to the same file share the same lock.
