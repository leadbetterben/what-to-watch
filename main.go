package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"what-to-watch/cmd/cli"
	apphttp "what-to-watch/cmd/http"
)

var (
	runCLI          = cli.Run
	startHTTPServer = func(port int) error {
		server := apphttp.NewServer(port)
		return server.Start()
	}
)

func main() {
	// Define command-line flags
	mode := flag.String("mode", "cli", "Run mode: 'cli' for interactive CLI or 'http' for HTTP server")
	port := flag.Int("port", 8080, "HTTP server port (only used in http mode)")
	flag.Parse()

	if exitCode := run(*mode, *port, os.Stderr); exitCode != 0 {
		os.Exit(exitCode)
	}
}

func run(mode string, port int, stderr io.Writer) int {
	switch mode {
	case "cli":
		runCLI()
		return 0
	case "http":
		if err := startHTTPServer(port); err != nil {
			fmt.Fprintf(stderr, "HTTP server error: %v\n", err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "Invalid mode: %s. Use 'cli' or 'http'.\n", mode)
		return 1
	}
}
