package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name            string
		mode            string
		httpResponse    error
		expectExitCode  int
		expectRunCalled bool
		expectPort      int
		expectErr       string
	}{
		{
			name:            "CLI mode",
			mode:            "cli",
			expectExitCode:  0,
			expectRunCalled: true,
			expectErr:       "",
		},
		{
			name:            "HTTP mode",
			mode:            "http",
			expectExitCode:  0,
			expectRunCalled: false,
			expectPort:      9090,
			expectErr:       "",
		},
		{
			name:            "HTTP mode with error",
			mode:            "http",
			httpResponse:    errors.New("listen failed"),
			expectExitCode:  1,
			expectRunCalled: false,
			expectPort:      9090,
			expectErr:       "HTTP server error: listen failed\n",
		},
		{
			name:            "Invalid mode",
			mode:            "invalid",
			expectExitCode:  1,
			expectRunCalled: false,
			expectErr:       "Invalid mode: invalid. Use 'cli' or 'http'.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalRunCLI := runCLI
			originalStartHTTPServer := startHTTPServer
			defer func() {
				runCLI = originalRunCLI
				startHTTPServer = originalStartHTTPServer
			}()

			runCalled := false
			runCLI = func() {
				runCalled = true
			}
			var gotPort int
			startHTTPServer = func(port int) error {
				gotPort = port
				return tt.httpResponse
			}

			var stderr bytes.Buffer
			exitCode := run(tt.mode, 9090, &stderr)

			if exitCode != tt.expectExitCode {
				t.Fatalf("expected exit code %d, got %d", tt.expectExitCode, exitCode)
			}
			if runCalled != tt.expectRunCalled {
				t.Fatalf("expected runCLI called: %v, got %v", tt.expectRunCalled, runCalled)
			}
			if gotPort != tt.expectPort {
				t.Fatalf("expected HTTP server port %d, got %d", tt.expectPort, gotPort)
			}
			if stderr.String() != tt.expectErr {
				t.Fatalf("expected stderr to be %q, got %q", tt.expectErr, stderr.String())
			}
		})
	}
}
