package gen

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"github.com/xwb1989/sqlparser"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type page struct {
	id                 int32
	title              string
	titleLowercase     string
	isRedirect         bool
	forwardLinkIndex   uint32
	forwardLinkLength  uint32
	backwardLinkIndex  uint32
	backwardLinkLength uint32
}

type pageLink struct {
	// id in wikipedia
	from  int
	title string
}

type pageTitleToIndex struct {
	title string
	index uint32
}

func Run(languageId, pageSQLFile, pageLinkSQLFile, outDir string) error {
	language := core.Language{
		Id:        languageId,
		PageFile:  "page.dat",
		TitleFile: "title.dat",
		LinkFile:  "link.dat",
	}
	var pages []page
	idToIndices := make(map[int]uint32)
	titleToIndices := make(map[string]uint32)
	pageFile := path.Join(outDir, language.PageFile)
	titleFile := path.Join(outDir, language.TitleFile)
	log.Printf("Load \"%s\".\n", pageSQLFile)
	pages = loadPages(pageSQLFile)
	language.PageCount = uint32(len(pages))

	for i, page := range pages {
		idToIndices[int(page.id)] = uint32(i)
		titleToIndices[page.title] = uint32(i)
	}
	log.Printf("%d pages.\n", language.PageCount)

	pageLinkFile := path.Join(outDir, language.LinkFile)
	log.Printf("Create \"%s\".\n", pageLinkFile)
	linkCount := generatePageLinks(pageLinkSQLFile, pageLinkFile, pages, idToIndices, titleToIndices)
	language.LinkCount = linkCount
	log.Printf("%v links loaded.\n", language.LinkCount)
	generatePages(pageFile, titleFile, pages)

	configFile := path.Join(outDir, "config.json")
	configBytes, err := json.MarshalIndent(language, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configFile, configBytes, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}

/**
 * Parse SQL insert statement of `pages` table.
 *
 * It returns pages which namespace == 0 (normal page)
 */
func parsePageStatement(stmt *sqlparser.Insert) ([]page, error) {
	pages := make([]page, 0)
	var columnIndex, pageID, pageNamespace, pageIsRedirect int
	var pageTitle string
	var err error
	const (
		columnID = iota
		columnNamespace
		columnTitle
		columnRestrictions
		columnCounter
		columnIsRedirect
	)
	err = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		//fmt.Printf("%T\n", node)
		switch node := node.(type) {
		case sqlparser.Values, sqlparser.Exprs, *sqlparser.NullVal:
			return true, nil
		case sqlparser.ValTuple:
			columnIndex = 0
			return true, nil
		case *sqlparser.SQLVal:
			if columnIndex == columnID {
				pageID, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
			} else if columnIndex == columnNamespace {
				pageNamespace, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
			} else if columnIndex == columnTitle {
				//pageTitle = sqlparser.String(node)
				pageTitle = string(node.Val)
			} else if columnIndex == columnIsRedirect {
				pageIsRedirect, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
				if pageNamespace == 0 {
					pages = append(pages, page{
						id:             int32(pageID),
						title:          pageTitle,
						titleLowercase: strings.ToLower(pageTitle),
						isRedirect:     pageIsRedirect != 0,
					})
				}
			}
			columnIndex++
			return false, nil
		default:
			panic(fmt.Sprintf("Unknown type! %T", node))
		}
	}, stmt.Rows)
	return pages, err
}

// Import page array from sql file
func loadPages(in string) []page {
	var allPages []page
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	tokens := sqlparser.NewTokenizer(reader)
	for {
		statement, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if sq, ok := statement.(*sqlparser.Insert); ok {
			pages, err := parsePageStatement(sq)
			if err != nil {
				panic(err)
			}
			allPages = append(allPages, pages...)
		}
	}
	// sort all pages by title in lowercase
	sort.Slice(allPages, func(i, j int) bool {
		return allPages[i].titleLowercase < allPages[j].titleLowercase
	})
	return allPages
}

// Parse SQL insert statement of `pagelinks` table.
func parsePageLinkStatement(stmt *sqlparser.Insert) ([]pageLink, error) {
	pagelinks := make([]pageLink, 0)
	var columnIndex, from, namespace, fromNamespace int
	var title string
	var err error
	const (
		columnFrom = iota
		columnNamespace
		columnTitle
		columnFromNamespace
	)
	err = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		switch node := node.(type) {
		case sqlparser.Values, sqlparser.Exprs, *sqlparser.NullVal:
			return true, nil
		case sqlparser.ValTuple:
			columnIndex = 0
			return true, nil
		case *sqlparser.SQLVal:
			if columnIndex == columnFrom {
				from, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
			} else if columnIndex == columnNamespace {
				namespace, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
			} else if columnIndex == columnTitle {
				title = string(node.Val)
			} else if columnIndex == columnFromNamespace {
				fromNamespace, err = strconv.Atoi(sqlparser.String(node))
				if err != nil {
					panic(fmt.Sprintf("Parse error %s", err))
				}
				if namespace == 0 && fromNamespace == 0 {
					pagelinks = append(pagelinks, pageLink{from: from, title: title})
				}
			}
			columnIndex++
			return false, nil
		default:
			panic(fmt.Sprintf("Unknown type! %T", node))
		}
	}, stmt.Rows)
	return pagelinks, err
}

func generatePageLinks(in string, out string, pages []page, idToIndices map[int]uint32, titleToIndices map[string]uint32) uint64 {
	forwardLinks := make([][]uint32, len(pages), len(pages))
	backwardLinks := make([][]uint32, len(pages), len(pages))
	var linkLength uint64 = 0
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	tokens := sqlparser.NewTokenizer(reader)
	for {
		statement, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if sq, ok := statement.(*sqlparser.Insert); ok {
			pageLinks, err := parsePageLinkStatement(sq)
			if err != nil {
				panic(err)
			}
			for _, pageLink := range pageLinks {
				if toIndex, ok := titleToIndices[pageLink.title]; ok {
					if fromIndex, ok := idToIndices[pageLink.from]; ok {
						forwardLinks[fromIndex] = append(forwardLinks[fromIndex], toIndex)
						backwardLinks[toIndex] = append(backwardLinks[toIndex], fromIndex)
						linkLength++
					}
				}
			}
		}
	}
	fp, err := os.Create(out)
	defer fp.Close()
	if err != nil {
		panic(err)
	}

	linkIndex := 0
	writer := bufio.NewWriter(fp)
	for i, links := range forwardLinks {
		pages[i].forwardLinkIndex = uint32(linkIndex)
		pages[i].forwardLinkLength = uint32(len(links))
		for _, toIndex := range links {
			err := binary.Write(writer, binary.LittleEndian, uint32(toIndex))
			if err != nil {
				panic(err)
			}
		}
		linkIndex += len(links)
	}
	for i, links := range backwardLinks {
		pages[i].backwardLinkIndex = uint32(linkIndex)
		pages[i].backwardLinkLength = uint32(len(links))
		for _, fromIndex := range links {
			err := binary.Write(writer, binary.LittleEndian, uint32(fromIndex))
			if err != nil {
				panic(err)
			}
		}
		linkIndex += len(links)
	}
	writer.Flush()
	return linkLength
}

func generatePages(pageFile, titleFile string, pages []page) {
	pageFp, err := os.Create(pageFile)
	if err != nil {
		panic(err)
	}
	defer pageFp.Close()
	titleFp, err := os.Create(titleFile)
	if err != nil {
		panic(err)
	}
	defer titleFp.Close()
	pageWriter := bufio.NewWriter(pageFp)
	titleWriter := bufio.NewWriter(titleFp)
	var titleOffset uint32
	for _, page := range pages {
		// Write page title + "\0"
		titleLength, err := titleWriter.WriteString(page.title)
		if err != nil {
			panic(err)
		}
		err = titleWriter.WriteByte(0)
		if err != nil {
			panic(err)
		}

		err = binary.Write(pageWriter, binary.LittleEndian, page.id)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, titleOffset)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, uint16(titleLength))
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, page.isRedirect)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, page.forwardLinkIndex)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, page.forwardLinkLength)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, page.backwardLinkIndex)
		if err != nil {
			panic(err)
		}
		err = binary.Write(pageWriter, binary.LittleEndian, page.backwardLinkLength)
		if err != nil {
			panic(err)
		}

		titleOffset += uint32(titleLength) + 1
	}
	err = pageWriter.Flush()
	if err != nil {
		panic(err)
	}
	err = titleWriter.Flush()
	if err != nil {
		panic(err)
	}
}

//func loadPageLinks(in string) []uint64 {
//	var value uint64
//	var allLinks []uint64
//	file, err := os.Open(in)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//	reader := bufio.NewReader(file)
//	for {
//		err := binary.Read(reader, binary.LittleEndian, &value)
//		if err != io.EOF {
//			allLinks = append(allLinks, value)
//			break
//		} else {
//			panic(err)
//		}
//	}
//	return allLinks
//}
