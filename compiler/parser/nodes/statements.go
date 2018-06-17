package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
)

type StatementsNode struct {
	Stmts []ast.Node
}

func (n StatementsNode) String() (out string) {
	for _, s := range n.Stmts {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}
