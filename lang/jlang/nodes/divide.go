package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
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

func (n *Divide) Children() []ast.Node {
	return []ast.Node{n.LHS, n.RHS}
}

func (n *Divide) Type() ast.Type {
	return n.LHS.Type()
}

func (n *Divide) TypeCheck() {
}
