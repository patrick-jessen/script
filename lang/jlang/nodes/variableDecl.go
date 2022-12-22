package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
)

type VariableDecl struct {
	Identifier *Identifier
	Value      Expression
}

func (n *VariableDecl) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "variableDecl",
		Pos:      n.Identifier.Info().Pos,
		Name:     n.Identifier.Info().Name,
		Children: []ast.Node{n.Value},
	}
}

// func (n *VariableDecl) TypeCheck() {
// 	n.Value.TypeCheck()

// 	if !n.Identifier.Typ.IsResolved {
// 		n.Identifier.Typ = n.Value.Type()
// 	}

// 	if !n.Identifier.Type().IsCompatible(n.Value.Type()) {
// 		n.Value.Pos().MarkError(fmt.Sprintf("cannot assign type %v to %v", n.Value.Type(), n.Identifier.Type()))
// 	}
// }
