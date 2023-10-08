package anagram

import (
	"bufio"
	"io"
	"sort"

	"github.com/ray1729/wordsearch/util"
)

type DB interface {
	FindAnagrams(s string) []string
	Add(s string)
}

type HashDBImpl map[string][]string

func New() HashDBImpl {
	return make(HashDBImpl)
}

func (db HashDBImpl) Add(s string) {
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
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d HashDBImpl) FindAnagrams(s string) []string {
	return d[toKey(s)]
}
