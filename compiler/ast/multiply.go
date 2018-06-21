package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Multiply struct {
	LHS   Expression
	RHS   Expression
	OpPos token.Pos
}

func (m *Multiply) Pos() token.Pos {
	return m.LHS.Pos()
}
func (m *Multiply) Type() Type {
	return m.LHS.Type()
}

func (m Multiply) String() string {
	lhs := fmt.Sprintf("  %v", m.LHS)
	rhs := fmt.Sprintf("  %v", m.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Multiply"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}

func (m *Multiply) TypeCheck(errFn ErrorFunc) {
	lhsTyp := m.LHS.Type()
	rhsTyp := m.RHS.Type()

	if lhsTyp.Return != rhsTyp.Return {
		errFn(m.RHS.Pos(), fmt.Sprintf("cannot multiply types %v and %v", lhsTyp, rhsTyp))
	}
}
