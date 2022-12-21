package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type Block struct {
	Statements []ast.Node
}

func (n *Block) Pos() file.Pos {
	return n.Statements[0].Pos()
}

func (n *Block) Children() []ast.Node {
	return n.Statements
}

func (n *Block) TypeCheck() {
	for _, s := range n.Statements {
		s.TypeCheck()
	}
}
