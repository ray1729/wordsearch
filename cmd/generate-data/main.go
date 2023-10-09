package main

import (
	"encoding/gob"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ray1729/wordsearch/anagram"
	"github.com/ray1729/wordsearch/match"
)

var wordlist = flag.String("wordlist", "", "Path to wordlist")
var dataDir = flag.String("data", "", "Path to output directory")

func main() {
	flag.Parse()
	os.MkdirAll(*dataDir, 0755)
	words, err := os.Open(*wordlist)
	if err != nil {
		log.Fatal(err)
	}
	anagramDB, err := anagram.Load(words)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeData(*dataDir, "anagram.bin", anagramDB); err != nil {
		log.Fatal(err)
	}

	words.Seek(0, 0)
	matchDB, err := match.Load(words)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeData(*dataDir, "match.bin", matchDB); err != nil {
		log.Fatal(err)
	}
}

func writeData(dirName string, fileName string, data interface{}) error {
	path := filepath.Join(dirName, fileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(data)
}
