package gen

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// 終了コード
const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
)

type CLI struct {
	OutStream, ErrStream io.Writer
}

type pageLink struct {
	// id in wikipedia
	from  int
	title string
}

type pageTitleToIndex struct {
	title string
	index int
}

func (cli *CLI) Run(pageSQLFile, pageLinkSQLFile, outDir string) int {
	flag.Parse()
	//walk()
	var pages []core.Page
	var pageLinks []uint64
	idToIndices := make(map[int]int)
	titleToIndices := make(map[string]int)
	pageFile := path.Join(outDir, "title.dat")
	if _, err := os.Stat(pageFile); err == nil {
		fmt.Fprintf(cli.ErrStream, "Load \"%s\".\n", pageFile)
		pages = core.LoadPages(pageFile)
	} else {
		fmt.Fprintf(cli.ErrStream, "Create \"%s\".\n", pageFile)
		pages = generatePages(pageSQLFile, pageFile)
	}

	for i, page := range pages {
		idToIndices[page.Id] = i
		titleToIndices[page.Title] = i
	}
	fmt.Fprintf(cli.ErrStream, "%d pages.\n", len(pages))

	titleIndicesFile := path.Join(outDir, "titleIndices.dat")
	if _, err := os.Stat(titleIndicesFile); err == nil {
		fmt.Fprintf(cli.ErrStream, "Skip to create \"%s\".\n", titleIndicesFile)
	} else {
		fmt.Fprintf(cli.ErrStream, "Create \"%s\".\n", titleIndicesFile)
		generateTitleIndices(titleIndicesFile, titleToIndices)
	}

	pageLinkFile := path.Join(outDir, "link.dat")
	if _, err := os.Stat(pageLinkFile); err == nil {
		fmt.Fprintf(cli.ErrStream, "Load \"%s\".\n", pageLinkFile)
		pageLinks = loadPageLinks(pageLinkFile)
	} else {
		fmt.Fprintf(cli.ErrStream, "Create \"%s\".\n", pageLinkFile)
		pageLinks = generatePageLinks(pageLinkSQLFile, pageLinkFile, idToIndices, titleToIndices)
	}
	fmt.Printf("%d pageLinks.\n", len(pageLinks))

	return ExitCodeOK
}

/**
 * Parse SQL insert statement of `pages` table.
 *
 * It returns only namespace = 0 (normal page)
 */
func parsePageStatement(stmt *sqlparser.Insert) ([]core.Page, error) {
	pages := make([]core.Page, 0)
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
					pages = append(pages, core.Page{Id: pageID, Title: pageTitle, IsRedirect: pageIsRedirect != 0})
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

// Import page array from sql file and generate new CSV file.
func generatePages(in string, out string) []core.Page {
	var allPages []core.Page
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
	fp, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	writer := csv.NewWriter(bufio.NewWriter(fp))
	for {
		statement, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if sq, ok := statement.(*sqlparser.Insert); ok {
			pages, err := parsePageStatement(sq)
			//fmt.Fprintf(os.Stderr, "page len %d\n", len(pages))
			if err != nil {
				panic(err)
			}
			for _, page := range pages {
				if err := writer.Write([]string{strconv.Itoa(page.Id), strconv.FormatBool(page.IsRedirect), page.Title}); err != nil {
					panic(err)
				}
			}
			allPages = append(allPages, pages...)
		}
	}
	writer.Flush()
	return allPages
}



// Generate a map from lowercase of title to indices and write it to file.
func generateTitleIndices(out string, titleToIndices map[string]int) {
	lowercaseTitleToIndices := make([]pageTitleToIndex, 0)
	for title, index := range titleToIndices {
		lowercaseTitleToIndices = append(lowercaseTitleToIndices, pageTitleToIndex{title: strings.ToLower(title), index: index})
	}
	sort.Slice(lowercaseTitleToIndices, func(i, j int) bool {
		return lowercaseTitleToIndices[i].title < lowercaseTitleToIndices[j].title
	})
	fp, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	writer := csv.NewWriter(bufio.NewWriter(fp))
	lastTitle := lowercaseTitleToIndices[0].title
	lastIndex := 0
	for i, entry := range lowercaseTitleToIndices {
		if lastTitle != entry.title {
			sameTitles := make([]string, 1)
			sameTitles[0] = lastTitle
			for j := lastIndex; j < i; j++ {
				sameTitles = append(sameTitles, strconv.Itoa(lowercaseTitleToIndices[j].index))
			}
			writer.Write(sameTitles)
			lastTitle = entry.title
			lastIndex = i
		}
	}
	writer.Flush()
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

func generatePageLinks(in string, out string, idToIndices map[int]int, titleToIndices map[string]int) []uint64 {
	var allLinks []uint64
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
	fp, err := os.Create(out)
	defer fp.Close()
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(fp)
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
						value := uint64(fromIndex)<<32 | uint64(toIndex)
						binary.Write(writer, binary.LittleEndian, value)
						allLinks = append(allLinks, value)
					}
				}
			}
		}
	}
	writer.Flush()
	return allLinks
}

func loadPageLinks(in string) []uint64 {
	var value uint64
	var allLinks []uint64
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		err := binary.Read(reader, binary.LittleEndian, &value)
		if err != io.EOF {
			allLinks = append(allLinks, value)
			break
		} else {
			panic(err)
		}
	}
	return allLinks
}
