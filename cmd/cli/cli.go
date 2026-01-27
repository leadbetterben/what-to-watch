package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"what-to-watch/handlers"
)

// Run starts the interactive CLI mode
func Run() {
	reader := bufio.NewReader(os.Stdin)

	// Display menu
	fmt.Println("What would you like to view?")
	fmt.Println("1. Currently watching shows")
	fmt.Println("2. Films")
	fmt.Println("3. Shows by genre")
	fmt.Print("Enter your choice (1, 2, or 3): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		viewShows(reader)
	case "2":
		viewFilms()
	case "3":
		viewShowsByGenre(reader)
	default:
		fmt.Println("Invalid input. Please enter 1, 2, or 3.")
	}
}

func viewShows(reader *bufio.Reader) {
	shows, err := handlers.GetCurrentlyWatchingShows()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(formatShowsTable(shows))

	// prompt user to mark a show as watched
	fmt.Print("Enter the Index of the show you watched (0 to cancel): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" || input == "0" {
		fmt.Println("No changes made.")
		return
	}

	idx, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid input: %s\n", input)
		return
	}

	isCompleted, err := handlers.MarkShowWatched(idx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if isCompleted {
		fmt.Printf("Show %d marked as watched and completed!\n", idx)
	} else {
		fmt.Printf("Show %d marked as watched.\n", idx)
	}
}

func viewFilms() {
	films, err := handlers.GetAllFilms()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(formatFilmsTable(films))
}

func viewShowsByGenre(reader *bufio.Reader) {
	// Get available genres
	genres, err := handlers.GetAvailableGenres()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if len(genres) == 0 {
		fmt.Println("No genres available.")
		return
	}

	// Display genres
	fmt.Println("Available genres:")
	for i, genre := range genres {
		fmt.Printf("%d. %s\n", i+1, genre)
	}

	// Get user selection
	fmt.Print("Enter the genre number (0 to cancel): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" || input == "0" {
		fmt.Println("No selection made.")
		return
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(genres) {
		fmt.Printf("Invalid input: %s\n", input)
		return
	}

	selectedGenre := genres[idx-1]

	// Get shows for selected genre
	shows, err := handlers.GetUnwatchedShowsByGenre(selectedGenre)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Unwatched shows in genre '%s':\n", selectedGenre)
	fmt.Println(formatShowsByGenreTable(shows))
}
