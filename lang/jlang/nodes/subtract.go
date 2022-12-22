package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type Subtract struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Subtract) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "subtract",
		Pos:      n.OpPos,
		Children: []ast.Node{n.LHS, n.RHS},
	}
}
