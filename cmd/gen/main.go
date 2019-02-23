package main

import (
	"flag"
	"log"
	"os"

	"github.com/mtgto/pediaroute-go/internal/app/gen"
)

// ./gen -ip page.sql.gz -il pagelinks.sql.gz -o [directory]
func main() {
	var (
		language     = flag.String("lang", "", "Language ID (ja or en)")
		pageFile     = flag.String("ip", "", "File path of page.sql.gz")
		pageLinkFile = flag.String("il", "", "File path of pagelinks.sql.gz")
		outDir       = flag.String("o", "", "Output directory")
	)
	flag.Parse()
	if fileMode, err := os.Stat(*pageFile); err != nil || !fileMode.Mode().IsRegular() {
		log.Println("No page sql file.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if fileMode, err := os.Stat(*pageLinkFile); err != nil || !fileMode.Mode().IsRegular() {
		log.Println("No pagelinks sql file.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if fileMode, err := os.Stat(*outDir); err != nil || !fileMode.IsDir() {
		log.Println("No output directory.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	gen.Run(*language, *pageFile, *pageLinkFile, *outDir)
}
