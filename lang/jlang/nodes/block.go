package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
)

type Block struct {
	Statements []ast.Node
}

func (n *Block) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "block",
		Pos:      n.Statements[0].Info().Pos, // TODO: handle empty block
		Children: n.Statements,
	}
}
