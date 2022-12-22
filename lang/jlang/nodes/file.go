package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
)

type File struct {
	Declarations []Declarable
}

func (n *File) Info() ast.NodeInfo {
	var children []ast.Node
	for _, decl := range n.Declarations {
		children = append(children, decl)
	}

	return ast.NodeInfo{
		Type:     "file",
		Pos:      n.Declarations[0].Info().Pos, // TODO: handle empty file
		Children: children,
	}
}
