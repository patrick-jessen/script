package ast

import (
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type String struct {
	Token token.Token
}

func (n *String) Pos() file.Pos {
	return n.Token.Pos
}

func (n *String) Children() []Node {
	return nil
}

func (n *String) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "string",
	}
}

func (n *String) Value() string {
	return n.Token.Value
}

func (n *String) TypeCheck() {
}
