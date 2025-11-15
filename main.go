package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Show struct {
	Name           string `json:"name"`
	Genre          string `json:"genre"`
	Episodes       []int  `json:"episodes"`
	Provider       string `json:"provider"`
	CurrentSeries  *int   `json:"currentSeries,omitempty"`
	CurrentEpisode *int   `json:"currentEpisode,omitempty"`
}

func main() {
	data, err := os.ReadFile("shows.json")
	if err != nil {
		log.Fatalf("reading shows.json: %v", err)
	}

	var shows []Show
	if err := json.Unmarshal(data, &shows); err != nil {
		log.Fatalf("parsing shows.json: %v", err)
	}

	// collect rows for shows that have currentSeries or currentEpisode
	type row struct {
		Name     string
		Genre    string
		Provider string
		Series   string
		Episode  string
	}

	var rows []row
	for _, s := range shows {
		if s.CurrentSeries != nil || s.CurrentEpisode != nil {
			series := "-"
			episode := "-"
			if s.CurrentSeries != nil {
				series = strconv.Itoa(*s.CurrentSeries)
			}
			if s.CurrentEpisode != nil {
				episode = strconv.Itoa(*s.CurrentEpisode)
			}
			rows = append(rows, row{
				Name:     s.Name,
				Genre:    s.Genre,
				Provider: s.Provider,
				Series:   series,
				Episode:  episode,
			})
		}
	}

	if len(rows) == 0 {
		return
	}

	// compute column widths
	wIndex := len("Index")
	wName := len("Name")
	wGenre := len("Genre")
	wProvider := len("Provider")
	wSeries := len("Series")
	wEpisode := len("Episode")

	for _, r := range rows {
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
	for _, r := range rows {
		fmt.Printf(format, strconv.Itoa(i), r.Name, r.Genre, r.Provider, r.Series, r.Episode)
		i++
	}
}
