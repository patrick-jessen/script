package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/utils/color"
)

type AddNode struct {
	LHS ast.Node
	RHS ast.Node
}

func (n AddNode) String() string {
	lhs := fmt.Sprintf("  %v", n.LHS)
	rhs := fmt.Sprintf("  %v", n.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Add"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}
