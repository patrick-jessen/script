package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type Add struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (n *Add) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "add",
		Pos:      n.OpPos,
		Children: []ast.Node{n.LHS, n.RHS},
	}
}

// func (n *Add) TypeCheck() {
// 	lhsTyp := n.LHS.Type()
// 	rhsTyp := n.RHS.Type()

// 	if lhsTyp.IsResolved && rhsTyp.IsResolved {
// 		if lhsTyp.Return != rhsTyp.Return {
// 			n.OpPos.MarkError(fmt.Sprintf("cannot add types %v and %v", lhsTyp, rhsTyp))
// 		}
// 	}
// }
