package ast

import (
	"github.com/patrick-jessen/script/utils/file"
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

func (e ExpressionNode) String(level int) string {
	return e.Expression.String(level)
}
func (*ExpressionNode) TypeCheck() {
}
