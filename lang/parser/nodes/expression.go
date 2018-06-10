package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/parser"
)

type ExpressionNode struct {
	Expression parser.ASTNode
}

func (n ExpressionNode) String() string {
	return fmt.Sprintf("%v", n.Expression)
}
