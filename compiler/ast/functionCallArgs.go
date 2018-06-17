package ast

import "github.com/patrick-jessen/script/compiler/token"

type FunctionCallArgs struct {
	Args []Node
}

func (f *FunctionCallArgs) Pos() token.Pos {
	return f.Args[0].Pos()
}
