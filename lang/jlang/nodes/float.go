package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

type Float struct {
	Token token.Token
}

func (n *Float) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "float",
		Pos:     n.Token.Pos,
		Literal: n.Token.Value,
	}
}
