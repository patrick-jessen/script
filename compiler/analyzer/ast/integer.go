package ast

import (
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Integer struct {
	Token token.Token
}

func (n *Integer) Pos() file.Pos {
	return n.Token.Pos
}

func (n *Integer) Children() []Node {
	return nil
}

func (n *Integer) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "int",
	}
}

func (n *Integer) Value() string {
	return n.Token.Value
}

func (n *Integer) TypeCheck() {
}
