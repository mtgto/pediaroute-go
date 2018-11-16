package main

import (
	"encoding/json"
	"flag"
	"github.com/mtgto/pediaroute-go/internal/app/web"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	var wikipedia web.Wikipedia

	var (
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
	wikipedia = web.Load(*titleFile, *titleIndicesFile, *linkFile)
	log.Println("Data loaded.")

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(statikFS)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	type Pair struct {
		From string `json:"from"`
		To string `json:"to"`
	}

	type SearchResult struct {
		Route []string `json:"route"`
		Error string `json:"error"`
	}

	http.HandleFunc("/api/search", func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var pair Pair
			var result SearchResult
			err := json.NewDecoder(r.Body).Decode(&pair)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if !wikipedia.IsWordExists(pair.From) {
				result = SearchResult{nil, pair.From + "というページがないみたい"}
			} else if !wikipedia.IsWordExists(pair.To) {
				result = SearchResult{nil, pair.To + "というページがないみたい"}
			} else {
				log.Printf("Search from \"%v\" to \"%v\"\n", pair.From, pair.To)
				route := wikipedia.Search(pair.From, pair.To)
				if route != nil {
					log.Printf("Found a route: %v\n", route)
					result = SearchResult{route, ""}
				} else {
					log.Printf("Not found a route from \"%v\" to \"%v\"\n", pair.From, pair.To)
					result = SearchResult{nil, "6回のリンクじゃ見つからなかった…ごめんね！"}
				}
			}
			err = json.NewEncoder(w).Encode(result)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	http.HandleFunc("/api/random", func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			from := wikipedia.GetRandomPage()
			to := wikipedia.GetRandomPage()
			pair := Pair{from, to}
			err := json.NewEncoder(w).Encode(pair)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
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

func isFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
