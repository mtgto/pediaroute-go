package main

import (
	"flag"
	"os"

	"github.com/mtgto/pediaroute-go/internal/app/gen"
)

// ./gen -ip page.sql.gz -il pagelinks.sql.gz -o [directory]
func main() {
	var (
		pageFile     = flag.String("ip", "", "File path of pages.sql.gz")
		pageLinkFile = flag.String("il", "", "File path of pages.sql.gz")
		outDir       = flag.String("o", "", "Output directory")
	)
	flag.Parse()
	if _, err := os.Stat(*pageFile); err == nil {
		if _, err := os.Stat(*pageLinkFile); err == nil {
			cli := &gen.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
			os.Exit(cli.Run(*pageFile, *pageLinkFile, *outDir))
		}
	}
	flag.PrintDefaults()
}
