package nodes

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type FunctionCall struct {
	Identifier    *Identifier
	Args          *FunctionCallArgs
	LastParentPos file.Pos
}

func (n *FunctionCall) Info() ast.NodeInfo {
	var children []ast.Node
	for _, arg := range n.Args.Args {
		children = append(children, arg)
	}

	return ast.NodeInfo{
		Type:     "functionCall",
		Pos:      n.Identifier.Info().Pos,
		Name:     n.Identifier.Info().Name,
		Children: children,
	}
}

// func (n *FunctionCall) TypeCheck() {
// 	if !n.Type().IsResolved {
// 		return
// 	}

// 	numArgs := 0
// 	if n.Args != nil {
// 		n.Args.TypeCheck()
// 		numArgs = len(n.Args.Args)
// 	}

// 	if numArgs != len(n.Type().Args) {
// 		n.LastParentPos.MarkError(fmt.Sprintf(
// 			"incorrect number of arguments. Expected %v, got %v",
// 			len(n.Type().Args), numArgs,
// 		))
// 	}
// 	for i, a := range n.Type().Args {
// 		if i == numArgs {
// 			break
// 		}

// 		r := n.Args.Args[i].Type().Return
// 		if r != a {
// 			n.Args.Args[i].Pos().MarkError(fmt.Sprintf(
// 				"expected %v, got %v", a, r,
// 			))
// 		}
// 	}
// }
