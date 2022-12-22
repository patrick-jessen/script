package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

type String struct {
	Token token.Token
}

func (n *String) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "string",
		Pos:     n.Token.Pos,
		Literal: n.Token.Value,
	}
}
