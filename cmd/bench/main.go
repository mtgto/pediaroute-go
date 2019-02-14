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
		titleFile = flag.String("ip", "title.dat", "File path of titles")
		titleIndicesFile = flag.String("is", "titleIndices.dat", "File path of same titles")
		linkFile = flag.String("il", "link.dat", "File path of links")
	)

	flag.Parse()
	// overwrite by environment variables
	if file, ok := os.LookupEnv("TITLE_FILE"); ok {
		titleFile = &file
	}
	log.Printf("title file: %v\n", *titleFile)
	if file, ok := os.LookupEnv("TITLE_INDICES_FILE"); ok {
		titleIndicesFile = &file
	}
	log.Printf("title indices file: %v\n", *titleIndicesFile)
	if file, ok := os.LookupEnv("LINK_FILE"); ok {
		linkFile = &file
	}
	log.Printf("link file: %v\n", *linkFile)

	if !isFileExists(*titleFile) {
		log.Fatalf("Title file does not exists: %v", *titleFile)
	}
	if !isFileExists(*titleIndicesFile) {
		log.Fatalf("Title indices file does not exists: %v", *titleIndicesFile)
	}
	if !isFileExists(*linkFile) {
		log.Fatalf("Link file does not exists: %v", *linkFile)
	}
	log.Println("Data loading...")
	wikipedia, err := web.Load(0, *titleFile, *titleIndicesFile, *linkFile)
	if err != nil {
		panic(err)
	}
	log.Println("Data loaded.")

	start := time.Now()
	for i := 0; i < *count; i++ {
		from := wikipedia.GetRandomPage()
		to := wikipedia.GetRandomPage()
		wikipedia.Search(from, to)
	}
	log.Printf("%v routes searched. %v seconds\n", *count, time.Since(start).Seconds())
}

func isFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
