package main

import (
	"encoding/json"
	"flag"
	"github.com/mtgto/pediaroute-go/internal/app/core"
	"github.com/mtgto/pediaroute-go/internal/app/web"
	"github.com/rakyll/statik/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
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

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	wikipedias := make(map[string]web.Wikipedia)

	var (
		japaneseConfigFile = flag.String("ja", "ja", "Configuration file path of Japanese data")
		englishConfigFile = flag.String("en", "en", "Configuration file path of English data")
	)

	flag.Parse()
	// overwrite by environment variables
	if configFile, ok := os.LookupEnv("JA"); ok {
		japaneseConfigFile = &configFile
	}
	if configFile, ok := os.LookupEnv("EN"); ok {
		englishConfigFile = &configFile
	}

	for lang, langFile := range map[string]string{"ja": *japaneseConfigFile, "en": *englishConfigFile} {
		language, err := core.LoadLanguage(langFile)
		if err != nil {
			log.Fatalf("Failed to load language file: %v", err)
		}
		pageFile := path.Join(path.Dir(langFile), language.PageFile)
		titleFile := path.Join(path.Dir(langFile), language.TitleFile)
		linkFile := path.Join(path.Dir(langFile), language.LinkFile)
		log.Printf("Start loading for language %v\n", language.Id)
		wikipedia, err := web.Load(language.PageCount, pageFile, titleFile, linkFile)
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
				title, err := wikipedia.GetRandomPage()
				if err != nil {
					log.Printf("Error while get random page: %v\n", err)
				}
				w.Header().Set("Content-Type", "application/json")
				log.Printf("Language: %v, Random page \"%v\"\n", lang, title)
				err = json.NewEncoder(w).Encode(title)
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
