package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
)

type VariableRef struct {
	Identifier *Identifier
}

func (n *VariableRef) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type: "variableRef",
		Pos:  n.Identifier.Info().Pos,
		Name: n.Identifier.Info().Name,
	}
}
