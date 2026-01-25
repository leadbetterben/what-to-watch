package main

import (
	"flag"
	"fmt"
	"os"

	"what-to-watch/cmd/cli"
	"what-to-watch/cmd/http"
)

func main() {
	// Define command-line flags
	mode := flag.String("mode", "cli", "Run mode: 'cli' for interactive CLI or 'http' for HTTP server")
	port := flag.Int("port", 8080, "HTTP server port (only used in http mode)")
	flag.Parse()

	switch *mode {
	case "cli":
		cli.Run()
	case "http":
		server := http.NewServer(*port)
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'cli' or 'http'.\n", *mode)
		os.Exit(1)
	}
}
