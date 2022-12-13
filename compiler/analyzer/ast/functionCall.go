package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/file"
)

type FunctionCall struct {
	Identifier    *Identifier
	Args          *FunctionCallArgs
	LastParentPos file.Pos
}

func (n *FunctionCall) Pos() file.Pos {
	return n.Identifier.Pos()
}

func (n *FunctionCall) Children() []Node {
	var nodes []Node
	for _, a := range n.Args.Args {
		nodes = append(nodes, a)
	}
	return nodes
}

func (n *FunctionCall) Name() string {
	return n.Identifier.Name()
}

func (n *FunctionCall) Type() Type {
	return n.Identifier.Type()
}

func (n *FunctionCall) SetType(t Type) {
	n.Identifier.Typ = t
}
func (n *FunctionCall) TypeCheck() {
	if !n.Type().IsResolved {
		return
	}

	numArgs := 0
	if n.Args != nil {
		n.Args.TypeCheck()
		numArgs = len(n.Args.Args)
	}

	if numArgs != len(n.Type().Args) {
		n.LastParentPos.MarkError(fmt.Sprintf(
			"incorrect number of arguments. Expected %v, got %v",
			len(n.Type().Args), numArgs,
		))
	}
	for i, a := range n.Type().Args {
		if i == numArgs {
			break
		}

		r := n.Args.Args[i].Type().Return
		if r != a {
			n.Args.Args[i].Pos().MarkError(fmt.Sprintf(
				"expected %v, got %v", a, r,
			))
		}
	}
}
