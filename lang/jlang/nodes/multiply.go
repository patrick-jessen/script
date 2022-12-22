package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type Multiply struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Multiply) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "multiply",
		Pos:      n.OpPos,
		Children: []ast.Node{n.LHS, n.RHS},
	}
}

// func (n *Multiply) TypeCheck() {
// 	lhsTyp := n.LHS.Type()
// 	rhsTyp := n.RHS.Type()

// 	if lhsTyp.Return != rhsTyp.Return {
// 		n.RHS.Pos().MarkError(fmt.Sprintf("cannot multiply types %v and %v", lhsTyp, rhsTyp))
// 	}
// }
