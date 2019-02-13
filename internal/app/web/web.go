package web

import (
	"encoding/binary"
	"errors"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Data structure
type Wikipedia struct {
	pages                   []core.Page
	lowercaseTitleToIndices map[string][]int32 // map from lowercase of title to index
	linkFile                *os.File
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
			result[i] = w.pages[index].Title
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
			if !w.pages[index].IsRedirect {
				return false
			}
		}
		return true
	}
	return false
}

func (w *Wikipedia) GetRandomPage() string {
	for {
		title := w.pages[rand.Intn(len(w.pages))].Title
		if !w.IsWordRedirect(title) {
			return title
		}
	}
}

func (w *Wikipedia) generateWordSet(word string) map[int32]struct{} {
	set := make(map[int32]struct{}, 1)
	if indices, ok := w.lowercaseTitleToIndices[strings.ToLower(word)]; ok {
		for _, index := range indices {
			set[index] = struct{}{}
		}
	}
	return set
}

func (w *Wikipedia) search(start, goal *map[int32]struct{}, depth int) ([]int32, error) {
	if depth >= MaxDepth {
		return nil, NotFound
	}
	if len(*start) < len(*goal) {
		nextStarts := make(map[int32]struct{}, 0)
		for from, _ := range *start {
			links, err := w.forwardLinks(w.pages[from])
			if err != nil {
				return nil, err
			}
			for _, to := range links {
				toIndex := to
				if _, ok := (*goal)[toIndex]; ok {
					return []int32{from, toIndex}, nil
				} else {
					nextStarts[toIndex] = struct{}{}
				}
			}
		}
		way, err := w.search(&nextStarts, goal, depth+1)
		if err != nil {
			return nil, err
		}
		links, err := w.backwardLinks(w.pages[way[0]])
		if err != nil {
			return nil, err
		}
		for _, from := range links {
			fromIndex := from
			if _, ok := (*start)[fromIndex]; ok {
				return append([]int32{fromIndex}, way...), nil
			}
		}
	} else {
		nextGoals := make(map[int32]struct{}, 0)
		for to, _ := range *goal {
			links, err := w.backwardLinks(w.pages[to])
			if err != nil {
				log.Printf("err: %v\n", err)
				return nil, err
			}
			for _, from := range links {
				fromIndex := from
				if _, ok := (*start)[fromIndex]; ok {
					return []int32{fromIndex, to}, nil
				} else {
					nextGoals[fromIndex] = struct{}{}
				}
			}
		}
		way, err := w.search(start, &nextGoals, depth+1)
		if err != nil {
			return nil, err
		}
		links, err := w.forwardLinks(w.pages[way[len(way)-1]])
		if err != nil {
			return nil, err
		}
		for _, to := range links {
			toIndex := to
			if _, ok := (*goal)[toIndex]; ok {
				return append(way, toIndex), nil
			}
		}
	}
	return nil, NotFound
}

func (w *Wikipedia) forwardLinks(page core.Page) ([]int32, error) {
	links := make([]int32, page.ForwardLinkLength, page.ForwardLinkLength)
	_, err := w.linkFile.Seek(int64(page.ForwardLinkIndex) * 4, 0)
	if err != nil {
		return nil, err
	}
	err = binary.Read(w.linkFile, binary.LittleEndian, links)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (w *Wikipedia) backwardLinks(page core.Page) ([]int32, error) {
	links := make([]int32, page.BackwardLinkLength, page.BackwardLinkLength)
	_, err := w.linkFile.Seek(int64(page.BackwardLinkIndex) * 4, 0)
	if err != nil {
		return nil, err
	}
	err = binary.Read(w.linkFile, binary.LittleEndian, links)
	if err != nil {
		return nil, err
	}
	return links, nil
}