package web

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"io"
	"os"
	"strconv"
)

func Load(titleFile, titleIndicesFile, linkFile string) Wikipedia {
	titles, redirectSet := loadTitles(titleFile)
	lowercaseTitleToIndices := loadLowercaseTitleToIndices(titleIndicesFile)
	forwardLinks, backwardLinks := loadLinks(linkFile, len(titles))
	return Wikipedia{titles: titles, redirectSet: redirectSet, lowercaseTitleToIndices: lowercaseTitleToIndices, forwardLinks: forwardLinks, backwardLinks: backwardLinks}
}

func loadTitles(in string) ([]string, map[int]struct{}) {
	pages := core.LoadPages(in)
	titles := make([]string, len(pages))
	redirects := make(map[int]struct{}, 0)
	for i, page := range pages {
		titles[i] = page.Title
		if page.IsRedirect {
			redirects[i] = struct{}{}
		}
	}
	return titles, redirects
}

func loadLowercaseTitleToIndices(in string) map[string][]int {
	lowercaseTitleToIndices := make(map[string][]int, 0)
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		lowercaseTitle := record[0]
		indices := make([]int, 0, len(record)-1)
		for _, indexStr := range record[1:] {
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				panic(err)
			}
			indices = append(indices, index)
		}
		lowercaseTitleToIndices[lowercaseTitle] = indices
	}
	return lowercaseTitleToIndices
}

func loadLinks(in string, titleCount int) ([][]int, [][]int) {
	var value uint64
	forwardLinks := make([][]int, titleCount)
	backwardLinks := make([][]int, titleCount)
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		err := binary.Read(reader, binary.LittleEndian, &value)
		if err == io.EOF {
			break
		} else if err == nil {
			from := int(value >> 32)
			to := int(value & 0xFFFFFFFF)
			forwardLinks[from] = append(forwardLinks[from], to)
			backwardLinks[to] = append(backwardLinks[to], from)
		} else {
			panic(err)
		}
	}
	return forwardLinks, backwardLinks
}
