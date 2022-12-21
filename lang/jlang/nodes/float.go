package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Float struct {
	Token token.Token
}

func (n *Float) Pos() file.Pos {
	return n.Token.Pos
}

func (n *Float) Children() []ast.Node {
	return nil
}

func (n *Float) Type() ast.Type {
	return ast.Type{
		IsResolved: true,
		Return:     "float",
	}
}

func (n *Float) TypeCheck() {
}
