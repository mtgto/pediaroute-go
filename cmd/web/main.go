package main

import (
	"encoding/json"
	"flag"
	"github.com/mtgto/pediaroute-go/internal/app/web"
	"github.com/rakyll/statik/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

/**
 * Error code of API
 */
const (
	NoError = iota
	NotFoundFrom
	NotFoundTo
	NotFoundRoute
)

type languageData struct {
	directory string
	pageCount uint
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	wikipedias := make(map[string]web.Wikipedia)

	var (
		japaneseDirectory = flag.String("ja", "ja", "Directory path of Japanese data")
		japanesePageCount = flag.Uint("ja-page-count", 0, "Page count of Japanese data")
		englishDirectory = flag.String("en", "en", "Directory path of English data")
		englishPageCount = flag.Uint("en-page-count", 0, "Page count of English data")
	)

	flag.Parse()
	// overwrite by environment variables
	if dir, ok := os.LookupEnv("JA"); ok {
		japaneseDirectory = &dir
	}
	if ja, ok := os.LookupEnv("JA_PAGES"); ok {
		count, err := strconv.Atoi(ja)
		if err != nil {
			panic(err)
		}
		*japanesePageCount = uint(count)
	}
	if dir, ok := os.LookupEnv("EN"); ok {
		englishDirectory = &dir
	}
	if en, ok := os.LookupEnv("EN_PAGES"); ok {
		count, err := strconv.Atoi(en)
		if err != nil {
			panic(err)
		}
		*englishPageCount = uint(count)
	}

	languageData := map[string]languageData{
		"ja": {directory: *japaneseDirectory, pageCount: *japanesePageCount},
		"en": {directory: *englishDirectory, pageCount: *englishPageCount},
	}
	for lang, langData := range languageData {
		titleFile := path.Join(langData.directory, "title.dat")
		titleIndicesFile := path.Join(langData.directory, "titleIndices.dat")
		linkFile := path.Join(langData.directory, "link.dat")
		log.Printf("Start loading for language %v\n", lang)
		wikipedia, err := web.Load(langData.pageCount, titleFile, titleIndicesFile, linkFile)
		if err != nil {
			log.Printf("Failed to load for lang %v: %v", lang, err)
		} else {
			wikipedias[lang] = wikipedia
		}
		log.Printf("Loaded for language %v\n", lang)
	}

	log.Println("Data loaded.")
	printMemory()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(statikFS)
	http.Handle("/", fileServer)
	http.Handle("/about", http.StripPrefix("/about", fileServer))
	http.Handle("/search", http.StripPrefix("/search", fileServer))

	type Pair struct {
		From string `json:"from"`
		To string `json:"to"`
	}

	type SearchResult struct {
		Route []string `json:"route"`
		Error int `json:"error"`
	}

	http.HandleFunc("/api/search", func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var pair Pair
			var result SearchResult
			lang := r.FormValue("lang")
			if wikipedia, ok := wikipedias[lang]; ok {
				err := json.NewDecoder(r.Body).Decode(&pair)
				if err != nil {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if !wikipedia.IsWordExists(pair.From) {
					result = SearchResult{nil, NotFoundFrom}
				} else if !wikipedia.IsWordExists(pair.To) {
					result = SearchResult{nil, NotFoundTo}
				} else {
					log.Printf("Language: %v, Search from \"%v\" to \"%v\"\n", lang, pair.From, pair.To)
					route := wikipedia.Search(pair.From, pair.To)
					if route != nil {
						log.Printf("Found a route: %v\n", route)
						result = SearchResult{route, NoError}
					} else {
						log.Printf("Not found a route from \"%v\" to \"%v\"\n", pair.From, pair.To)
						result = SearchResult{nil, NotFoundRoute}
					}
				}
				err = json.NewEncoder(w).Encode(result)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	http.HandleFunc("/api/random", func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			lang := r.FormValue("lang")
			if wikipedia, ok := wikipedias[lang]; ok {
				title := wikipedia.GetRandomPage()
				w.Header().Set("Content-Type", "application/json")
				log.Printf("Language: %v, Random page \"%v\"\n", lang, title)
				err := json.NewEncoder(w).Encode(title)
				if err != nil {
					log.Printf("Error while encoding json: %v\n", err)
				}
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}

func printMemory() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Printf("Memory: %+v", mem.HeapAlloc)
}

func isFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
