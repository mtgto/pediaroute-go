package core

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// WikipediaFiles defines JSON structure for data set of language.
type Language struct {
	Id                          string `json:"id"`
	PageCount                   uint32 `json:"page_count"`
	PageFile                    string `json:"page_file"`
	TitleFile                   string `json:"title_file"`
	LinkCount                   uint64 `json:"link_count"`
	LinkFile                    string `json:"link_file"`
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
	for i := 0; uint32(i) < pageCount; i += 1 {
		page := Page{}
		err := binary.Read(reader, binary.LittleEndian, &page.Id)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.TitleOffset)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.TitleLength)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.IsRedirect)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.ForwardLinkIndex)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.ForwardLinkLength)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.BackwardLinkIndex)
		if err != nil {
			panic(err)
		}
		err = binary.Read(reader, binary.LittleEndian, &page.BackwardLinkLength)
		if err != nil {
			panic(err)
		}
		pages = append(pages, page)
	}
	return pages
}

// load config.json
func LoadLanguage(in string) (*Language, error) {
	bytes, err := ioutil.ReadFile(in)
	if err != nil {
		log.Fatal(err)
	}
	language := Language{}
	err = json.Unmarshal(bytes, &language)
	if err != nil {
		return nil, err
	}
	return &language, nil
}
