package web

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	titles := []string{"A", "B", "C", "D", "E"}
	lowercaseTitleToIndices := make(map[string][]int, 0)
	for index, title := range titles {
		lowercaseTitleToIndices[strings.ToLower(title)] = []int{index}
	}
	forwardLinks := [][]int{
		[]int{1},    // A -> B
		[]int{2, 4}, // B -> C, E
		[]int{4},    // C -> E
		[]int{0},    // D -> A
		[]int{3},    // E -> D
	}
	backwardLinks := [][]int{
		[]int{3},
		[]int{0},
		[]int{1},
		[]int{4},
		[]int{1, 2},
	}
	w := Wikipedia{titles: titles, lowercaseTitleToIndices: lowercaseTitleToIndices, forwardLinks: forwardLinks, backwardLinks: backwardLinks}

	testSuccess := func(start, goal string, expected []string) {
		result, err := w.Search(start, goal)
		if err == NotFound {
			t.Fatalf("Not found the way from %s to %s\n", start, goal)
		} else if len(result) > len(expected) {
			t.Fatalf("Wrong way found from %s to %s\n Expected: %v, Actual: %v\n", start, goal, expected, result)
		} else {
			t.Logf("Found the way from %s to %s: %v\n", start, goal, result)
		}
	}

	testSuccess("A", "B", []string{"A", "B"})
	testSuccess("A", "C", []string{"A", "B", "C"})
	testSuccess("A", "E", []string{"A", "B", "E"})
	testSuccess("A", "D", []string{"A", "B", "E", "D"})
}
