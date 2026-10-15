# Plan: Harden the HTTP server defaults and lifecycle

Issue: [#26](https://github.com/leadbetterben/what-to-watch/issues/26)

## Goal

Keep HTTP mode safe and useful for local-only use: bind to loopback by default, apply bounded connection timeouts, and keep server setup configurable and testable. Remote deployment, authentication, and access by web or Android clients are explicitly deferred until that deployment is in scope.

## Steps

1. **Add safe server configuration**
   - Replace the all-interface listener with a loopback listener (`127.0.0.1`) plus the selected port.
   - Preserve the existing `-port` behavior and default. Avoid adding remote bind configuration in this local-only change.
   - Create an `http.Server` with documented `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` values, and start it using the configured address.
   - Keep construction/startup seams injectable for tests; avoid changing the existing route-handler interface unnecessarily.

2. **Keep future deployment work separable**
   - Do not add authentication or remote bind options for this local-only change.
   - Keep server construction/configuration sufficiently encapsulated that a future deployment can add an authentication strategy for `POST /shows/watch` alongside its client and hosting requirements.

3. **Add focused tests**
   - Verify default startup is configured for loopback.
   - Verify all four server timeout fields are populated with the intended values.
   - Keep existing handler tests passing and use an injectable listener/server seam rather than opening externally reachable sockets.

4. **Document operation and security**
   - Update the HTTP section of `README.md` with the default loopback behavior and timeout behavior.
   - State that HTTP mode is intended for local use and that remote deployment/access is not currently supported; avoid documenting an insecure remote workaround.

5. **Validate**
   - Run `go vet ./...`, `go test ./...`, `go test ./... -cover`, and `go build ./...`.
   - Review `git diff --check` and confirm runtime JSON files remain untouched.

## Acceptance criteria

- Default HTTP startup listens only on loopback and cannot be reached by other machines on the network.
- The configured `http.Server` sets read-header, read, write, and idle timeouts.
- Server construction and loopback configuration have focused automated coverage.
- README explains the local-only scope and does not suggest exposing the unauthenticated write route remotely.
- Remote deployment, authentication, and web/Android client access remain future work.

## Likely files

- `main.go` — startup configuration, if needed to pass through the loopback address.
- `cmd/http/http.go` — loopback server configuration, timeout values, and lifecycle.
- `cmd/http/http_test.go` — loopback configuration and timeout coverage.
- `README.md` — HTTP bind and remote-access guidance.
