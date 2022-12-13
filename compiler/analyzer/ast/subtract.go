package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type Subtract struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Subtract) Pos() file.Pos {
	return n.LHS.Pos()
}

func (n *Subtract) Children() []Node {
	return []Node{n.LHS, n.RHS}
}

func (n *Subtract) Type() Type {
	return n.LHS.Type()
}

func (n *Subtract) TypeCheck() {
}
