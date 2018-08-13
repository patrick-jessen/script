package ast

import (
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Float struct {
	Token token.Token
}

func (f *Float) Pos() file.Pos {
	return f.Token.Pos
}

func (f *Float) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "float",
	}
}

func (f *Float) String() string {
	return f.Token.String()
}
func (*Float) TypeCheck() {
}
