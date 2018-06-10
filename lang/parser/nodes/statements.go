package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/parser"
)

type StatementsNode struct {
	Stmts []parser.ASTNode
}

func (n StatementsNode) String() (out string) {
	for _, s := range n.Stmts {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}
