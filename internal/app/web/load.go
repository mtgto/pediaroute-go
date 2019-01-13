package web

import (
	"encoding/csv"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"io"
	"os"
	"strconv"
)

func Load(titleFile, titleIndicesFile, linkFile string) Wikipedia {
	pages := core.LoadPages(titleFile)
	lowercaseTitleToIndices := loadLowercaseTitleToIndices(titleIndicesFile)
	links, err := os.Open(linkFile)
	if err != nil {
		panic(err)
	}
	return Wikipedia{
		pages: pages,
		lowercaseTitleToIndices: lowercaseTitleToIndices,
		linkFile: links,
	}
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
