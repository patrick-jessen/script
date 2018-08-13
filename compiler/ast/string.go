package ast

import (
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
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
