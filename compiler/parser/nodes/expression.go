package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
)

type ExpressionNode struct {
	Expression ast.Node
	Immutable  bool
}

func (n ExpressionNode) String() string {
	return fmt.Sprintf("%v", n.Expression)
}

func (n *ExpressionNode) Type() string {
	panic("not implemented") // for type checking
}
