package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableAssignNode struct {
	Identifier parser.ASTNode
	Value      parser.ASTNode
}

func (n VariableAssignNode) String() string {
	val := fmt.Sprintf("  %v", n.Value)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("VariableAssign"),
		n.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}

func (n *VariableAssignNode) Analyze(mod module.Module) {}
