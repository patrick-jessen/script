package ast

import (
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Integer struct {
	Token token.Token
}

func (i *Integer) Pos() file.Pos {
	return i.Token.Pos
}

func (i *Integer) String(level int) string {
	return i.Token.String()
}

func (i *Integer) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "int",
	}
}
func (*Integer) TypeCheck() {
}
