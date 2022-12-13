package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/file"
)

type Multiply struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Multiply) Pos() file.Pos {
	return n.LHS.Pos()
}

func (n *Multiply) Children() []Node {
	return []Node{n.LHS, n.RHS}
}

func (n *Multiply) Type() Type {
	return n.LHS.Type()
}

func (n *Multiply) TypeCheck() {
	lhsTyp := n.LHS.Type()
	rhsTyp := n.RHS.Type()

	if lhsTyp.Return != rhsTyp.Return {
		n.RHS.Pos().MarkError(fmt.Sprintf("cannot multiply types %v and %v", lhsTyp, rhsTyp))
	}
}
