package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"what-to-watch/db"
	"what-to-watch/shows"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Display menu
	fmt.Println("What would you like to view?")
	fmt.Println("1. Currently watching shows")
	fmt.Println("2. Films")
	fmt.Print("Enter your choice (1 or 2): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		viewShows(reader)
	case "2":
		viewFilms()
	default:
		fmt.Println("Invalid input. Please enter 1 or 2.")
	}
}

func viewShows(reader *bufio.Reader) {
	s, err := db.ReadShows()
	if err != nil {
		log.Fatalf("error reading shows: %s", err)
	}

	cw, err := shows.GetCurrentlyWatching(s)
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

	// prompt user to mark a show as watched
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

	// update show watched
	updatedShows, msg, err := shows.MarkEpisodeWatched(s, idx)
	if err != nil {
		fmt.Printf("error updating show: %s\n", err)
		return
	}
	fmt.Println(msg)

	// save updated shows
	if err := db.WriteShows(updatedShows); err != nil {
		fmt.Printf("error saving updated shows: %s\n", err)
		return
	}
}

func viewFilms() {
	films, err := db.ReadFilms()
	if err != nil {
		log.Fatalf("error reading films: %s", err)
	}

	if len(films) == 0 {
		fmt.Println("No films found.")
		return
	}

	// compute column widths
	wIndex := len("Index")
	wName := len("Name")
	wGenre := len("Genre")
	wProvider := len("Provider")

	for _, f := range films {
		if l := len(f.Name); l > wName {
			wName = l
		}
		if l := len(f.Genre); l > wGenre {
			wGenre = l
		}
		if l := len(f.Provider); l > wProvider {
			wProvider = l
		}
	}

	// build format string (left-aligned columns, two spaces between)
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds\n",
		wIndex, wName, wGenre, wProvider)

	// header
	fmt.Printf(format, "Index", "Name", "Genre", "Provider")

	// separator line
	parts := []string{
		strings.Repeat("-", wIndex),
		strings.Repeat("-", wName),
		strings.Repeat("-", wGenre),
		strings.Repeat("-", wProvider),
	}
	fmt.Printf(format, parts[0], parts[1], parts[2], parts[3])

	// rows
	for i, f := range films {
		fmt.Printf(format, strconv.Itoa(i+1), f.Name, f.Genre, f.Provider)
	}
}
