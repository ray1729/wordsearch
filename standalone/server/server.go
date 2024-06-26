package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/ray1729/wordsearch/anagram"
	"github.com/ray1729/wordsearch/match"
)

func New(matchDB match.DB, anagramDB anagram.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", searchHandler(matchDB, anagramDB))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, home, nil)
	})
	return withRequestLogger(mux)
}

func withRequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.URL.Query())
		h.ServeHTTP(w, r)
	})
}

func searchHandler(matchDB match.DB, anagramDB anagram.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("error parsing form: %v", err)
			http.Error(w, "error parsing form", http.StatusBadRequest)
			return
		}
		switch r.Form.Get("mode") {
		case "match":
			params := matchResults(matchDB, r.Form.Get("pattern"))
			renderTemplate(w, results, params)
		case "anagrams":
			params := anagramResults(anagramDB, r.Form.Get("pattern"))
			renderTemplate(w, results, params)
		default:
			renderTemplate(w, results, ResultParams{})
		}
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
