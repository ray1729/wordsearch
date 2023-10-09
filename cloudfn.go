package wordsearch

import (
	"embed"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/ray1729/wordsearch/anagram"
	"github.com/ray1729/wordsearch/match"
	"github.com/rs/cors"
)

var anagramDB anagram.DB
var matchDB match.DB

func initializeDB() error {
	anagrams, err := fs.Open("data/anagram.bin")
	if err != nil {
		return err
	}
	defer anagrams.Close()
	if err := gob.NewDecoder(anagrams).Decode(&anagramDB); err != nil {
		return err
	}
	matches, err := fs.Open("data/match.bin")
	if err != nil {
		return err
	}
	defer matches.Close()
	if err := gob.NewDecoder(matches).Decode(&matchDB); err != nil {
		return err
	}
	return nil
}

func init() {
	log.Println("Initializing databases")
	if err := initializeDB(); err != nil {
		log.Fatal(err)
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodPost},
		AllowCredentials: false,
		MaxAge:           3600,
	})
	log.Println("Registering HTTP function with the Functions Framework")
	functions.HTTP("WordSearch", func(w http.ResponseWriter, r *http.Request) {
		corsHandler.ServeHTTP(w, r, handleFormSubmission)
	})
}

func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	mode := r.Form.Get("mode")
	pattern := r.Form.Get("pattern")
	if len(pattern) == 0 {
		http.Error(w, "Missing pattern", http.StatusBadRequest)
		return
	}
	switch mode {
	case "match":
		results := matchResults(matchDB, pattern)
		renderTemplate(w, resultsTmpl, results)
	case "anagrams":
		results := anagramResults(anagramDB, pattern)
		renderTemplate(w, resultsTmpl, results)
	default:
		log.Printf("invalid mode: %s", mode)
		http.Error(w, fmt.Sprintf("Invalid mode: %s", mode), http.StatusBadRequest)
	}
}

func anagramResults(db anagram.DB, pattern string) ResultParams {
	var params ResultParams
	params.Results = db.FindAnagrams(pattern)
	if len(params.Results) > 0 {
		params.Preamble = fmt.Sprintf("Anagrams of %q:", pattern)
	} else {
		params.Preamble = fmt.Sprintf("Found no anagrams of %q", pattern)
	}
	sort.Slice(params.Results, func(i, j int) bool { return params.Results[i] < params.Results[j] })
	return params
}

func matchResults(db match.DB, pattern string) ResultParams {
	var params ResultParams
	params.Results = db.FindMatches(pattern)
	if len(params.Results) > 0 {
		params.Preamble = fmt.Sprintf("Matches for %q:", pattern)
	} else {
		params.Preamble = fmt.Sprintf("Found no matches for %q", pattern)
	}
	sort.Slice(params.Results, func(i, j int) bool { return params.Results[i] < params.Results[j] })
	return params
}

func renderTemplate(w http.ResponseWriter, t *template.Template, params any) {
	err := t.Execute(w, params)
	if err != nil {
		log.Printf("Error rendering template %s: %v", t.Name(), err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type ResultParams struct {
	Preamble string
	Results  []string
}

var resultsTmpl = template.Must(template.New("results").Parse(`
{{ with .Preamble }}
<p>{{ . }}</p>
{{ end }}
<ul>
{{ range .Results }}
  <li>{{.}}</li>
{{ end }}
</ul>
`))

//go:embed data/*
var fs embed.FS
