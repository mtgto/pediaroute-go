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

type idIndexEntry struct {
	id    int32
	index uint32
}

type titleIndexEntry struct {
	title string
	index uint32
}

type linkPair struct {
	from uint32
	to   uint32
}

func buildIDIndex(pages []page) []idIndexEntry {
	idIndex := make([]idIndexEntry, len(pages))
	for i, p := range pages {
		idIndex[i] = idIndexEntry{id: p.id, index: uint32(i)}
	}
	sort.Slice(idIndex, func(i, j int) bool {
		return idIndex[i].id < idIndex[j].id
	})
	return idIndex
}

func buildTitleIndex(pages []page) []titleIndexEntry {
	titleIndex := make([]titleIndexEntry, len(pages))
	for i, p := range pages {
		titleIndex[i] = titleIndexEntry{title: p.title, index: uint32(i)}
	}
	sort.Slice(titleIndex, func(i, j int) bool {
		return titleIndex[i].title < titleIndex[j].title
	})
	return titleIndex
}

func lookupID(idIndex []idIndexEntry, id int) (uint32, bool) {
	i := sort.Search(len(idIndex), func(k int) bool {
		return int(idIndex[k].id) >= id
	})
	if i < len(idIndex) && int(idIndex[i].id) == id {
		return idIndex[i].index, true
	}
	return 0, false
}

func lookupTitle(titleIndex []titleIndexEntry, title string) (uint32, bool) {
	i := sort.Search(len(titleIndex), func(k int) bool {
		return titleIndex[k].title >= title
	})
	if i < len(titleIndex) && titleIndex[i].title == title {
		return titleIndex[i].index, true
	}
	return 0, false
}

func Run(languageId, pageSQLFile, pageLinkSQLFile, outDir string) error {
	language := core.Language{
		Id:        languageId,
		PageFile:  "page.dat",
		TitleFile: "title.dat",
		LinkFile:  "link.dat",
	}
	var pages []page
	pageFile := path.Join(outDir, language.PageFile)
	titleFile := path.Join(outDir, language.TitleFile)
	log.Printf("Load \"%s\".\n", pageSQLFile)
	pages = loadPages(pageSQLFile)
	language.PageCount = uint32(len(pages))

	idIndex := buildIDIndex(pages)
	titleIndex := buildTitleIndex(pages)
	log.Printf("%d pages.\n", language.PageCount)

	pageLinkFile := path.Join(outDir, language.LinkFile)
	log.Printf("Create \"%s\".\n", pageLinkFile)
	linkCount := generatePageLinks(pageLinkSQLFile, pageLinkFile, pages, idIndex, titleIndex)
	language.LinkCount = linkCount
	log.Printf("%v links loaded.\n", language.LinkCount)
	generatePages(pageFile, titleFile, pages)

	configFile := path.Join(outDir, "config.json")
	configBytes, err := json.MarshalIndent(language, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configFile, configBytes, 0644)
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
						id:         int32(pageID),
						title:      pageTitle,
						isRedirect: pageIsRedirect != 0,
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
		return strings.ToLower(allPages[i].title) < strings.ToLower(allPages[j].title)
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

func generatePageLinks(in string, out string, pages []page, idIndex []idIndexEntry, titleIndex []titleIndexEntry) uint64 {
	// Phase A: stream pagelinks SQL and write valid (fromIndex, toIndex) pairs to temp file
	tmpFile, err := os.CreateTemp("", "pediaroute-pairs-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	pairWriter := bufio.NewWriterSize(tmpFile, 4*1024*1024)

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
			for _, pl := range pageLinks {
				toIndex, ok := lookupTitle(titleIndex, pl.title)
				if !ok {
					continue
				}
				fromIndex, ok := lookupID(idIndex, pl.from)
				if !ok {
					continue
				}
				if err := binary.Write(pairWriter, binary.LittleEndian, fromIndex); err != nil {
					panic(err)
				}
				if err := binary.Write(pairWriter, binary.LittleEndian, toIndex); err != nil {
					panic(err)
				}
			}
		}
	}
	if err := pairWriter.Flush(); err != nil {
		panic(err)
	}

	// Phase B: read all pairs into memory
	fi, err := tmpFile.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := fi.Size()
	pairCount := fileSize / 8

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	pairReader := bufio.NewReaderSize(tmpFile, 64*1024*1024)

	pairs := make([]linkPair, pairCount)
	for i := range pairs {
		if err := binary.Read(pairReader, binary.LittleEndian, &pairs[i].from); err != nil {
			panic(err)
		}
		if err := binary.Read(pairReader, binary.LittleEndian, &pairs[i].to); err != nil {
			panic(err)
		}
	}
	tmpFile.Close()

	// Phase C: open output file
	fp, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	linkWriter := bufio.NewWriterSize(fp, 4*1024*1024)

	// Phase D: sort by from, write forward links
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].from != pairs[j].from {
			return pairs[i].from < pairs[j].from
		}
		return pairs[i].to < pairs[j].to
	})

	var linkIndex uint32
	for i := 0; i < len(pairs); {
		currentFrom := pairs[i].from
		j := i
		for j < len(pairs) && pairs[j].from == currentFrom {
			if err := binary.Write(linkWriter, binary.LittleEndian, pairs[j].to); err != nil {
				panic(err)
			}
			j++
		}
		pages[currentFrom].forwardLinkIndex = linkIndex
		pages[currentFrom].forwardLinkLength = uint32(j - i)
		linkIndex += uint32(j - i)
		i = j
	}

	// Phase E: sort by to, write backward links
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].to != pairs[j].to {
			return pairs[i].to < pairs[j].to
		}
		return pairs[i].from < pairs[j].from
	})

	for i := 0; i < len(pairs); {
		currentTo := pairs[i].to
		j := i
		for j < len(pairs) && pairs[j].to == currentTo {
			if err := binary.Write(linkWriter, binary.LittleEndian, pairs[j].from); err != nil {
				panic(err)
			}
			j++
		}
		pages[currentTo].backwardLinkIndex = linkIndex
		pages[currentTo].backwardLinkLength = uint32(j - i)
		linkIndex += uint32(j - i)
		i = j
	}

	if err := linkWriter.Flush(); err != nil {
		panic(err)
	}
	return uint64(len(pairs))
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
