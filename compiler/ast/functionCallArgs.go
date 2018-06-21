package ast

import "github.com/patrick-jessen/script/compiler/token"

type FunctionCallArgs struct {
	Args []Expression
}

func (f *FunctionCallArgs) Pos() token.Pos {
	return f.Args[0].Pos()
}
func (f *FunctionCallArgs) TypeCheck(errFn ErrorFunc) {
	for _, a := range f.Args {
		a.TypeCheck(errFn)
	}
}
