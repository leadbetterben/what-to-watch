# Issue #2 — Update the storage when watching a show

## Overview

After listing the shows the user is currently watching, prompt the user to select a show (by the list number) that they have just watched. Update the stored show progress as follows:

- Increment the current episode number.
- If the new episode number exceeds the number of episodes in the current series, increment the series number and reset the episode to 1.
- If the new series number exceeds the number of series in the show, remove `CurrentSeries` and `CurrentEpisode` (the show is finished) and display a congratulations message.
- Persist the changes to `shows.json`.

## Acceptance criteria

- When a user selects a show from the currently-watching list, the show's `CurrentEpisode` and/or `CurrentSeries` are updated according to the rules above.
- When a show is finished (series > number of series), `CurrentSeries` and `CurrentEpisode` are removed (set to `nil`) in the stored JSON.
- The storage (`shows.json`) is updated so the change persists across runs.
- There are unit tests verifying the increment, rollover, and completion behaviors.

## Files to change

- `shows/shows.go` — add a new exported function to update progress (e.g. `MarkEpisodeWatched(index int) error`).
- `main.go` — after showing currently-watching, prompt the user for a selection and call the new `shows` function.
- `db/db.go` — already contains `ReadShows()`; we'll add a `WriteShows([]data.Show) error` function if it doesn't exist.
- `data/data.go` — inspect `Show` struct to ensure fields `CurrentSeries`, `CurrentEpisode` and counts for `NumSeries`/`Series`/`EpisodesPerSeries` (or whichever fields exist) are present.
- `shows/shows_test.go` — add tests for the update logic.

## Implementation plan (detailed steps)

1. Add `db.WriteShows(shows []data.Show) error`.
   - Implement a safe write using temporary file + rename to avoid corrupting `shows.json` on failures.
   - Use same path resolution as `readFile` (use helper `writeFile(path string, data []byte) error`).

2. Implement `shows.MarkEpisodeWatched(idx int) error`.
   - Load shows via `db.ReadShows()`.
   - Determine the mapping between the displayed "currently watching" list and the full shows slice. (Option A: `GetCurrentlyWatching()` returns a filtered list — use the chosen index relative to that list and map back to original shows by matching a unique identifier such as `Title` or an `ID` field. Option B: Present numbered list from the full slice and use the same ordering — prefer Option A to keep the UI consistent.)
   - For the selected show, apply increment logic:
     - If `CurrentEpisode` is nil or `CurrentSeries` is nil, return an error — those shows shouldn't be selectable.
     - Increment `CurrentEpisode`.
     - If `CurrentEpisode` > number of episodes in the current series, set `CurrentEpisode = 1` and increment `CurrentSeries`.
     - If `CurrentSeries` &gt; number of series in the show, set both to `nil` and prepare a congratulations message for the caller.
   - Persist updated shows via `db.WriteShows(shows)`.
   - Return `nil` and optionally a `string` message (or an exported type) with the congratulations text. Simpler approach: return `(string, error)` where the string is non-empty for congratulations.

3. Wire the CLI in `main.go`.
   - After showing currently-watching (via `shows.GetCurrentlyWatching()`), prompt the user for a number (or allow `0`/Enter to cancel).
   - Parse input, call `shows.MarkEpisodeWatched(selectedIndex)` and print any returned message.
   - Re-display updated currently-watching or confirmation.

4. Add unit tests in `shows/shows_test.go`.
   - Test cases:
     - Increment within a series (episode++).
     - Episode rollover -> increment series and reset episode to 1.
     - Series rollover -> clear current series/episode and return congratulations.
   - Use in-memory fixtures: create slice of `data.Show` values, write them to a temp `shows.json` in a temp directory and override path resolution during tests if necessary (or add helper to `db` to allow specifying file path).

5. Manual verification
   - Run `go run .` and exercise the flow: list currently watching, select a show, confirm `shows.json` updated and UI message printed.

## API / function signatures (suggested)

- `db.WriteShows(shows []data.Show) error`
- `func MarkEpisodeWatched(listIndex int) (string, error)` — `listIndex` matches the index in `GetCurrentlyWatching()` (0-based). Returns a message (e.g., congratulations) which may be empty.

## Risks & notes

- Mapping from the filtered currently-watching list back to the full shows slice must be deterministic; prefer using an `ID` field (if `data.Show` contains one) or matching on `Title` as a fallback.
- Tests will be easier if `db` allows writing to a test-specific path or if `ReadShows/WriteShows` accept a path parameter. If modifying `db`'s API is undesirable, tests can run against a temporary working directory and set `os.Executable` behavior — more brittle.

## Next steps

- I will implement `db.WriteShows` and `shows.MarkEpisodeWatched`, then add tests and wire the CLI.


---

*Generated plan for Issue #2.*
