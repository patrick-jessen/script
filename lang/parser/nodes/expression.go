package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

type ExpressionNode struct {
	Expression parser.ASTNode
	Immutable  bool
}

func (n ExpressionNode) String() string {
	return fmt.Sprintf("%v", n.Expression)
}

func (n *ExpressionNode) Analyze(mod module.Module) {}

func (n *ExpressionNode) Type() string {
	panic("not implemented") // for type checking
}
