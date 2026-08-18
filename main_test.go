package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCLIMode(t *testing.T) {
	restore := stubMainFunctions()
	t.Cleanup(restore)

	called := false
	runCLI = func() {
		called = true
	}

	var stderr bytes.Buffer
	exitCode := run("cli", 8080, &stderr)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if !called {
		t.Fatalf("expected CLI runner to be called")
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunHTTPMode(t *testing.T) {
	restore := stubMainFunctions()
	t.Cleanup(restore)

	var gotPort int
	startHTTPServer = func(port int) error {
		gotPort = port
		return nil
	}

	var stderr bytes.Buffer
	exitCode := run("http", 9090, &stderr)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if gotPort != 9090 {
		t.Fatalf("expected HTTP server port 9090, got %d", gotPort)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunHTTPModeReturnsError(t *testing.T) {
	restore := stubMainFunctions()
	t.Cleanup(restore)

	startHTTPServer = func(int) error {
		return errors.New("listen failed")
	}

	var stderr bytes.Buffer
	exitCode := run("http", 8080, &stderr)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "HTTP server error: listen failed") {
		t.Fatalf("expected HTTP error message, got %q", stderr.String())
	}
}

func TestRunInvalidMode(t *testing.T) {
	restore := stubMainFunctions()
	t.Cleanup(restore)

	var stderr bytes.Buffer
	exitCode := run("bogus", 8080, &stderr)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "Invalid mode: bogus. Use 'cli' or 'http'.") {
		t.Fatalf("expected invalid mode message, got %q", stderr.String())
	}
}

func stubMainFunctions() func() {
	originalRunCLI := runCLI
	originalStartHTTPServer := startHTTPServer

	runCLI = func() {}
	startHTTPServer = func(int) error { return nil }

	return func() {
		runCLI = originalRunCLI
		startHTTPServer = originalStartHTTPServer
	}
}
