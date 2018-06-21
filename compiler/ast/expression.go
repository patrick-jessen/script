package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
)

type ExpressionNode struct {
	Expression Expression
	Immutable  bool
}

func (e *ExpressionNode) Type() Type {
	return e.Expression.Type()
}

func (e *ExpressionNode) Pos() token.Pos {
	return e.Expression.Pos()
}

func (e ExpressionNode) String() string {
	return fmt.Sprintf("%v", e.Expression)
}
func (*ExpressionNode) TypeCheck(errFn ErrorFunc) {
}
