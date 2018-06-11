package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableDeclNode struct {
	Identifier parser.ASTNode
	Value      parser.ASTNode
}

func (n VariableDeclNode) String() string {
	val := fmt.Sprintf("  %v", n.Value)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("VariableDecl"),
		n.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}
