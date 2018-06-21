package ast

type FunctionDeclArgs struct {
	Names []*Identifier
	Types []*Identifier
}

func (*FunctionDeclArgs) TypeCheck(errFn ErrorFunc) {
}
