package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/utils/color"
)

type MultiplyNode struct {
	LHS ast.Node
	RHS ast.Node
}

func (n MultiplyNode) String() string {
	lhs := fmt.Sprintf("  %v", n.LHS)
	rhs := fmt.Sprintf("  %v", n.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Multiply"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}

func (n *MultiplyNode) Analyze(mod module.Module) {}
