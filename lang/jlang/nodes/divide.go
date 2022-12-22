package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type Divide struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Divide) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "divide",
		Pos:      n.OpPos,
		Children: []ast.Node{n.LHS, n.RHS},
	}
}
