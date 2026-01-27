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

// formatShowsTable formats shows into a table string
func formatShowsTable(s []data.Show) string {
	if len(s) == 0 {
		return "No shows currently being watched.\n"
	}

	// compute column widths
	wIndex := len("Index")
	wName := len("Name")
	wGenre := len("Genre")
	wProvider := len("Provider")
	wSeries := len("Series")
	wEpisode := len("Episode")

	for _, r := range s {
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

	var buf strings.Builder

	// header
	buf.WriteString(fmt.Sprintf(format, "Index", "Name", "Genre", "Provider", "Series", "Episode"))

	// separator line
	parts := []string{
		strings.Repeat("-", wIndex),
		strings.Repeat("-", wName),
		strings.Repeat("-", wGenre),
		strings.Repeat("-", wProvider),
		strings.Repeat("-", wSeries),
		strings.Repeat("-", wEpisode),
	}
	buf.WriteString(fmt.Sprintf(format, parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]))

	// rows
	for i, r := range s {
		buf.WriteString(fmt.Sprintf(format, strconv.Itoa(i+1), r.Name, r.Genre, r.Provider, r.Series, r.Episode))
	}

	return buf.String()
}

// formatShowsByGenreTable formats unwatched shows into a table string
func formatShowsByGenreTable(s []data.Show) string {
	if len(s) == 0 {
		return "No unwatched shows in this genre.\n"
	}

	// compute column widths
	wIndex := len("Index")
	wName := len("Name")
	wProvider := len("Provider")

	for _, r := range s {
		if l := len(r.Name); l > wName {
			wName = l
		}
		if l := len(r.Provider); l > wProvider {
			wProvider = l
		}
	}

	// build format string (left-aligned columns, two spaces between)
	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds\n",
		wIndex, wName, wProvider)

	var buf strings.Builder

	// header
	buf.WriteString(fmt.Sprintf(format, "Index", "Name", "Provider"))

	// separator line
	parts := []string{
		strings.Repeat("-", wIndex),
		strings.Repeat("-", wName),
		strings.Repeat("-", wProvider),
	}
	buf.WriteString(fmt.Sprintf(format, parts[0], parts[1], parts[2]))

	// rows
	for i, r := range s {
		buf.WriteString(fmt.Sprintf(format, strconv.Itoa(i+1), r.Name, r.Provider))
	}

	return buf.String()
}

// formatFilmsTable formats films into a table string
func formatFilmsTable(films []data.Film) string {
	if len(films) == 0 {
		return "No films found.\n"
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

	var buf strings.Builder

	// header
	buf.WriteString(fmt.Sprintf(format, "Index", "Name", "Genre", "Provider"))

	// separator line
	parts := []string{
		strings.Repeat("-", wIndex),
		strings.Repeat("-", wName),
		strings.Repeat("-", wGenre),
		strings.Repeat("-", wProvider),
	}
	buf.WriteString(fmt.Sprintf(format, parts[0], parts[1], parts[2], parts[3]))

	// rows
	for i, f := range films {
		buf.WriteString(fmt.Sprintf(format, strconv.Itoa(i+1), f.Name, f.Genre, f.Provider))
	}

	return buf.String()
}
