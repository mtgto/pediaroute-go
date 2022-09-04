package web

import (
	"encoding/binary"
	"errors"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// Data structure
type Wikipedia struct {
	pages     []core.Page
	titleFile *os.File
	linkFile  *os.File
}

var NotFound = errors.New("not found.")

const MaxDepth = 6

func (w *Wikipedia) Search(from, to string) []string {
	start := w.generateWordSet(from)
	goal := w.generateWordSet(to)
	log.Printf("start: %v, goal: %v\n", start, goal)
	way, err := w.search(&start, &goal, 0)
	if err != nil {
		return nil
	} else {
		result := make([]string, len(way))
		for i, index := range way {
			title, err := w.title(w.pages[index])
			if err != nil {
				log.Printf("err: %v", err)
				return nil
			}
			result[i] = title
		}
		return result
	}
}

func (w *Wikipedia) IsWordExists(word string) bool {
	indices := w.generateWordSet(word)
	return len(indices) > 0
}

// return whether word is only redirect page title
func (w *Wikipedia) IsWordRedirect(word string) bool {
	indices := w.generateWordSet(word)
	if len(indices) > 0 {
		for index, _ := range indices {
			if !w.pages[index].IsRedirect {
				return false
			}
		}
		return true
	}
	return false
}

func (w *Wikipedia) GetRandomPage() (string, error) {
	for {
		title, err := w.title(w.pages[rand.Intn(len(w.pages))])
		if err != nil {
			return "", err
		}
		if !w.IsWordRedirect(title) {
			return title, nil
		}
	}
}

func (w *Wikipedia) generateWordSet(word string) map[uint32]struct{} {
	lowercaseWord := strings.ToLower(word)
	set := make(map[uint32]struct{}, 0)
	index := sort.Search(len(w.pages), func(i int) bool {
		title, err := w.title(w.pages[i])
		if err != nil {
			log.Printf("Error while generateWordSet: %v", err)
		}
		return strings.ToLower(title) >= lowercaseWord
	})
	if index == len(w.pages) {
		return set
	}
	for i := index; i >= 0; i -= 1 {
		title, err := w.title(w.pages[i])
		if err != nil {
			log.Printf("Error while generateWordSet: %v", err)
			continue
		}
		if strings.ToLower(title) == lowercaseWord {
			set[uint32(i)] = struct{}{}
		} else {
			break
		}
	}
	for i := index + 1; i < len(w.pages); i += 1 {
		title, err := w.title(w.pages[i])
		if err != nil {
			log.Printf("Error while generateWordSet: %v", err)
			continue
		}
		if strings.ToLower(title) == lowercaseWord {
			set[uint32(i)] = struct{}{}
		} else {
			break
		}
	}
	return set
}

func (w *Wikipedia) search(start, goal *map[uint32]struct{}, depth int) ([]uint32, error) {
	if depth >= MaxDepth {
		return nil, NotFound
	}
	if len(*start) < len(*goal) {
		nextStarts := make(map[uint32]struct{}, 0)
		for from, _ := range *start {
			links, err := w.forwardLinks(w.pages[from])
			if err != nil {
				return nil, err
			}
			for _, to := range links {
				toIndex := to
				if _, ok := (*goal)[toIndex]; ok {
					return []uint32{from, toIndex}, nil
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
				return append([]uint32{fromIndex}, way...), nil
			}
		}
	} else {
		nextGoals := make(map[uint32]struct{}, 0)
		for to, _ := range *goal {
			links, err := w.backwardLinks(w.pages[to])
			if err != nil {
				log.Printf("err: %v\n", err)
				return nil, err
			}
			for _, from := range links {
				fromIndex := from
				if _, ok := (*start)[fromIndex]; ok {
					return []uint32{fromIndex, to}, nil
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

func (w *Wikipedia) forwardLinks(page core.Page) ([]uint32, error) {
	links := make([]uint32, page.ForwardLinkLength, page.ForwardLinkLength)
	_, err := w.linkFile.Seek(int64(page.ForwardLinkIndex)*4, 0)
	if err != nil {
		return nil, err
	}
	err = binary.Read(w.linkFile, binary.LittleEndian, links)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (w *Wikipedia) backwardLinks(page core.Page) ([]uint32, error) {
	links := make([]uint32, page.BackwardLinkLength, page.BackwardLinkLength)
	_, err := w.linkFile.Seek(int64(page.BackwardLinkIndex)*4, 0)
	if err != nil {
		return nil, err
	}
	err = binary.Read(w.linkFile, binary.LittleEndian, links)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (w *Wikipedia) title(page core.Page) (string, error) {
	_, err := w.titleFile.Seek(int64(page.TitleOffset), 0)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, page.TitleLength)
	_, err = w.titleFile.ReadAt(bytes, int64(page.TitleOffset))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
