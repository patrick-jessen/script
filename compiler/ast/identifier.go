package ast

import (
	"github.com/patrick-jessen/script/compiler/token"
)

type Identifier struct {
	Token token.Token
}

func (i *Identifier) Pos() token.Pos {
	return i.Token.Pos
}

func (i Identifier) String() string {
	return i.Token.String()
}

func (i *Identifier) Type() string {
	panic("not implemented")
}
