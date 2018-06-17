package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionDeclNode struct {
	Identifier *TokenNode
	Block      ast.Node
}

func (n FunctionDeclNode) String() string {
	block := fmt.Sprintf("  %v", n.Block)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("FunctionDecl"),
		n.Identifier,
		strings.Replace(block, "\n", "\n  ", -1),
	)
}

func (n *FunctionDeclNode) ForwardDeclare(s *Scope) *Scope {
	return s.DeclareFunction(n.Identifier.Token, "()")
}
