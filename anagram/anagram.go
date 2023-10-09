package anagram

import (
	"bufio"
	"io"
	"sort"

	"github.com/ray1729/wordsearch/util"
)

type DB map[string][]string

func New() DB {
	return make(DB)
}

func (db DB) Add(s string) {
	k := toKey(s)
	db[k] = append(db[k], s)
}

func toKey(s string) string {
	xs := util.LowerCaseAlpha(s)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	return string(xs)
}

func Load(r io.Reader) (DB, error) {
	db := New()
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		db.Add(sc.Text())

	}
	return db, sc.Err()
}

func (d DB) FindAnagrams(s string) []string {
	return d[toKey(s)]
}
