package ast

import "github.com/patrick-jessen/script/compiler/token"

type Float struct {
	Token token.Token
}

func (f *Float) Pos() token.Pos {
	return f.Token.Pos
}

func (f *Float) Type() string {
	return "float"
}

func (f *Float) String() string {
	return f.Token.String()
}
