package core

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Page struct {
	Id                 int
	Title              string
	IsRedirect         bool
	ForwardLinkIndex   int
	ForwardLinkLength  int
	BackwardLinkIndex  int
	BackwardLinkLength int
}

func LoadPages(in string) []Page {
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	pages := make([]Page, 0, len(records))
	for _, record := range records {
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
			Id:                 pageID,
			Title:              record[2],
			IsRedirect:         pageIsRedirect,
			ForwardLinkIndex:   forwardLinkIndex,
			ForwardLinkLength:  forwardLinkLength,
			BackwardLinkIndex:  backwardLinkIndex,
			BackwardLinkLength: backwardLinkLength,
		})
	}
	return pages
}
