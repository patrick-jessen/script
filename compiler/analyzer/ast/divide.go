package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type Divide struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Divide) Pos() file.Pos {
	return n.LHS.Pos()
}

func (n *Divide) Children() []Node {
	return []Node{n.LHS, n.RHS}
}

func (n *Divide) Type() Type {
	return n.LHS.Type()
}

func (n *Divide) TypeCheck() {
}
