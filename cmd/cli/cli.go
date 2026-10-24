package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"what-to-watch/data"
	"what-to-watch/handlers"
)

var (
	getCurrentlyWatchingShows = handlers.GetCurrentlyWatchingShows
	markShowWatched           = handlers.MarkShowWatched
	getAllFilms               = handlers.GetAllFilms
	getAvailableGenres        = handlers.GetAvailableGenres
	getUnwatchedShowsByGenre  = handlers.GetUnwatchedShowsByGenre
	getHalloweenEpisodes      = func() ([]data.ShowSpecialEpisode, error) {
		return handlers.GetSpecialEpisodes(data.SpecialCategoryHalloween)
	}
	getChristmasEpisodes = func() ([]data.ShowSpecialEpisode, error) {
		return handlers.GetSpecialEpisodes(data.SpecialCategoryChristmas)
	}
	getRandomHalloweenEpisode = func() (data.ShowSpecialEpisode, error) {
		return handlers.GetRandomSpecialEpisode(data.SpecialCategoryHalloween)
	}
	getRandomChristmasEpisode = func() (data.ShowSpecialEpisode, error) {
		return handlers.GetRandomSpecialEpisode(data.SpecialCategoryChristmas)
	}
)

// Run starts the interactive CLI mode
func Run() {
	reader := bufio.NewReader(os.Stdin)

	// Display menu
	fmt.Println("What would you like to view?")
	fmt.Println("1. Currently watching shows")
	fmt.Println("2. Films")
	fmt.Println("3. Shows by genre")
	fmt.Println("4. Halloween episodes")
	fmt.Println("5. Christmas episodes")
	fmt.Println("6. Random Halloween episode")
	fmt.Println("7. Random Christmas episode")
	fmt.Print("Enter your choice (1, 2, 3, 4, 5, 6, or 7): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		viewShows(reader)
	case "2":
		viewFilms()
	case "3":
		viewShowsByGenre(reader)
	case "4":
		viewSeasonalSpecials("Halloween", getHalloweenEpisodes)
	case "5":
		viewSeasonalSpecials("Christmas", getChristmasEpisodes)
	case "6":
		viewRandomSeasonalSpecial("Halloween", getRandomHalloweenEpisode)
	case "7":
		viewRandomSeasonalSpecial("Christmas", getRandomChristmasEpisode)
	default:
		fmt.Println("Invalid input. Please enter a number from 1 to 7.")
	}
}

func viewShows(reader *bufio.Reader) {
	shows, err := getCurrentlyWatchingShows()
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

	isCompleted, err := markShowWatched(idx)
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
	films, err := getAllFilms()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println(formatFilmsTable(films))
}

func viewShowsByGenre(reader *bufio.Reader) {
	// Get available genres
	genres, err := getAvailableGenres()
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
	shows, err := getUnwatchedShowsByGenre(selectedGenre)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Unwatched shows in genre '%s':\n", selectedGenre)
	fmt.Println(formatShowsByGenreTable(shows))
}

func viewSeasonalSpecials(category string, getEpisodes func() ([]data.ShowSpecialEpisode, error)) {
	episodes, err := getEpisodes()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("%s episodes:\n", category)
	fmt.Println(formatSpecialEpisodesTable(category, episodes))
}

func viewRandomSeasonalSpecial(category string, getEpisode func() (data.ShowSpecialEpisode, error)) {
	episode, err := getEpisode()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Random %s episode: %s (%s), Series %d Episode %d\n", strings.ToLower(category), episode.ShowName, episode.Provider, episode.Series, episode.Episode)
}
