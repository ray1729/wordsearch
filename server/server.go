package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ray1729/puzzle-solver/anagram"
	"github.com/ray1729/puzzle-solver/grep"
)

func New(assetsPath string, grepDB grep.DB, anagramDB anagram.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))
	mux.HandleFunc("/", getHomePage)
	mux.HandleFunc("/search", getResults(grepDB, anagramDB))
	return withRequestLogger(mux)
}

func withRequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.URL.Query())
		h.ServeHTTP(w, r)
	})
}

func getResults(grepDB grep.DB, anagramDB anagram.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("error parsing form: %v", err)
			http.Error(w, "error parsing form", http.StatusBadRequest)
			return
		}
		mode := r.Form.Get("mode")
		pattern := r.Form.Get("pattern")
		templateParams := struct {
			Preamble string
			Results  []string
		}{}
		switch mode {
		case "match":
			templateParams.Results = grepDB.FindMatches(pattern)
			if len(templateParams.Results) > 0 {
				templateParams.Preamble = fmt.Sprintf("Matches for %q:", pattern)
			} else {
				templateParams.Preamble = fmt.Sprintf("Found no matches for %q", pattern)
			}
		case "anagrams":
			templateParams.Results = anagramDB.FindAnagrams(pattern)
			if len(templateParams.Results) > 0 {
				templateParams.Preamble = fmt.Sprintf("Anagrams of %q:", pattern)
			} else {
				templateParams.Preamble = fmt.Sprintf("Found no anagrams of %q", pattern)
			}
		case "clear":
			// pass
		default:
			log.Printf("invalid mode: %s", mode)
			http.Error(w, "invalid mode", http.StatusBadRequest)
			return
		}
		renderTemplate(w, results, templateParams)
	}
}

func getHomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, home, nil)
}

func renderTemplate(w http.ResponseWriter, t *template.Template, params any) {
	err := t.Execute(w, params)
	if err != nil {
		log.Printf("Error rendering template %s: %v", t.Name(), err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
