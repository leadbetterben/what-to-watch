package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	i := 1
	for _, s := range shows {
		if s.CurrentSeries != nil || s.CurrentEpisode != nil {
			cs := 0
			ce := 0
			if s.CurrentSeries != nil {
				cs = *s.CurrentSeries
			}
			if s.CurrentEpisode != nil {
				ce = *s.CurrentEpisode
			}

			fmt.Printf("%d. %s: %s on %s - series %d episode %d\n", i, s.Name, s.Genre, s.Provider, cs, ce)

			i++
		}
	}
}
