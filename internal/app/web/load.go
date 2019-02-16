package web

import (
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"os"
)

func Load(pageCount uint32, pageFile, titleFile, linkFile string) (Wikipedia, error) {
	var w Wikipedia
	pages := core.LoadPages(pageCount, pageFile)
	titleFp, err := os.Open(titleFile)
	if err != nil {
		return w, err
	}
	linkFp, err := os.Open(linkFile)
	if err != nil {
		return w, err
	}
	w = Wikipedia{
		pages: pages,
		titleFile: titleFp,
		linkFile: linkFp,
	}
	return w, nil
}
