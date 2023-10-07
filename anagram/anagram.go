package anagram

import (
	"bufio"
	"io"
	"sort"

	"github.com/ray1729/puzzle-solver/util"
)

type DB interface {
	FindAnagrams(s string) []string
}

func toKey(s string) string {
	xs := util.LowerCaseAlpha(s)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	return string(xs)
}

type HashDBImpl map[string][]string

func Load(r io.Reader) (DB, error) {
	db := make(HashDBImpl)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		s := sc.Text()
		k := toKey(s)
		db[k] = append(db[k], s)

	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d HashDBImpl) FindAnagrams(s string) []string {
	return d[toKey(s)]
}
