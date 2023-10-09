package match

import (
	"bufio"
	"io"

	"github.com/ray1729/wordsearch/util"
)

type DB struct {
	Root *Node
}

func New() DB {
	return DB{Root: &Node{}}
}

func (db DB) Add(s string) {
	xs := util.LowerCaseAlpha(s)
	db.Root.add(xs, s)
}

type Node struct {
	Value    byte
	Children []*Node
	Results  []string
}

func (n *Node) add(xs []byte, s string) {
	if len(xs) == 0 {
		n.Results = append(n.Results, s)
		return
	}
	x := xs[0]
	var child *Node
	for _, c := range n.Children {
		if c.Value == x {
			child = c
			break
		}
	}
	if child == nil {
		child = &Node{Value: x}
		n.Children = append(n.Children, child)
	}
	child.add(xs[1:], s)
}

func Load(r io.Reader) (DB, error) {
	db := New()
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		db.Add(sc.Text())
	}
	return db, sc.Err()
}

func (db DB) FindMatches(s string) []string {
	return db.Root.find(util.LowerCaseAlphaOrDot(s))
}

func (n *Node) find(xs []byte) []string {
	if len(xs) == 0 {
		return n.Results
	}
	x := xs[0]
	if x == '.' {
		var result []string
		for _, c := range n.Children {
			result = append(result, c.find(xs[1:])...)
		}
		return result
	}
	for _, c := range n.Children {
		if c.Value == x {
			return c.find(xs[1:])
		}
	}
	return nil
}
