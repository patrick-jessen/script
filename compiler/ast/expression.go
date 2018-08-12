package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/file"
)

type ExpressionNode struct {
	Expression Expression
	Immutable  bool
}

func (e *ExpressionNode) Type() Type {
	return e.Expression.Type()
}

func (e *ExpressionNode) Pos() file.Pos {
	return e.Expression.Pos()
}

func (e ExpressionNode) String() string {
	return fmt.Sprintf("%v", e.Expression)
}
func (*ExpressionNode) TypeCheck() {
}
