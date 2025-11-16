# AI Agent Instructions for this Repository

Purpose
- This document gives concise, practical instructions for AI agents (and human reviewers) who will make, test, and validate changes in the `what-to-watch` Go CLI repository.

Environment & Tooling
- Go toolchain: `go 1.25.4` (use the version in `go.mod` and CI).
- Shell used for commands in examples: Windows PowerShell (v5.1).
- Formatting: run `gofmt -w .` before committing. Optionally run `go vet ./...`.

Key files and intent
- `main.go` — program entrypoint; reads shows, prompts user, updates JSON.
- `db/db.go` — read/write helpers and `getFullPath` logic for `db/shows.json`.
- `data/data.go` — `Show` struct and model types used across packages.
- `shows/shows.go` — business logic: `GetCurrentlyWatching` and `MarkEpisodeWatched`.
- `shows/shows_test.go` — unit tests for `shows` package (fast, in-memory).
- `db/shows.json` — canonical on-disk data used when running `go run .`.

Common Commands (PowerShell)
```
go version
gofmt -w .
go build ./...
go test ./...
go run .
```

Agent Workflow Expectations
- Use the repository's Go version (`go 1.25.4`) to avoid CI drift.
- Run `go build ./...` and `go test ./...` locally before proposing changes.
- Keep changes minimal and focused; do not refactor unrelated code.

Editing & Patching
- Use the provided `apply_patch` workflow (or equivalent) to modify files.
- When editing files:
	- Preserve existing coding style and indentation.
	- Avoid unrelated formatting changes.
	- If a change touches file I/O, verify `db.getFullPath` semantics (built binary vs `go run`).

Testing & Validation
- Unit tests live under `shows/`. They are fast and authoritative for business logic.
- Tests should not rely on `db/shows.json` being modified — tests use in-memory data.
- After implementing changes, run:
```
gofmt -w .
go vet ./...
go test ./...
```

Committing & PRs
- Commit messages should be concise and descriptive.
- Before opening a PR, ensure `go build ./...` and `go test ./...` pass.
- CI runs `go build -v ./...` and `go test -v ./...` on Go 1.25.4; target the same locally.

Documentation
- When making changes, update documentation files where relevant:
	- `README.md` — end-user usage, examples, and high-level project description.
	- `copilot-instructions.md` — any repository-specific agent guidance or workflow updates.
	- `AGENTS.md` — keep this file in sync with any changes to agent expectations, tooling, or workflows.
- This is especially important for new features, public API changes, large refactors, or structural changes. Include brief notes in the PR description summarizing doc updates.

Data and File I/O Notes
- `db.ReadShows()` searches for `shows.json` next to the executable and falls back to the source `db` directory. Built binaries may expect `shows.json` next to the binary, whereas `go run .` finds the source `db/shows.json`.
- `db.WriteShows()` writes atomically (temp file then rename). Be careful not to accidentally commit runtime-modified JSON files.

Agent Behaviour & Tooling Conventions
- Use the repo-provided `manage_todo_list` for planning. Always:
	1) Create a short todo list for multi-step tasks.
	2) Mark exactly one todo `in-progress` while working on it.
 3) Mark todos `completed` as you finish them.
- Provide brief preambles before automated tool calls (1–2 short sentences) describing intent.
- After groups of changes or several tool calls, return a concise progress update (1–2 sentences).
-
Plans Directory
- All AI-generated plans must be saved in the `plans/` directory as Markdown files (`*.md`). Keep filenames descriptive (for example `ISSUE-2-plan.md`) and include clear steps, acceptance criteria, and any commands or files changed.

Safety and Scope
- Do not introduce untrusted external dependencies unless necessary and approved.
- Avoid making broad stylistic or structural changes without explicit instructions.

If You Need More
- If any CI or test failures are unclear, run the failing commands locally and include the exact error output in your report.
- When in doubt about changing `db/shows.json` on-disk, ask for guidance or open a draft PR.

Contact / Handoff
- When work is complete, include in the final message:
	- Files changed.
	- Commands you ran to validate (`build`, `test`).
	- Any remaining manual checks you recommend.

-- End of agent instructions --

