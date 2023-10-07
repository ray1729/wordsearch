package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ray1729/puzzle-solver/anagram"
	"github.com/ray1729/puzzle-solver/grep"
	"github.com/ray1729/puzzle-solver/server"
)

var grepDB grep.DB
var anagramDB anagram.DB

func init() {
	f, err := os.Open("/usr/share/dict/british-english-huge")
	if err != nil {
		log.Fatalf("Error opening word list: %v", err)
	}
	defer f.Close()
	grepDB, err = grep.Load(f)
	if err != nil {
		log.Fatalf("Error loading grep database:  %v", err)
	}
	f.Seek(0, 0)
	anagramDB, err = anagram.Load(f)
	if err != nil {
		log.Fatalf("Error loading anagram database: %v", err)
	}
}

func main() {
	s := server.New("./assets", grepDB, anagramDB)
	address := ":8000"
	log.Printf("Listening on %s", address)
	if err := http.ListenAndServe(address, s); err != nil {
		log.Fatal(err)
	}
}
