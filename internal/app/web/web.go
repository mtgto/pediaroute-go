package web

import (
	"errors"
	"math/rand"
	"strings"
)

// Data structure
type Wikipedia struct {
	titles                  []string
	redirectSet             map[int]struct{} // whether page index is redirected
	lowercaseTitleToIndices map[string][]int // map from lowercase of title to index
	forwardLinks            [][]int
	backwardLinks           [][]int
}

var NotFound = errors.New("Not found")

const MaxDepth = 6

func (w *Wikipedia) Search(from, to string) []string {
	start := w.generateWordSet(from)
	goal := w.generateWordSet(to)
	way, err := w.search(&start, &goal, 0)
	if err != nil {
		return nil
	} else {
		result := make([]string, len(way))
		for i, index := range way {
			result[i] = w.titles[index]
		}
		return result
	}
}

func (w *Wikipedia) IsWordExists(word string) bool {
	_, exists := w.lowercaseTitleToIndices[strings.ToLower(word)]
	return exists
}

// return whether word is only redirect page title
func (w *Wikipedia) IsWordRedirect(word string) bool {
	if indices, exists := w.lowercaseTitleToIndices[strings.ToLower(word)]; exists {
		for _, index := range indices {
			if _, exists := w.redirectSet[index]; !exists {
				return false
			}
		}
		return true
	}
	return false
}

func (w *Wikipedia) GetRandomPage() string {
	for {
		title := w.titles[rand.Intn(len(w.titles))]
		if !w.IsWordRedirect(title) {
			return title
		}
	}
}

func (w *Wikipedia) generateWordSet(word string) map[int]struct{} {
	set := make(map[int]struct{}, 1)
	if indices, ok := w.lowercaseTitleToIndices[strings.ToLower(word)]; ok {
		for _, index := range indices {
			set[index] = struct{}{}
		}
	}
	return set
}

func (w *Wikipedia) search(start, goal *map[int]struct{}, depth int) ([]int, error) {
	if depth >= MaxDepth {
		return nil, NotFound
	}
	if len(*start) < len(*goal) {
		nextStarts := make(map[int]struct{}, 0)
		for from, _ := range *start {
			for _, to := range w.forwardLinks[from] {
				if _, ok := (*goal)[to]; ok {
					return []int{from, to}, nil
				} else {
					nextStarts[to] = struct{}{}
				}
			}
		}
		way, err := w.search(&nextStarts, goal, depth+1)
		if err != nil {
			return nil, err
		}
		for _, index := range w.backwardLinks[way[0]] {
			if _, ok := (*start)[index]; ok {
				return append([]int{index}, way...), nil
			}
		}
	} else {
		nextGoals := make(map[int]struct{}, 0)
		for to, _ := range *goal {
			for _, from := range w.backwardLinks[to] {
				if _, ok := (*start)[from]; ok {
					return []int{from, to}, nil
				} else {
					nextGoals[from] = struct{}{}
				}
			}
		}
		way, err := w.search(start, &nextGoals, depth+1)
		if err != nil {
			return nil, err
		}
		for _, index := range w.forwardLinks[way[len(way)-1]] {
			if _, ok := (*goal)[index]; ok {
				return append(way, index), nil
			}
		}
	}
	return nil, NotFound
}
