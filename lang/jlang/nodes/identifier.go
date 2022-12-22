package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

type Identifier struct {
	Symbol token.Token
	Module token.Token
}

func (n *Identifier) Info() ast.NodeInfo {
	name := n.Symbol.Value
	if len(n.Module.Value) > 0 {
		name = n.Module.Value + "." + name
	}
	return ast.NodeInfo{
		Type: "identifier",
		Pos:  n.Symbol.Pos,
		Name: name,
	}
}
