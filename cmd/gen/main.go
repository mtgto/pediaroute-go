package main

import (
	"flag"
	"os"

	"github.com/mtgto/pediaroute-go/internal/app/gen"
)

// ./gen -ip page.sql.gz -il pagelinks.sql.gz -o [directory]
func main() {
	var (
		language     = flag.String("lang", "", "Language ID (ja or en)")
		pageFile     = flag.String("ip", "", "File path of pages.sql.gz")
		pageLinkFile = flag.String("il", "", "File path of pages.sql.gz")
		outDir       = flag.String("o", "", "Output directory")
	)
	flag.Parse()
	if fileMode, err := os.Stat(*pageFile); err == nil && fileMode.Mode().IsRegular() {
		if fileMode, err := os.Stat(*pageLinkFile); err == nil && fileMode.Mode().IsRegular() {
			if fileMode, err := os.Stat(*outDir); err == nil && fileMode.IsDir() {
				gen.Run(*language, *pageFile, *pageLinkFile, *outDir)
				os.Exit(0)
			}
		}
	}
	flag.PrintDefaults()
}
