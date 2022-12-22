package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
)

type FunctionDecl struct {
	Identifier *Identifier
	Args       []*Identifier
	Block      *Block
}

func (n *FunctionDecl) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "functionDecl",
		Pos:      n.Identifier.Info().Pos,
		Name:     n.Identifier.Info().Name,
		Children: []ast.Node{n.Block},
	}
}

// func (n *File) FunctionDecl() ast.NodeInfo {
// 	var children []ast.Node
// 	for _, decl := range n.Declarations {
// 		children = append(children, decl)
// 	}

// 	return ast.NodeInfo{
// 		Type: "functionDecl",

// 		Pos:      n.Declarations[0].Info().Pos, // TODO: handle empty file
// 		Children: children,
// 	}
// }

// func (n *FunctionDecl) Init() {
// 	var args []string
// 	for _, a := range n.Args {
// 		args = append(args, a.Type().Return)
// 	}

// 	n.Identifier.Typ = ast.Type{
// 		IsResolved: true,
// 		IsFunction: true,
// 		Return:     "void",
// 		Args:       args,
// 	}
// }
