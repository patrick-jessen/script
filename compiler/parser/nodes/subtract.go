package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/utils/color"
)

type SubtractNode struct {
	LHS ast.Node
	RHS ast.Node
}

func (n SubtractNode) String() string {
	lhs := fmt.Sprintf("  %v", n.LHS)
	rhs := fmt.Sprintf("  %v", n.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Subtract"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}

func (n *SubtractNode) Analyze(mod module.Module) {}
