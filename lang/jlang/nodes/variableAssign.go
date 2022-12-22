package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type VariableAssign struct {
	Identifier *Identifier
	Value      Expression
	EqPos      file.Pos
}

func (n *VariableAssign) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "variableAssign",
		Pos:      n.EqPos,
		Name:     n.Identifier.Info().Name,
		Children: []ast.Node{n.Value},
	}
}

// func (n *VariableAssign) TypeCheck() {
// 	n.Value.TypeCheck()

// 	if !n.Identifier.Type().IsCompatible(n.Value.Type()) {
// 		n.EqPos.MarkError(fmt.Sprintf("cannot assign type %v to %v",
// 			n.Value.Type(), n.Identifier.Type()))
// 	}
// }
