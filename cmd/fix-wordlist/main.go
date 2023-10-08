package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	dec := charmap.ISO8859_1.NewDecoder()
	sc := bufio.NewScanner(dec.Reader(os.Stdin))
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}
