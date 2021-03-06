package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type Block struct {
	Statements []Node
}

func (b *Block) Pos() file.Pos {
	return b.Statements[0].Pos()
}

func (b *Block) String(level int) (out string) {
	for _, s := range b.Statements {
		out += s.String(level)
	}
	return
}

func (b *Block) TypeCheck() {
	for _, s := range b.Statements {
		s.TypeCheck()
	}
}
