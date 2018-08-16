package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
)

type Multiply struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (m *Multiply) Pos() file.Pos {
	return m.LHS.Pos()
}
func (m *Multiply) Type() Type {
	return m.LHS.Type()
}

func (m Multiply) String(level int) string {
	lhs := fmt.Sprintf("  %v", m.LHS)
	rhs := fmt.Sprintf("  %v", m.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Multiply"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}

func (m *Multiply) TypeCheck() {
	lhsTyp := m.LHS.Type()
	rhsTyp := m.RHS.Type()

	if lhsTyp.Return != rhsTyp.Return {
		m.RHS.Pos().MarkError(fmt.Sprintf("cannot multiply types %v and %v", lhsTyp, rhsTyp))
	}
}
