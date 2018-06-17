package ast

import "github.com/patrick-jessen/script/compiler/token"

type Identifier struct {
	token token.Token
}

func (i *Identifier) Pos() token.Pos {
	return i.token.Pos
}
