package ast

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

type String struct {
	Token token.Token
}

func (s *String) String() string {
	return s.Token.String()
}

func (s *String) Pos() file.Pos {
	return s.Token.Pos
}

func (s *String) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "string",
	}
}
func (*String) TypeCheck() {
}
