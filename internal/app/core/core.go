package core

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"os"
)

// Language defines JSON structure for data set of wikipediaPages.
type Language struct {
	Id        string `json:"id"`
	PageCount uint32 `json:"page_count"`
	PageFile  string `json:"page_file"`
	TitleFile string `json:"title_file"`
	LinkCount uint64 `json:"link_count"`
	LinkFile  string `json:"link_file"`
	Version   string `json:"version"` // YYYYMMDD
}

type Page struct {
	Id                 uint32
	TitleOffset        uint32
	TitleLength        uint16
	IsRedirect         bool
	ForwardLinkIndex   uint32
	ForwardLinkLength  uint32
	BackwardLinkIndex  uint32
	BackwardLinkLength uint32
}

func LoadPages(pageCount uint32, in string) []Page {
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	pages := make([]Page, 0, pageCount)
	for i := 0; uint32(i) < pageCount; i++ {
		page := Page{}
		err := binary.Read(reader, binary.LittleEndian, &page)
		if err != nil {
			panic(err)
		}
		pages = append(pages, page)
	}
	return pages
}

// load config.json
func LoadLanguage(in string) (*Language, error) {
	bytes, err := os.ReadFile(in)
	if err != nil {
		return nil, err
	}
	language := Language{}
	err = json.Unmarshal(bytes, &language)
	if err != nil {
		return nil, err
	}
	return &language, nil
}
