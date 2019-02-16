package main

import (
	"flag"
	"github.com/mtgto/pediaroute-go/internal/app/web"
	"log"
	"os"
	"time"
)

func main() {
	var wikipedia web.Wikipedia

	var (
		count = flag.Int("n", 10000, "Count of search")
		pageFile = flag.String("ip", "page.dat", "File path of pages")
		titleFile = flag.String("it", "title.dat", "File path of titles")
		linkFile = flag.String("il", "link.dat", "File path of links")
	)

	flag.Parse()
	// overwrite by environment variables
	if file, ok := os.LookupEnv("PAGE_FILE"); ok {
		pageFile = &file
	}
	log.Printf("page file: %v\n", *pageFile)
	if file, ok := os.LookupEnv("TITLE_FILE"); ok {
		titleFile = &file
	}
	log.Printf("title file: %v\n", *titleFile)
	if file, ok := os.LookupEnv("LINK_FILE"); ok {
		linkFile = &file
	}
	log.Printf("link file: %v\n", *linkFile)

	if !isFileExists(*titleFile) {
		log.Fatalf("Title file does not exists: %v", *titleFile)
	}
	if !isFileExists(*linkFile) {
		log.Fatalf("Link file does not exists: %v", *linkFile)
	}
	log.Println("Data loading...")
	wikipedia, err := web.Load(0, *pageFile, *titleFile, *linkFile)
	if err != nil {
		panic(err)
	}
	log.Println("Data loaded.")

	start := time.Now()
	for i := 0; i < *count; i++ {
		from, _ := wikipedia.GetRandomPage()
		to, _ := wikipedia.GetRandomPage()
		wikipedia.Search(from, to)
	}
	log.Printf("%v routes searched. %v seconds\n", *count, time.Since(start).Seconds())
}

func isFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
