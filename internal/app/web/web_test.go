package web

import (
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	// A -> B
	// B -> C, E
	// C -> E
	// D -> A
	// E -> D
	// F -> A
	pages := []core.Page{
		{Id: 0, TitleOffset: 0, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 0, ForwardLinkLength: 1, BackwardLinkIndex: 7, BackwardLinkLength: 1},
		{Id: 1, TitleOffset: 1, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 1, ForwardLinkLength: 2, BackwardLinkIndex: 8, BackwardLinkLength: 1},
		{Id: 2, TitleOffset: 2, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 3, ForwardLinkLength: 1, BackwardLinkIndex: 9, BackwardLinkLength: 1},
		{Id: 3, TitleOffset: 3, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 4, ForwardLinkLength: 1, BackwardLinkIndex: 10, BackwardLinkLength: 1},
		{Id: 4, TitleOffset: 4, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 5, ForwardLinkLength: 1, BackwardLinkIndex: 11, BackwardLinkLength: 2},
		{Id: 5, TitleOffset: 5, TitleLength: 1, IsRedirect: false, ForwardLinkIndex: 6, ForwardLinkLength: 1, BackwardLinkIndex: 13, BackwardLinkLength: 0},
	}
	titleFile, err := os.Open("testdata/title.dat")
	if err != nil {
		t.Fatalf("Failed to open title data file: %v\n", err)
	}
	linkFile, err := os.Open("testdata/link.dat")
	if err != nil {
		t.Fatalf("Failed to open link data: %v\n", err)
	}
	w := Wikipedia{
		pages:     pages,
		titleFile: titleFile,
		linkFile:  linkFile,
	}

	testSearch := func(start, goal string, expected []string) {
		result := w.Search(start, goal)
		if len(result) > len(expected) {
			t.Fatalf("Wrong way found from %s to %s\n Expected: %v, Actual: %v\n", start, goal, expected, result)
		} else if len(result) < len(expected) {
			t.Fatalf("Impossible way found from %s to %s\n Expected: %v, Actual: %v\n", start, goal, expected, result)
		} else {
			t.Logf("Found the way from %s to %s: %v\n", start, goal, result)
		}
	}

	testSearch("A", "B", []string{"A", "B"})
	testSearch("A", "C", []string{"A", "B", "C"})
	testSearch("A", "E", []string{"A", "B", "E"})
	testSearch("A", "D", []string{"A", "B", "E", "D"})
	testSearch("A", "F", []string{})
	testSearch("AA", "B", []string{})
}
