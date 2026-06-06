package cli

import (
	"testing"

	"what-to-watch/data"
)

func TestFormatShowsTable(t *testing.T) {
	shows := []data.Show{
		{Name: "A", Genre: "Drama", Provider: "BBC", Series: "1", Episode: "2"},
		{Name: "Long Show", Genre: "Comedy", Provider: "Netflix", Series: "12", Episode: "10"},
	}

	expected := "" +
		"Index  Name       Genre   Provider  Series  Episode\n" +
		"-----  ---------  ------  --------  ------  -------\n" +
		"1      A          Drama   BBC       1       2      \n" +
		"2      Long Show  Comedy  Netflix   12      10     \n"

	if got := formatShowsTable(shows); got != expected {
		t.Fatalf("expected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestFormatShowsTableEmpty(t *testing.T) {
	expected := "No shows currently being watched.\n"

	if got := formatShowsTable(nil); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestFormatShowsByGenreTable(t *testing.T) {
	shows := []data.Show{
		{Name: "Long Drama", Provider: "Netflix"},
		{Name: "A", Provider: "BBC"},
	}

	expected := "" +
		"Index  Name        Provider\n" +
		"-----  ----------  --------\n" +
		"1      Long Drama  Netflix \n" +
		"2      A           BBC     \n"

	if got := formatShowsByGenreTable(shows); got != expected {
		t.Fatalf("expected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestFormatShowsByGenreTableEmpty(t *testing.T) {
	expected := "No unwatched shows in this genre.\n"

	if got := formatShowsByGenreTable(nil); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestFormatFilmsTable(t *testing.T) {
	films := []data.Film{
		{Name: "Arrival", Genre: "Sci-Fi", Provider: "Netflix"},
		{Name: "Up", Genre: "Family", Provider: "Disney+"},
	}

	expected := "" +
		"Index  Name     Genre   Provider\n" +
		"-----  -------  ------  --------\n" +
		"1      Arrival  Sci-Fi  Netflix \n" +
		"2      Up       Family  Disney+ \n"

	if got := formatFilmsTable(films); got != expected {
		t.Fatalf("expected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestFormatFilmsTableEmpty(t *testing.T) {
	expected := "No films found.\n"

	if got := formatFilmsTable(nil); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
