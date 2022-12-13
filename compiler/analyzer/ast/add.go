package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/file"
)

type Add struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Add) Pos() file.Pos {
	return n.LHS.Pos()
}

func (n *Add) Children() []Node {
	return []Node{n.LHS, n.RHS}
}

func (n *Add) Type() Type {
	return n.LHS.Type()
}

func (n *Add) TypeCheck() {
	lhsTyp := n.LHS.Type()
	rhsTyp := n.RHS.Type()

	if lhsTyp.IsResolved && rhsTyp.IsResolved {
		if lhsTyp.Return != rhsTyp.Return {
			n.OpPos.MarkError(fmt.Sprintf("cannot add types %v and %v", lhsTyp, rhsTyp))
		}
	}
}
