package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
)

type Statements struct {
	Stmts []Node
}

func (s *Statements) Pos() token.Pos {
	return s.Stmts[0].Pos()
}

func (n Statements) String() (out string) {
	for _, s := range n.Stmts {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}
