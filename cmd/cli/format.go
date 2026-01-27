package cli

import (
	"fmt"
	"strconv"
	"strings"
	"what-to-watch/data"
)

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
