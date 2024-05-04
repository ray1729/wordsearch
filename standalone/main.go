package main

import (
	"embed"
	"encoding/gob"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/ray1729/wordsearch/anagram"
	"github.com/ray1729/wordsearch/match"
	"github.com/ray1729/wordsearch/standalone/server"
)

//go:embed data/*
var fs embed.FS

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
	if err := initializeDB(); err != nil {
		log.Fatalf("error initializing databases: %v", err)
	}
}

func main() {
	var listenAddr = flag.String("listen", ":8000", "Address or port to listen on, prefix with unix: to listen on a Unix domain socket")
	flag.Parse()

	server := http.Server{
		Handler: server.New(matchDB, anagramDB),
	}

	var listener net.Listener
	var err error
	if strings.HasPrefix(*listenAddr, "unix:") {
		listener, err = net.Listen("unix", strings.TrimPrefix(*listenAddr, "unix:"))
	} else {
		listener, err = net.Listen("tcp", *listenAddr)
	}
	if err != nil {
		log.Fatalf("Error listening on %s: %v", *listenAddr, err)
	}
	log.Printf("Listening on %s", *listenAddr)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
