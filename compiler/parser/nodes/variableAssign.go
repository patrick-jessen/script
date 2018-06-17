package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableAssignNode struct {
	Identifier ast.Node
	Value      ast.Node
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
