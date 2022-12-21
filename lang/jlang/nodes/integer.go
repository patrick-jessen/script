package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Integer struct {
	Token token.Token
}

func (n *Integer) Pos() file.Pos {
	return n.Token.Pos
}

func (n *Integer) Children() []ast.Node {
	return nil
}

func (n *Integer) Type() ast.Type {
	return ast.Type{
		IsResolved: true,
		Return:     "int",
	}
}

func (n *Integer) Value() string {
	return n.Token.Value
}

func (n *Integer) TypeCheck() {
}
