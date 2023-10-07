package grep

import (
	"bufio"
	"io"

	"github.com/ray1729/puzzle-solver/util"
)

type DB interface {
	FindMatches(s string) []string
}

type Node struct {
	Value    byte
	Children []*Node
	Results  []string
}

func (n *Node) Add(xs []byte, s string) {
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
	child.Add(xs[1:], s)
}

func Load(r io.Reader) (DB, error) {
	db := &Node{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		s := sc.Text()
		xs := util.LowerCaseAlpha(s)
		db.Add(xs, s)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return db, nil
}

func (n *Node) FindMatches(s string) []string {
	return n.find(util.LowerCaseAlphaOrDot(s))
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
