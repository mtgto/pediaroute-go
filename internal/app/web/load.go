package web

import (
	"encoding/csv"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"io"
	"os"
	"strconv"
)

func Load(titleFile, titleIndicesFile, linkFile string) (Wikipedia, error) {
	var w Wikipedia
	pages := core.LoadPages(titleFile)
	lowercaseTitleToIndices := loadLowercaseTitleToIndices(titleIndicesFile)
	links, err := os.Open(linkFile)
	if err != nil {
		return w, err
	}
	w = Wikipedia{
		pages: pages,
		lowercaseTitleToIndices: lowercaseTitleToIndices,
		linkFile: links,
	}
	return w, nil
}

func loadLowercaseTitleToIndices(in string) map[string][]int32 {
	lowercaseTitleToIndices := make(map[string][]int32, 0)
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
		lowercaseTitle := core.CopyString(record[0])
		indices := make([]int32, 0, len(record)-1)
		for _, indexStr := range record[1:] {
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				panic(err)
			}
			indices = append(indices, int32(index))
		}
		lowercaseTitleToIndices[lowercaseTitle] = indices
	}
	return lowercaseTitleToIndices
}
