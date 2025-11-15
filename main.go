package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"what-to-watch/shows"
)

func main() {
	cw, err := shows.GetCurrentlyWatching()
	if err != nil {
		log.Fatalf("error getting currently watching shows: %s", err)
	}

	if len(cw) == 0 {
		fmt.Println("You are not currently watching any shows.")
		return
	}

	// compute column widths
	wIndex := len("Index")
	wName := len("Name")
	wGenre := len("Genre")
	wProvider := len("Provider")
	wSeries := len("Series")
	wEpisode := len("Episode")

	for _, r := range cw {
		if l := len(r.Name); l > wName {
			wName = l
		}
		if l := len(r.Genre); l > wGenre {
			wGenre = l
		}
		if l := len(r.Provider); l > wProvider {
			wProvider = l
		}
		if l := len(r.Series); l > wSeries {
			wSeries = l
		}
		if l := len(r.Episode); l > wEpisode {
			wEpisode = l
		}
	}

	// build format string (left-aligned columns, two spaces between)
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds\n",
		wIndex, wName, wGenre, wProvider, wSeries, wEpisode)

	// header
	fmt.Printf(format, "Index", "Name", "Genre", "Provider", "Series", "Episode")

	// separator line
	parts := []string{
		strings.Repeat("-", wIndex),
		strings.Repeat("-", wName),
		strings.Repeat("-", wGenre),
		strings.Repeat("-", wProvider),
		strings.Repeat("-", wSeries),
		strings.Repeat("-", wEpisode),
	}
	fmt.Printf(format, parts[0], parts[1], parts[2], parts[3], parts[4], parts[5])

	// rows
	i := 1
	for _, r := range cw {
		fmt.Printf(format, strconv.Itoa(i), r.Name, r.Genre, r.Provider, r.Series, r.Episode)
		i++
	}

	// Prompt user to mark a show as watched
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the Index of the show you watched (0 to cancel): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" || input == "0" {
		fmt.Println("No changes made.")
		return
	}
	idx, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("invalid input: %s\n", input)
		return
	}

	msg, err := shows.MarkEpisodeWatched(idx)
	if err != nil {
		fmt.Printf("error updating show: %s\n", err)
		return
	}
	if msg != "" {
		fmt.Println(msg)
	} else {
		fmt.Println("Updated show progress.")
	}
}
