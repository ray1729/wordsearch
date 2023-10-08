package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	"github.com/ray1729/puzzle-solver/anagram"
	"github.com/ray1729/puzzle-solver/match"
	"github.com/ray1729/puzzle-solver/server"
)

var matchDB match.DB
var anagramDB anagram.DB

func init() {
	f, err := os.Open("wordlist.txt")
	if err != nil {
		log.Fatalf("Error opening word list: %v", err)
	}
	defer f.Close()
	matchDB = match.New()
	anagramDB = anagram.New()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := sc.Text()
		matchDB.Add(s)
		anagramDB.Add(s)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("Error loading databases: %v", err)
	}
}

func main() {
	s := server.New("./assets", matchDB, anagramDB)
	address := ":8000"
	log.Printf("Listening on %s", address)
	if err := http.ListenAndServe(address, s); err != nil {
		log.Fatal(err)
	}
}
