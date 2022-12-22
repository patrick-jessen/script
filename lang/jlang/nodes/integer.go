package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

type Integer struct {
	Token token.Token
}

func (n *Integer) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "integer",
		Pos:     n.Token.Pos,
		Literal: n.Token.Value,
	}
}
