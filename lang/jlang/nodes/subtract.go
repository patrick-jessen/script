package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
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

func (n *Subtract) Children() []ast.Node {
	return []ast.Node{n.LHS, n.RHS}
}

func (n *Subtract) Type() ast.Type {
	return n.LHS.Type()
}

func (n *Subtract) TypeCheck() {
}
