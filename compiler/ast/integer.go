package ast

import "github.com/patrick-jessen/script/compiler/token"

type Integer struct {
	Token token.Token
}

func (i *Integer) Pos() token.Pos {
	return i.Token.Pos
}

func (i *Integer) String() string {
	return i.Token.String()
}

func (i *Integer) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "int",
	}
}
func (*Integer) TypeCheck(errFn ErrorFunc) {
}