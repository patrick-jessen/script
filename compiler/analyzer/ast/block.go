package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type Block struct {
	Statements []Node
}

func (n *Block) Pos() file.Pos {
	return n.Statements[0].Pos()
}

func (n *Block) Children() []Node {
	return n.Statements
}

func (n *Block) TypeCheck() {
	for _, s := range n.Statements {
		s.TypeCheck()
	}
}
