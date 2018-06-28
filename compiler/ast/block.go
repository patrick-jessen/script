package ast

import "github.com/patrick-jessen/script/compiler/token"

type Block struct {
	Statements []Node
}

func (b *Block) Pos() token.Pos {
	return b.Statements[0].Pos()
}

func (b *Block) String() (out string) {
	for _, s := range b.Statements {
		out += s.String() + "\n"
	}
	return
}
func (b *Block) TypeCheck(errFn ErrorFunc) {
	for _, s := range b.Statements {
		s.TypeCheck(errFn)
	}
}
