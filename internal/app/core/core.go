package core

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type Page struct {
	Id                 int32
	Title              string
	IsRedirect         bool
	ForwardLinkIndex   int32
	ForwardLinkLength  uint32
	BackwardLinkIndex  int32
	BackwardLinkLength uint32
}

func LoadPages(in string) []Page {
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 7
	reader.ReuseRecord = true
	pages := make([]Page, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		pageID, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		pageIsRedirect, err := strconv.ParseBool(record[1])
		if err != nil {
			panic(err)
		}
		forwardLinkIndex, err := strconv.Atoi(record[3])
		if err != nil {
			panic(err)
		}
		forwardLinkLength, err := strconv.Atoi(record[4])
		if err != nil {
			panic(err)
		}
		backwardLinkIndex, err := strconv.Atoi(record[5])
		if err != nil {
			panic(err)
		}
		backwardLinkLength, err := strconv.Atoi(record[6])
		if err != nil {
			panic(err)
		}
		pages = append(pages, Page{
			Id:                 int32(pageID),
			Title:              CopyString(record[2]),
			IsRedirect:         pageIsRedirect,
			ForwardLinkIndex:   int32(forwardLinkIndex),
			ForwardLinkLength:  uint32(forwardLinkLength),
			BackwardLinkIndex:  int32(backwardLinkIndex),
			BackwardLinkLength: uint32(backwardLinkLength),
		})
	}
	return pages
}

// deepcopy
func CopyString(s string) string {
	bytes := make([]byte, len(s))
	copy(bytes, s)
	return string(bytes)
}
