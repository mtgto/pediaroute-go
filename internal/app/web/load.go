package web

import (
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"os"
)

func Load(pageCount uint32, pageFile, titleFile, forwardLinkFile, backwardLinkFile string) (Wikipedia, error) {
	var w Wikipedia
	pages := core.LoadPages(pageCount, pageFile)
	titleFp, err := os.Open(titleFile)
	if err != nil {
		return w, err
	}
	forwardFp, err := os.Open(forwardLinkFile)
	if err != nil {
		return w, err
	}
	backwardFp, err := os.Open(backwardLinkFile)
	if err != nil {
		return w, err
	}
	w = Wikipedia{
		pages:            pages,
		titleFile:        titleFp,
		forwardLinkFile:  forwardFp,
		backwardLinkFile: backwardFp,
	}
	return w, nil
}
